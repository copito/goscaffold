package controller

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/copito/goscaffold/core"
	"github.com/copito/goscaffold/entity"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rgxHooksFolder *regexp.Regexp = regexp.MustCompile(`\/hooks$`)
	rgxHooksFile   *regexp.Regexp = regexp.MustCompile(`hooks(\/|\\)(pre_prompt|pre_gen_project|post_gen_project)\\.(py|go|sh)$`)
)

func Run(cmd *cobra.Command, args []string) {
	// 1. Getting Path
	var runPath string
	if len(args) == 0 {
		fmt.Println("No path provided, assuming path is current path: .")
		runPath = "."
	} else {
		runPath = args[0]
	}

	// Check if path exists
	isExists, err := core.PathExists(runPath)
	if err != nil {
		fmt.Printf("Path provided does not exist...")
		os.Exit(1)
	}

	if !isExists {
		fmt.Printf("Path provided does not exist...")
		os.Exit(1)
	}

	// 2. Load config file
	configFilePath, err := cmd.Flags().GetString("config")
	if err != nil {
		configFilePath = "./config.yaml"
	}

	extension := path.Ext(configFilePath)[1:]
	dirConfig := path.Dir(configFilePath)
	basePath := core.FileNameWithoutExtension(path.Base(configFilePath))

	if basePath == "base" && extension == "yaml" {
		fmt.Println("base.yaml is the only name that cannot be used for the configuration file")
		os.Exit(1)
	}

	viper.SetConfigName(basePath)
	viper.SetConfigType(extension)
	viper.AddConfigPath(dirConfig)

	// Setting prefix for all env variables: SCAFFOLD_ID => viper.Get("ID")
	// viper.SetEnvPrefix("SCAFFOLD")

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println("Unable to find and/or load config file:", configFilePath)
		os.Exit(1)
	}

	// Parse configurations
	promptConfig := entity.Prompt{}
	err = viper.Unmarshal(&promptConfig)
	if err != nil {
		fmt.Println("Unable to load config file:", configFilePath)
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

	// TODO: ask questions about config (settle variables)
	// Loop through all prompt based configs (based on data type)
	for index, key := range keys {

		item := promptConfig.Items[key]
		item.Key = key
		fmt.Println("Running index: ", index, key)

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
			fmt.Printf("unexpected type %T", v)
			os.Exit(1)
		}
	}

	// TODO: send it to a file (if running under debug)
	fmt.Println("New Compiled Results: ", paramChoice)

	// 3. pre-hooks
	hasPreGenProjectHook, _ := core.PathExists(path.Join(runPath, "hooks", "pre_gen_project.go"))
	if hasPreGenProjectHook {
		fmt.Println("Running pre_gen_hook...")
	}

	// 4. Generate output folder to add output here
	outputFolder := "output"
	outputBasePath := path.Join(runPath, outputFolder)
	fmt.Println("Output Path =>", outputBasePath)
	err = os.Mkdir(outputBasePath, os.FileMode(0o755))
	if err != nil {
		fmt.Println("Could not create output folder")
		os.Exit(1)
	}

	// rollback output folder
	rollbackChan := make(chan interface{}, 1)
	rollbackOutput := func(rollbackChan chan interface{}, outputBasePath string) {
		defer close(rollbackChan)
		<-rollbackChan

		err := os.RemoveAll(outputBasePath)
		if err != nil {
			fmt.Println("Error cleaning up output folder...")
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
		fmt.Println(pathValue, info.Size(), info.Mode().IsDir(), info.Mode().IsRegular())
		deltaPath := core.DeltaRelativePath(runPath, pathValue)
		newFullPath := path.Join(outputBasePath, deltaPath)

		switch mode := info.Mode(); {
		case mode.IsDir():
			// Folder/Directory
			// TODO: Create folder
			err = os.MkdirAll(newFullPath, os.FileMode(0o755))
			if err != nil {
				rollbackChan <- true
			}

		case mode.IsRegular():
			// File
			bytesProcessed, err := core.PathCopy(pathValue, newFullPath)
			if err != nil {
				rollbackChan <- true
			}
			fmt.Print("Processed: ", bytesProcessed)
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
