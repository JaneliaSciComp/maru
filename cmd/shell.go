package cmd

import (
	Utils "maru/utils"
	"strings"

	"github.com/spf13/cobra"
)

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Starts a Bash shell into the current container",
	Long:  `Starts a Bash shell into the current container. Mainly used for debugging.`,
	Run: func(cmd *cobra.Command, args []string) {
		RunInteractive([]string{"/bin/bash"}, nil)
	},
}

func init() {
	rootCmd.AddCommand(shellCmd)
}

// A lot of this code was adapted from:
// https://stackoverflow.com/questions/58732588/accept-user-input-os-stdin-to-container-using-golang-docker-sdk-interactive-co
func RunInteractive(entrypoint []string, args []string) {

	config := Utils.ReadMandatoryProjectConfig()
	versionTag := config.GetNameVersion()
	Utils.PrintInfo("Creating interactive shell for %s", versionTag)

	cmdArgs := make([]string, 5+2*len(EnvParam))

	cmdArgs[0] = "run"
	cmdArgs[1] = "-it"

	i := 2
	if EnvParam != nil {
		for i, v := range EnvParam {
			cmdArgs[i] = "-e"
			cmdArgs[i+1] = v
			i += 2
		}
	}

	cmdArgs[i] = "--entrypoint"
	cmdArgs[i+1] = "/bin/bash"
	cmdArgs[i+2] = versionTag

	Utils.PrintHint("%% docker %s", strings.Join(cmdArgs, " "))

	err := Utils.RunCommand("docker", cmdArgs...)
	if err != nil {
		Utils.PrintError("Command `docker run` exited with %s", err)
	}

	/*
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			Utils.PrintFatal("%s", err)
		}
		defer cli.Close()

		ctx := context.Background()
		resp, err := cli.ContainerCreate(ctx, &container.Config{
			Image:        versionTag,
			Cmd:          args,
			Tty:          true,
			OpenStdin:    true,
			AttachStdin:  true,
			AttachStdout: true,
			StdinOnce:    true,
			Entrypoint:   entrypoint,
			Env:          EnvParam,
			User:         UserParam,
		}, nil, nil, nil, "")
		if err != nil {
			Utils.PrintFatal("%s", err)
		}

		// It's important to attach before starting, otherwise we'll miss the first prompt
		waiter, err := cli.ContainerAttach(ctx, resp.ID, types.ContainerAttachOptions{
			Stdout:       true,
			Stdin:        true,
			Stream:       true,
		})
		if err != nil {
			Utils.PrintFatal("%s", err)
		}

		// Set up IO pipes to run in the background
		go io.Copy(os.Stdout, waiter.Reader)
		go io.Copy(os.Stderr, waiter.Reader)
		go io.Copy(waiter.Conn, os.Stdin)

		// Ensure the terminal is raw
		fd := int(os.Stdin.Fd())
		var oldState *terminal.State
		if terminal.IsTerminal(fd) {
			oldState, err = terminal.MakeRaw(fd)
			if err != nil {
				Utils.PrintError("%s", err)
			}
			defer terminal.Restore(fd, oldState)
		}

		// Start the container process
		if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			Utils.PrintFatal("%s", err)
		}

		// Wait until the container is done
		statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
		select {
		case err := <-errCh:
			if err != nil {
				Utils.PrintFatal("%s", err)
			}
		case <-statusCh:
		}
	*/
}
