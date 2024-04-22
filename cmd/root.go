package command

import (
	"fmt"
	"os"

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
	RunCmd.PersistentFlags().StringP("config", "c", "./scaffold.yaml", "configuration file")
	RunCmd.PersistentFlags().StringP("token", "t", "", "Auth token used when connecting to a secure Hoarder")

	// connect to viper
	viper.BindPFlag("config", RunCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("token", RunCmd.PersistentFlags().Lookup("token"))

	// local flags;
	// RunCmd.Flags().StringVar(&config, "config", "", "/path/to/config.yml")
	// RunCmd.Flags().BoolVar(&daemon, "server", false, "Run hoarder as a server")
	// RunCmd.Flags().BoolVarP(&version, "version", "v", false, "Display the current version of this CLI")

	rootCmd.AddCommand(VerisonCmd)
	rootCmd.AddCommand(RunCmd)
	// rootCmd.AddCommand(initCmd)
}
