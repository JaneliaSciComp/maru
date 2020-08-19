package cmd

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	Utils "maru/utils"
	"os"
	"strings"
)

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Starts a Bash shell into the current container",
	Long: `Starts a Bash shell into the current container. Mainly used for debugging.`,
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

	config := Utils.ReadProjectConfig()
	versionTag := config.GetNameVersion()
	Utils.PrintInfo("Creating interactive shell for %s", versionTag)

	Utils.PrintMessage("%% ^docker run -it %s--entrypoint=/bin/bash %s %s^",
		GetEnvVariableString(), versionTag, strings.Join(args, " "))

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
}