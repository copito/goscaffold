package command

import (
	"github.com/copito/goscaffold/controller"
	"github.com/spf13/cobra"
)

var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs the builder on a scaffold project",
	Long:  `Runs the builder on a scaffold project`,
	Run:   controller.Run,
}

func init() {
	RunCmd.Flags().IntP("number", "n", 10, "A help for number")
	// RunCmd.PersistentFlags().StringVar(&developer, "developer", "Unknown Developer!", "Developer name.")
}
