package controller

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Version(cmd *cobra.Command, args []string) {
	version := viper.GetString("global.version")

	termCmd := exec.Command("git", "rev-parse", "HEAD")
	commit, err := termCmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Scaffold version %s, commit hash %s\n", version, commit)
}
