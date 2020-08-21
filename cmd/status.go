package cmd

import (
	"github.com/spf13/cobra"
	Utils "maru/utils"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the status of the current project",
	Long: `Prints information about the project in the current working directory.`,
	Run: func(cmd *cobra.Command, args []string) {

		var config = Utils.ReadMandatoryProjectConfig()
		Utils.PrintInfo("%s %s", config.Name, config.Version)
		Utils.PrintMessage("flavor: %s", config.Flavor)
		Utils.PrintMessage("local tags:")
		Utils.PrintMessage("- %s", config.GetNameLatest())
		Utils.PrintMessage("- %s", config.GetNameVersion())
		if config.HasRemotes() {
			Utils.PrintMessage("remote tags:")
			for _, n := range config.Remotes {
				Utils.PrintMessage("- %s", config.GetDockerTag(n))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
