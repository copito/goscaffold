package command

import (
	"github.com/copito/goscaffold/controller"
	"github.com/spf13/cobra"
)

var VerisonCmd = &cobra.Command{
	Use:   "version",
	Short: "version of the scaffold cli application",
	Long:  `version of the scaffold cli application`,
	Run:   controller.Version,
}
