package cmd

import (
	"github.com/spf13/cobra"
	Utils "jape/utils"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set [variable] [value]",
	Short: "Set the values of variables in the Jape configuration",
	Long: `Convenience command for quickly updating the project version and other variables 
without manually editing the jape.yaml. Can be useful for automation, too. 
Supported variables:
- version
- repo_tag
`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		varName := args[0]
		varValue := args[1]

		var config = Utils.ReadProjectConfig()

		switch varName {
		case "version":
			config.Version = varValue
		case "repo_tag":
			config.Config.Build.RepoTag = varValue
		default:
			Utils.PrintFatal("Unrecognized variable %s", varName)
		}

		Utils.WriteProjectConfig(config)
		Utils.PrintSuccess("Updated %s to %s", varName, varValue)
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
