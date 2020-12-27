package cmd

import (
	Utils "maru/utils"
	"strings"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [args]",
	Short: "Run the container for the current project",
	Long: `Runs a docker container for the current Maru project, passing any arguments directly to the container's entrypoint. 
The current directory must contain a maru.yaml file describing the project. You can create a runnable project using the init and build commands. 
`,
	Run: func(cmd *cobra.Command, args []string) {
		RunContainer(nil, args)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	// Disable parsing because we want to pass through flags to the containerized application
	runCmd.DisableFlagParsing = true
}

func RunContainer(entrypoint []string, args []string) {

	config := Utils.ReadMandatoryProjectConfig()
	versionTag := config.GetNameVersion()
	Utils.PrintInfo("Running %s", versionTag)

	cmdArgs := make([]string, 3+2*len(EnvParam)+len(args))

	cmdArgs[0] = "run"
	cmdArgs[1] = "-i"

	i := 2
	if EnvParam != nil {
		for i, v := range EnvParam {
			cmdArgs[i] = "-e"
			cmdArgs[i+1] = v
			i += 2
		}
	}

	cmdArgs[i] = versionTag

	if args != nil {
		for i, v := range args {
			cmdArgs[i] = v
			i++
		}
	}

	Utils.PrintHint("%% docker %s", strings.Join(cmdArgs, " "))

	err := Utils.RunCommand("docker", cmdArgs...)
	if err != nil {
		Utils.PrintError("Command `docker run` exited with %s", err)
	}

	// 	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	// 	if err != nil {
	// 		Utils.PrintFatal("%s", err)
	// 	}
	// 	defer cli.Close()

	// 	ctx := context.Background()

	// 	// Create the container using the current project context and user arguments
	// 	resp, err := cli.ContainerCreate(ctx, &container.Config{
	// 		Image:      versionTag,
	// 		Cmd:        args,
	// 		Tty:        true,
	// 		Entrypoint: entrypoint,
	// 		Env:        EnvParam,
	// 		User:       UserParam,
	// 	}, nil, nil, nil, "")
	// 	if err != nil {
	// 		Utils.PrintFatal("%s", err)
	// 	}

	// 	// Run the container
	// 	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
	// 		Utils.PrintFatal("%s", err)
	// 	}

	// 	// Monitor until the container is finished
	// 	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	// 	select {
	// 	case err := <-errCh:
	// 		if err != nil {
	// 			Utils.PrintFatal("%s", err)
	// 		}
	// 	case <-statusCh:
	// 	}

	// 	// Copy the container logs so that the user can view them
	// 	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
	// 		ShowStdout: true,
	// 	})
	// 	if err != nil {
	// 		Utils.PrintFatal("%s", err)
	// 	}

	// 	if _, err := io.Copy(os.Stdout, out); err != nil {
	// 		Utils.PrintFatal("%s", err)
	// 	}
}

func GetEnvVariableString() string {
	envParams := make([]string, len(EnvParam)+1)
	if EnvParam != nil {
		for i, v := range EnvParam {
			envParams[i] = "-e " + v
		}
		envParams[len(EnvParam)] = ""
	}
	return strings.Join(envParams, " ")
}
