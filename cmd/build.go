package cmd

import (
	Utils "maru/utils"
	"os"
	"strings"

	"github.com/spf13/cobra"
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

	Utils.PrintMessage("Building image...")

	// Interpolate map values to string pointers
	args := make([]string, 6+2*len(config.BuildArgs))

	args[0] = "build"

	i := 1
	for key := range config.BuildArgs {
		args[i] = "--build-arg"
		args[i+1] = key + "=" + config.GetBuildArg(key)
		i += 2
	}

	args[i] = "-t"
	args[i+1] = config.Name + ":latest"
	args[i+2] = "-t"
	args[i+3] = versionTag
	args[i+4] = "."

	Utils.PrintHint("%% docker %s", strings.Join(args, " "))

	err := Utils.RunCommand("docker", args...)
	if err != nil {
		Utils.PrintError("Command `docker build` failed with %s", err)
	} else {
		Utils.PrintSuccess("Successfully built %s", versionTag)
		Utils.PrintInfo("Next use `maru run` to run the container")
	}

	// To get the Docker client working, I had to `go get github.com/docker/docker@master`
	// as per https://github.com/moby/moby/issues/40185

	/*
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
	*/
}
