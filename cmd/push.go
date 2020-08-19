package cmd

import (
	"github.com/spf13/cobra"
	Utils "maru/utils"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push the container to all its configured namespaces",
	Long: `Deploys the container to all of its configured namespaces. The container must be already built using the build command. Use namespace command to list namespaces or add a new one.`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		config := Utils.ReadMandatoryProjectConfig()
		if !config.HasNamespaces() {
			Utils.PrintMessage("There are no namespaces configured for the current project.")
			Utils.PrintInfo("Use `maru namespace add` to add a new namespace.")
		} else {
			imageName := config.GetNameVersion()
			Utils.PrintInfo("Pushing %s to %d repositories", imageName, len(config.Namespaces))

			for _, n := range config.Namespaces {
				registryTag := config.GetDockerTag(n)

				Utils.PrintMessage("%% ^docker tag %s %s^", imageName, registryTag)
				err := Utils.RunCommand("docker", "tag", imageName, registryTag)
				if err != nil {
					Utils.PrintError("Command `docker tag` failed with %s", err)
				} else {
					Utils.PrintMessage("%% ^docker push %s^", registryTag)
					err := Utils.RunCommand("docker", "push", registryTag)
					if err != nil {
						Utils.PrintError("Command `docker push` failed with %s", err)
					} else {
						Utils.PrintSuccess("Successfully pushed to %s", registryTag)
					}
				}
			}

		}
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
