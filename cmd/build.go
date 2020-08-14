package cmd

import (
	"context"
	"fmt"
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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("build called")

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

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		Utils.PrintFatal("%s", err)
	}
	defer cli.Close()

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

	options := types.ImageBuildOptions{
		SuppressOutput: false,
		Tags:           []string{"testimage:latest"},
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
}
