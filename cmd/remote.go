package cmd

import (
	"github.com/spf13/cobra"
	Utils "maru/utils"
)

var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Lists the remote repositories configured for the current Maru project",
	Long: "Lists the remote repositories configured for the current Maru project. "+
		"Remotes allow you to push your container to a remote registry and share it with others. "+
		"If the remote does not begin with a hostname, then DockerHub is assumed.",
	Run: func(cmd *cobra.Command, args []string) {

		var config = Utils.ReadMandatoryProjectConfig()
		if config.HasRemotes() {
			Utils.PrintInfo("Remote repositories: ")
			for _, remote := range config.Remotes {
				Utils.PrintMessage("- %s", remote)
			}
		} else {
			Utils.PrintMessage("There are no remote repositories configured for the current project.")
		}
		Utils.PrintInfo("\nUse `maru remote add` to add a new remote, and `maru remote rm` to delete one.")

	},
}

var remoteAddCmd = &cobra.Command{
	Use:   "add [remote]",
	Short: "Add a remote to the current project",
	Long: `Add the given remote to the current project.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var remote = args[0]
		var config = Utils.ReadMandatoryProjectConfig()
		config.Remotes = append(config.Remotes, remote)
		Utils.WriteProjectConfig(config)
		Utils.PrintSuccess("Added remote %s", remote)
	},
}

var remoteRmCmd = &cobra.Command{
	Use:   "rm [remote]",
	Short: "Remove a remote from the current project",
	Long: `Does not remove the container from any remote servers.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var remote = args[0]
		var config = Utils.ReadMandatoryProjectConfig()
		if config.HasRemotes() {
			idx := indexOf(remote, config.Remotes)
			if idx < 0 {
				Utils.PrintFatal("Remote '%s' not found.", remote)
			} else {
				config.Remotes = remove(config.Remotes, idx)
				Utils.WriteProjectConfig(config)
				Utils.PrintSuccess("Removed remote %s", remote)
			}
		} else {
			Utils.PrintFatal("There are no remotes configured for the current project.")
		}
	},
}

func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1
}

func remove(slice []string, i int) []string {
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}

func init() {
	rootCmd.AddCommand(remoteCmd)
	remoteCmd.AddCommand(remoteAddCmd)
	remoteCmd.AddCommand(remoteRmCmd)
}
