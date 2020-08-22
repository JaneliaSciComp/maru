package cmd

import (
	"github.com/spf13/cobra"
	Utils "maru/utils"
	"os/exec"
)

var singularityCmd = &cobra.Command{
	Use:   "singularity",
	Short: "Run containers using Singularity",
	Long: "Run containers using Singularity. This is used mainly for running on HPC clusters.",
}

var singularityBuildCmd = &cobra.Command{
	Use:   "build [output image file]",
	Short: "Builds a Singularity container from the existing Docker container",
	Long: "Builds a Singularity container (in Singularity Image Format) from the built Docker container.\n"+
	"This assumes that `maru build` was already run successfully and the Docker container exists on disk.",
	Args: cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {

		if !isCommandAvailable("singularity") {
			Utils.PrintFatal("You need to install Singularity before using this command.")
		}

		var config = Utils.ReadMandatoryProjectConfig()
		imageName := config.GetNameVersion()

		// Default to temp directory
		outFile := "/tmp/"+config.Name+"_"+config.Version+".sif"
		if len(args) > 0 {
			outFile = args[0]
		}

		Utils.PrintHint("%% singularity build %s docker-daemon://%s", outFile, imageName)
		err := Utils.RunCommand("singularity", "build", outFile, "docker-daemon://"+imageName)
		if err != nil {
			Utils.PrintError("Singularity build failed: %s", err)
		}

		Utils.PrintSuccess("Singularity container saved to %s", outFile)
		Utils.PrintInfo("You can now run the container: ^singularity run %s^", outFile)
	},
}

var singularityRunCmd = &cobra.Command{
	Use:   "run [args]",
	Short: "Runs the current Maru project using Singularity",
	Long: "Runs the current Maru project using Singularity, passing any arguments to the container's entrypoint.\n"+
	"This first runs an implicit command equivalent to `maru singularity build` in order to to convert the container \n" +
	"to Singularity Image Format. Environment variables may be passed using the -e flag, but the user flag -u will \n" +
	"have no affect because Singularity always runs as the current user.",
	Run: func(cmd *cobra.Command, args []string) {

		if !isCommandAvailable("singularity") {
			Utils.PrintFatal("You need to install Singularity before using this command.")
		}

		var config = Utils.ReadMandatoryProjectConfig()
		imageName := config.GetNameVersion()

		// The usual slashes after docker-daemon are not accepted here
		// https://github.com/hpcng/singularity/issues/4734
		Utils.PrintHint("%% singularity run docker-daemon:%s", imageName)
		err := Utils.RunCommand("singularity", "run", "docker-daemon:"+imageName)
		if err != nil {
			Utils.PrintError("Singularity run failed: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(singularityCmd)
	singularityCmd.AddCommand(singularityBuildCmd)
	singularityCmd.AddCommand(singularityRunCmd)
}

// From https://siongui.github.io/2018/03/16/go-check-if-command-exists/
func isCommandAvailable(name string) bool {
	cmd := exec.Command("/bin/sh", "-c", "command -v "+name)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}