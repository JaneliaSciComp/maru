package cmd

import (
	"github.com/spf13/cobra"
	Utils "maru/utils"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the current version",
	Long: `Prints the current version of Maru'`,
	Run: func(cmd *cobra.Command, args []string) {
		Utils.PrintSuccess("Maru %s", Utils.MaruVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
