package cmd

import (
	"github.com/spf13/cobra"
	Utils "maru/utils"
)

var namespaceCmd = &cobra.Command{
	Use:   "namespace",
	Short: "Lists the namespaces configured for the current Maru project",
	Long: `Namespaces allow you to push your container to a remote registry and share it with others. If the
namespace does not begin with a hostname, then DockerHub is assumed.`,
	Run: func(cmd *cobra.Command, args []string) {

		var config = Utils.ReadMandatoryProjectConfig()
		if config.HasNamespaces() {
			Utils.PrintMessage("Namespaces configured for the current project: ")
			for i, n := range config.Namespaces {
				Utils.PrintMessage("%d) %s", i+1, n)
			}
		} else {
			Utils.PrintMessage("There are no namespaces configured for the current project.")
		}
		Utils.PrintInfo("Use `maru namespace add` to add a new namespace, and `maru namespace rm` to delete one.")

	},
}

var namespaceAddCmd = &cobra.Command{
	Use:   "add [namespace]",
	Short: "Add a namespace to the current project",
	Long: `Add the given namespace to the current project.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var namespace = args[0]
		var config = Utils.ReadMandatoryProjectConfig()
		config.Namespaces = append(config.Namespaces, namespace)
		Utils.WriteProjectConfig(config)
		Utils.PrintSuccess("Added namespace %s", namespace)
	},
}

var namespaceRmCmd = &cobra.Command{
	Use:   "rm [namespace]",
	Short: "Remove a namespace from the current project",
	Long: `Does not remove the container from any remote servers.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var namespace = args[0]
		var config = Utils.ReadMandatoryProjectConfig()
		if config.HasNamespaces() {
			idx := indexOf(namespace, config.Namespaces)
			if idx < 0 {
				Utils.PrintFatal("Namespace '%s' not found.", namespace)
			} else {
				config.Namespaces = remove(config.Namespaces, idx)
				Utils.WriteProjectConfig(config)
				Utils.PrintSuccess("Removed namespace %s", namespace)
			}
		} else {
			Utils.PrintFatal("There are no namespaces configured for the current project.")
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
	rootCmd.AddCommand(namespaceCmd)
	namespaceCmd.AddCommand(namespaceAddCmd)
	namespaceCmd.AddCommand(namespaceRmCmd)
}
