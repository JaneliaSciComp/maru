package cmd

import (
	"context"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/jhoonb/archivex"
	"github.com/moby/term"
	"github.com/spf13/cobra"
	"io/ioutil"
	Utils "jape/utils"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build a container for the current project",
	Long: `Runs a docker build for the current Jape project. The current directory must contain a jape.yaml 
file describing the project. You can initialize a project using the init command.
`,
	Run: func(cmd *cobra.Command, args []string) {
		runBuild()
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runBuild() {
	// To get the Docker client working, I had to `go get github.com/docker/docker@master`
	// as per https://github.com/moby/moby/issues/40185

	config := Utils.ReadProjectConfig()
	versionTag := config.Name+":"+config.Version

	Utils.PrintInfo("Building %s", versionTag)

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		Utils.PrintFatal("%s", err)
	}
	defer cli.Close()

	Utils.PrintMessage("Creating build context...")

	file, err := ioutil.TempFile("", "jape_build_ctx_*.tar")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())

	tar := new(archivex.TarFile)
	tar.Create(file.Name())
	tar.AddAll(".", true)
	tar.Close()
	dockerBuildContext, err := os.Open(file.Name())
	defer dockerBuildContext.Close()

	Utils.PrintMessage("Building image...")

	options := types.ImageBuildOptions{
		SuppressOutput: false,
		Tags:           []string{config.Name+":latest",versionTag},
		Dockerfile:     "./Dockerfile",
		//BuildArgs:      args,
	}
	r, err := cli.ImageBuild(context.Background(), dockerBuildContext, options)
	if err != nil {
		Utils.PrintFatal("%s", err)
	}

	termFd, isTerm := term.GetFdInfo(os.Stderr)
	defer r.Body.Close()
	err2 := jsonmessage.DisplayJSONMessagesStream(r.Body, os.Stderr, termFd, isTerm, nil)
	if err2 != nil {
		Utils.PrintFatal("%s", err2)
	}

	Utils.PrintSuccess("Successfully built %s", versionTag)
	Utils.PrintInfo("Next use `jape run` to run the container")
}
