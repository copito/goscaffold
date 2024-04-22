package controller

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Version(cmd *cobra.Command, args []string) {
	version := viper.GetString("global.version")

	// Getting last commit version hash
	getCommitHash := func() string {
		if info, ok := debug.ReadBuildInfo(); ok {
			for _, setting := range info.Settings {
				if setting.Key == "vcs.revision" {
					return setting.Value
				}
			}
		}

		return ""
	}
	build := getCommitHash()
	fmt.Printf("Scaffold version %s, commit hash %s", version, build)
}
