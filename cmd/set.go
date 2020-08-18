package cmd

import (
	"github.com/spf13/cobra"
	Utils "maru/utils"
)

var setCmd = &cobra.Command{
	Use:   "set [variable] [value]",
	Short: "Set the values of variables in the Maru configuration",
	Long: `Convenience command for quickly updating the project version and other variables without 
manually editing the maru.yaml.`,
}

var setVersionCmd = &cobra.Command{
	Use:   "version [value]",
	Short: "Set the version of the container",
	Long: `Convenience command for quickly updating the current project's version.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		varValue := args[0]
		var config = Utils.ReadMandatoryProjectConfig()
		config.Version = varValue
		Utils.WriteProjectConfig(config)
		Utils.PrintSuccess("Updated version to %s", varValue)
	},
}

var setGitTagCmd = &cobra.Command{
	Use:   "git_tag [value]",
	Short: "Set the GIT tag to build",
	Long: `Convenience command for quickly updating the GIT tag to build.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		varValue := args[0]
		var config = Utils.ReadMandatoryProjectConfig()
		config.BuildArgs["GIT_TAG"] = varValue
		Utils.WriteProjectConfig(config)
		Utils.PrintSuccess("Updated git_tag to %s", varValue)
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
	setCmd.AddCommand(setVersionCmd)
	setCmd.AddCommand(setGitTagCmd)
}
