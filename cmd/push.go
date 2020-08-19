package cmd

import (
	//"context"
	//"encoding/base64"
	//"encoding/json"
	//"github.com/docker/docker/pkg/jsonmessage"
	//"github.com/docker/docker/api/types"
	//"github.com/docker/docker/client"
	//"github.com/moby/term"
	"github.com/spf13/cobra"
	Utils "maru/utils"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push the built container to all its registered namespaces",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {

		config := Utils.ReadProjectConfig()
		if !config.HasNamespaces() {
			Utils.PrintMessage("There are no namespaces configured for the current project.")
			Utils.PrintInfo("Use `maru namespace add` to add a new namespace, and `maru namespace rm` to delete one.")
		} else {
			imageName := config.GetNameVersion()
			Utils.PrintInfo("Pushing %s to %d repositories", imageName, len(config.Namespaces))

			for _, n := range config.Namespaces {
				registryTag := config.GetDockerTag(n)

				Utils.PrintMessage("Tagging with %s", registryTag)
				err := Utils.RunCommand("docker", "tag", imageName, registryTag)
				if err != nil {
					Utils.PrintError("Command `docker tag` failed with %s", err)
				} else {
					Utils.PrintMessage("Pushing to %s", registryTag)
					err := Utils.RunCommand("docker", "push", registryTag)
					if err != nil {
						Utils.PrintError("Command `docker push` failed with %s", err)
					} else {
						Utils.PrintSuccess("Successfully pushed to %s", registryTag)
					}
				}

				/*
				cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
				if err != nil {
					Utils.PrintFatal("%s", err)
				}
				defer cli.Close()

				ctx := context.Background()

				// Tag the image to include the registry and namespace
				err = cli.ImageTag(ctx, imageName, registryTag)
				if err != nil {
					Utils.PrintFatal("%s", err)
				}

				authConfig := types.AuthConfig{
					Username: username,
					Password: password,
				}
				encodedJSON, err := json.Marshal(authConfig)
				if err != nil {
					panic(err)
				}
				authStr := base64.URLEncoding.EncodeToString(encodedJSON)

				reader, err := cli.ImagePush(ctx, registryTag, types.ImagePushOptions{
					RegistryAuth: authStr,
				})
				if err != nil {
					Utils.PrintFatal("%s", err)
				}

				termFd, isTerm := term.GetFdInfo(os.Stderr)
				defer reader.Close()
				err2 := jsonmessage.DisplayJSONMessagesStream(reader, os.Stderr, termFd, isTerm, nil)
				if err2 != nil {
					Utils.PrintFatal("%s", err2)
				}
				*/
			}

		}
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
