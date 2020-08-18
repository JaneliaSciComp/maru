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
		Utils.PrintMessage("  Flavor: %s", config.Flavor)
		Utils.PrintMessage("  Namespace Tags:")
		if config.HasNamespaces() {
			for _, n := range config.Namespaces {
				Utils.PrintMessage("  - %s", config.GetDockerTag(n))
			}
		} else {
			Utils.PrintMessage("  - %s", config.GetNameVersion())
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
