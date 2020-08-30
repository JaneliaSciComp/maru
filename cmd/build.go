package cmd

import (
	"context"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/jhoonb/archivex"
	"github.com/moby/term"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	Utils "maru/utils"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build container image for the current project",
	Long: `Runs a Docker build for the current Maru project. The current directory must contain a maru.yaml 
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

	config := Utils.ReadMandatoryProjectConfig()
	versionTag := config.GetNameVersion()

	if config.TemplateArgs.Flavor != "" {
		checksum := config.GetTemplateArgsChecksum()
		if !Utils.TestChecksum(checksum) {
			Utils.PrintDebug("Checksum does not match: %s", checksum)
			if Utils.AskForBool("The project configuration has changed. Do you want to regenerate the Dockerfile?", true) {
				Init()
				if !Utils.AskForBool("Proceed with container build?", true) {
					os.Exit(0)
				}
			}
		}
	}

	if config.TemplateArgs.Build.RepoUrl == "" {
		Utils.PrintInfo("Building %s", versionTag)
	} else {
		Utils.PrintInfo("Building %s from %s @ %s", versionTag,
			config.GetRepoTag(), config.TemplateArgs.Build.RepoUrl)
	}

	Utils.PrintHint("%% docker build . -t %s -t %s", config.GetNameLatest(), versionTag)

	// To get the Docker client working, I had to `go get github.com/docker/docker@master`
	// as per https://github.com/moby/moby/issues/40185

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
	Utils.PrintDebug("Creating temporary build context at %s", file.Name())

	tar := new(archivex.TarFile)
	tar.Create(file.Name())
	tar.AddAll(".", true)
	tar.Close()
	dockerBuildContext, err := os.Open(file.Name())
	defer dockerBuildContext.Close()

	Utils.PrintMessage("Building image...")

	// Interpolate map values to string pointers
	buildArgs := make(map[string]*string)
	for key, _ := range config.BuildArgs {
		v := config.GetBuildArg(key)
		buildArgs[key] = &v
	}

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
