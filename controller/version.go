package controller

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Version(cmd *cobra.Command, args []string) {
	version := viper.GetString("global.version")
	build := viper.GetString("global.build")
	fmt.Println(fmt.Sprintf("Scaffold version %s, build %s", version, build))
}
