package cmd

import (
	"context"
	"io/ioutil"
	Utils "maru/utils"
	"log"
	"os"

	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/jhoonb/archivex"
	"github.com/moby/term"
	"github.com/spf13/cobra"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build a container for the current project",
	Long: `Runs a docker build for the current Maru project. The current directory must contain a maru.yaml 
file describing the project. You can initialize a project using the init command.
`,
	Run: func(cmd *cobra.Command, args []string) {
		runBuild()
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}

func runBuild() {
	// To get the Docker client working, I had to `go get github.com/docker/docker@master`
	// as per https://github.com/moby/moby/issues/40185

	config := Utils.ReadProjectConfig()
	versionTag := config.Name + ":" + config.Version

	Utils.PrintInfo("Building %s from %s @ %s", versionTag,
		config.Config.Build.RepoTag, config.Config.Build.RepoUrl)

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		Utils.PrintFatal("%s", err)
	}
	defer cli.Close()

	Utils.PrintMessage("Creating build context...")

	file, err := ioutil.TempFile("", "maru_build_ctx_*.tar")
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

	buildArgs := make(map[string]*string)
	buildArgs["APP_TAG"] = &config.Config.Build.RepoTag

	options := types.ImageBuildOptions{
		SuppressOutput: false,
		Tags:           []string{config.Name + ":latest", versionTag},
		Dockerfile:     "./Dockerfile",
		BuildArgs:      buildArgs,
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
	Utils.PrintInfo("Next use `maru run` to run the container")
}
