package command

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/copito/goscaffold/setup"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "scaffold",
	Short: "scaffold - a scaffolding tool to generate cookiecuttered tools",
	Long: `Scaffold is a CLI tool for Go that empowers application developers generate projects quickly.
   
One can use scaffold to generate projects based on YAML configurations`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Check if color should be used
		useNoColor, _ := cmd.Flags().GetBool("no-color")
		useVerbose, _ := cmd.Flags().GetBool("verbose")
		isDryRun, _ := cmd.Flags().GetBool("dry-run")

		// Setup Logging
		lvl := new(slog.LevelVar)
		lvl.Set(slog.LevelWarn)
		if useVerbose {
			lvl.Set(slog.LevelInfo)
		}
		if isDryRun {
			lvl.Set(slog.LevelDebug)
		}
		logger := setup.SetupLogging(!useNoColor, lvl.Level())

		// Setup all configuration for this application
		setup.SetupConfig(logger)

		// Create a new context with the logger attached
		ctx := context.Background()
		ctx = context.WithValue(ctx, "logger", logger)
		ctx = context.WithValue(ctx, "dry_run", isDryRun)
		cmd.SetContext(ctx)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}

func init() {
	// set config defaults
	viper.SetDefault("garbage-collect", false)

	// persistent flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "activate verbose mode for more details")
	rootCmd.PersistentFlags().BoolP("no-color", "n", false, "disable colors on logs/debugs")
	rootCmd.PersistentFlags().BoolP("dry-run", "d", false, "allows to test command without side-effects")

	// connect to viper
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("no-color", rootCmd.PersistentFlags().Lookup("no-color"))
	viper.BindPFlag("dry-run", rootCmd.PersistentFlags().Lookup("dry-run"))

	// local flags;
	// RunCmd.Flags().StringVar(&config, "config", "", "/path/to/config.yml")
	// RunCmd.Flags().BoolVar(&daemon, "server", false, "Run hoarder as a server")
	// RunCmd.Flags().BoolVarP(&version, "version", "v", false, "Display the current version of this CLI")

	rootCmd.AddCommand(VerisonCmd)
	rootCmd.AddCommand(RunCmd)
	// rootCmd.AddCommand(initCmd)
}
