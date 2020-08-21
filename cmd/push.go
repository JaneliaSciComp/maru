package cmd

import (
	"github.com/spf13/cobra"
	Utils "maru/utils"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push the container to all its configured remotes",
	Long: `Deploys the container to all of its configured remotes. The container must be already built using the build command. Use remote command to list remotes or add a new one.`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		config := Utils.ReadMandatoryProjectConfig()
		if !config.HasRemotes() {
			Utils.PrintMessage("There are no remotes configured for the current project.")
			Utils.PrintInfo("Use `maru remote add` to add a new remote.")
		} else {
			imageName := config.GetNameVersion()
			Utils.PrintInfo("Pushing %s to %d repositories", imageName, len(config.Remotes))

			for _, n := range config.Remotes {
				registryTag := config.GetDockerTag(n)

				Utils.PrintHint("%% docker tag %s %s", imageName, registryTag)
				err := Utils.RunCommand("docker", "tag", imageName, registryTag)
				if err != nil {
					Utils.PrintError("Command `docker tag` failed with %s", err)
				} else {
					Utils.PrintHint("%% docker push %s", registryTag)
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
