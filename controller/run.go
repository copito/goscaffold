package controller

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/copito/goscaffold/core"
	"github.com/copito/goscaffold/entity"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/kluctl/go-jinja2"
)

var (
	rgxHooksFolder *regexp.Regexp = regexp.MustCompile(`\/hooks$`)
	rgxHooksFile   *regexp.Regexp = regexp.MustCompile(`hooks(\/|\\)(pre_prompt|pre_gen_project|post_gen_project)\\.(py|go|sh)$`)

	// rgxPrePromptHooksFile   *regexp.Regexp = regexp.MustCompile(`hooks(\/|\\)(pre_prompt)\\.(py|go|sh)$`)
	// rgxPreProjectHooksFile  *regexp.Regexp = regexp.MustCompile(`hooks(\/|\\)(pre_gen_project)\\.(py|go|sh)$`)
	// rgxPostProjectHooksFile *regexp.Regexp = regexp.MustCompile(`hooks(\/|\\)(post_gen_project)\\.(py|go|sh)$`)
)

func Run(cmd *cobra.Command, args []string) {
	// Get Logger
	logger := cmd.Context().Value("logger").(*slog.Logger)
	// logger.Debug("Testing debug logger")
	// logger.Info("Testing info logger")
	// logger.Warn("Testing warn logger")
	// logger.Error("Testing error logger")

	// 1. Getting Path
	logger.Debug("Getting path provided...")
	var runPath string
	if len(args) == 0 {
		logger.Info("No path provided, assuming path is current path: .")
		runPath = "."
	} else {
		runPath = args[0]
	}

	// Check if path exists
	isExists, err := core.PathExists(runPath)
	if err != nil {
		logger.Error("Path provided does not exist...")
		os.Exit(1)
	}

	if !isExists {
		logger.Error("Path provided does not exist...")
		os.Exit(1)
	}

	// 2. Load config file
	logger.Debug("Loading configuration file...")
	configFilePath, err := cmd.Flags().GetString("config")
	if err != nil {
		configFilePath = "./config.yaml"
	}

	extension := path.Ext(configFilePath)[1:]
	dirConfig := path.Dir(configFilePath)
	basePath := core.FileNameWithoutExtension(path.Base(configFilePath))

	if basePath == "base" && extension == "yaml" {
		logger.Error("base.yaml is the only name that cannot be used for the configuration file")
		os.Exit(1)
	}

	viper.SetConfigName(basePath)
	viper.SetConfigType(extension)
	viper.AddConfigPath(dirConfig)

	// Setting prefix for all env variables: SCAFFOLD_ID => viper.Get("ID")
	// viper.SetEnvPrefix("SCAFFOLD")

	err = viper.ReadInConfig()
	if err != nil {
		logger.Error("Unable to find and/or load config file", "config", configFilePath)
		os.Exit(1)
	}

	// Parse configurations
	promptConfig := entity.Prompt{}
	err = viper.Unmarshal(&promptConfig)
	if err != nil {
		logger.Error("Unable to load config file", "config", configFilePath)
		os.Exit(1)
	}

	// Extract keys (for sorting index)
	keys := make([]string, 0, len(promptConfig.Items))
	for key := range promptConfig.Items {
		keys = append(keys, key)
	}

	// Sort keys based on OrderID
	sort.Slice(keys, func(i, j int) bool {
		return promptConfig.Items[keys[i]].OrderID < promptConfig.Items[keys[j]].OrderID
	})

	paramChoice := make(map[string]string)

	// ask questions about config (settle variables)
	// Loop through all prompt based configs (based on data type)
	for _, key := range keys {

		item := promptConfig.Items[key]
		item.Key = key

		// Private variables check (_)
		if strings.HasPrefix(key, "_") {
			paramChoice[key] = fmt.Sprintf("%v", item.DefaultValue)
			continue
		}

		// Request Data from Users
		if len(item.Options) >= 1 {
			itemOptions := core.InterfaceSliceToStringSlice(item.Options)
			result := core.SingleSelectPrompt(fmt.Sprintf("Select %s [%s]", key, item.DefaultValue.(string)), itemOptions)
			paramChoice[key] = result
			continue
		}

		switch v := item.DefaultValue.(type) {
		case string:
			result := core.StringPrompt(fmt.Sprintf("Select %s [%s]", key, v), v)
			paramChoice[key] = result
			continue
		case int:
			result := core.NumberPrompt(fmt.Sprintf("Select %s [%v]", key, item.DefaultValue), fmt.Sprintf("%v", item.DefaultValue))
			paramChoice[key] = result
			continue
		case float32:
			result := core.NumberPrompt(fmt.Sprintf("Select %s [%v]", key, item.DefaultValue), fmt.Sprintf("%v", item.DefaultValue))
			paramChoice[key] = result
			continue
		case float64:
			result := core.NumberPrompt(fmt.Sprintf("Select %s [%v]", key, item.DefaultValue), fmt.Sprintf("%v", item.DefaultValue))
			paramChoice[key] = result
			continue
		case bool:
			result := core.BoolPrompt(fmt.Sprintf("Select %s [%t]", key, item.DefaultValue), "FALSE", false)
			paramChoice[key] = result
			continue
		default:
			logger.Error("unexpected type %T", v)
			os.Exit(1)
		}
	}

	jj, err := jinja2.NewJinja2("FolderFileName", 1, jinja2.WithGlobal("scaffold", paramChoice))
	if err != nil {
		logger.Error("Unable prepare chosen parameters for Jinja Templating...")
		os.Exit(1)
	}
	defer jj.Close()

	// TODO: send it to a file (if running under debug)
	logger.Debug("New Compiled Results", "params", paramChoice)

	// 3. pre-hooks
	preHookPath := path.Join(runPath, "hooks", "pre_gen_project.go")
	hasPreGenProjectHook, _ := core.PathExists(preHookPath)
	if hasPreGenProjectHook {
		logger.Info("Running pre_gen_hook...")
		err = core.RenderFileContent(preHookPath, jj)
		if err != nil {
			logger.Error("Rendering pre-hook caused the application to crash...")
			os.Exit(1)
		}
	}

	// 4. Generate output folder to add output here
	isDryRun, _ := cmd.Flags().GetBool("dry-run")
	outputFolder := "output"
	outputBasePath := path.Join(runPath, outputFolder)
	logger.Info(fmt.Sprintf("Output Path => %s", outputBasePath))
	if !isDryRun {
		err = os.Mkdir(outputBasePath, os.FileMode(0o755))
		if err != nil {
			logger.Error("Could not create output folder")
			os.Exit(1)
		}
	}

	// rollback output folder
	rollbackChan := make(chan interface{}, 1)
	rollbackOutput := func(rollbackChan chan interface{}, outputBasePath string) {
		defer close(rollbackChan)
		<-rollbackChan

		logger.Info("Invoked rollback - removing output folder...")
		err := os.RemoveAll(outputBasePath)
		if err != nil {
			logger.Error("Error cleaning up output folder...")
			os.Exit(1)
		}
	}
	go rollbackOutput(rollbackChan, outputBasePath)

	// 5. Walk through every folder/file and change names + file data
	err = filepath.Walk(runPath, func(pathValue string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip own project folder
		if runPath == pathValue {
			// If using: /home/user/Documents/scaffold/example => runPath
			// and the run path is the same then it should be skipped
			return nil
		}

		// Skip - Bypass config file
		if path.Base(pathValue) == path.Base(configFilePath) {
			// Skip configuration file from walk
			return nil
		}

		// Hooks folder bypass
		matchedHookFolder := rgxHooksFolder.MatchString(pathValue)
		if matchedHookFolder {
			// Skip any hooks folders
			return nil
		}

		// Hooks bypass
		matchedHookFile := rgxHooksFile.MatchString(pathValue)
		if matchedHookFile {
			// Skip any hooks
			return nil
		}

		// Skip Output created folder
		if strings.HasPrefix(pathValue, outputBasePath) {
			// If it being outputted to the output file then
			// it should be ignored (especially to avoid recursive loop)
			return nil
		}

		// TODO: Copy file to output folder -> also transforming using Jinja2
		logger.Info(pathValue, "size", info.Size(), "is_dir", info.Mode().IsDir(), "is_file", info.Mode().IsRegular())
		deltaPath := core.DeltaRelativePath(runPath, pathValue)
		newFullPath := path.Join(outputBasePath, deltaPath)

		// Jinja template path name
		newFullPathRendered, err := jj.RenderString(newFullPath)
		if err != nil {
			// rollbackChan <- true
			fmt.Println("rendered new path name but failed!!")
			os.Exit(1)
		}

		logger.Info("jinja template", "templated", newFullPath, "rendered", newFullPathRendered)

		switch mode := info.Mode(); {
		case mode.IsDir():
			// Folder/Directory
			// Create folder
			err = os.MkdirAll(newFullPathRendered, os.FileMode(0o755))
			if err != nil {
				rollbackChan <- true
				time.Sleep(time.Second)
				os.Exit(1)
			}

		case mode.IsRegular():
			// File
			bytesProcessed, err := core.PathCopy(pathValue, newFullPathRendered)
			if err != nil {
				rollbackChan <- true
				time.Sleep(time.Second)
				os.Exit(1)
			}
			fmt.Print("Processed: ", bytesProcessed)

			// Render this file content
			err = core.RenderFileContent(newFullPathRendered, jj)
			if err != nil {
				fmt.Println("rendering file (using jinja) failed!!")
				rollbackChan <- true
				time.Sleep(time.Second)
				os.Exit(1)
			}
		}

		// core.PathCopy(pathValue, out)
		// file := path.Base(pathValue)

		return nil
	})
	if err != nil {
		log.Println(err)
	}

	// 6. TODO: post-hook
	hasPostGenProjectHook, _ := core.PathExists(path.Join(runPath, "hooks", "post_gen_project.go"))
	if hasPostGenProjectHook {
		fmt.Println("Running post_gen_project...")
	}

	//
}
