package cmd

import (
	Utils "maru/utils"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var buildArgs []string

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
	buildCmd.Flags().StringArrayVar(&buildArgs, "build-arg", nil, "Set build-time arguments for the container, e.g. when using run or shell")
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
	args := make([]string, 1)

	args[0] = "build"
	i := 1

	// Process command line build args first
	set := make(map[string]bool)
	for _, buildArg := range buildArgs {
		s := strings.Split(buildArg, "=")
		set[s[0]] = true
		args = append(args, "--build-arg")
		args = append(args, buildArg)
		i += 2
	}

	// Add any build args from the config file which were not overridden on the command line
	for key := range config.BuildArgs {
		if !set[key] {
			args = append(args, "--build-arg")
			args = append(args, key+"="+config.GetBuildArg(key))
		}
	}

	args = append(args, "-t")
	args = append(args, config.Name+":latest")
	args = append(args, "-t")
	args = append(args, versionTag)
	args = append(args, ".")

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
