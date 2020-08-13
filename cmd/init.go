package cmd

import (
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/tcnksm/go-input"
	Utils "jape/utils"
)

var initCmd = &cobra.Command{
	Use:   "init <flavor>",
	Short: "Initialize a new Jade project of the given flavor in the current directory",
	Long: `This command initializes a new Jade project for containerizing the code 
in the current directory. The flavor determines how the code is packaged. 
Valid flavors include: empty, python, java, fiji, matlab, bash
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		flavor := args[0]

		Utils.PrintInfo("Initializing Jape Project")

		if flavor == "fiji" {
			fiji()
		} else {
			msg := fmt.Errorf("invalid flavor specified: %s", flavor)
			Utils.PrintFatal("%s", msg)
			os.Exit(1)
		}

		Utils.PrintInfo("Jape project was successfully initialized.")
		Utils.PrintInfo("Next run `jape build` to build the container.")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

var ui = &input.UI{
	Writer: os.Stdout,
	Reader: os.Stdin,
}

type fijiParameters struct {
	PluginDir string
	MacroDir  string
	MacroName string
}

func fiji() {

	tpl := template.Must(template.ParseFiles("templates/fiji.got"))
	data := fijiParameters{"fiji_plugins", "fiji_macros", "macro.ijm"}

	data.PluginDir = read(data.PluginDir, "Relative path to the Fiji plugins which should be included in the container?")
	data.MacroDir = read(data.MacroDir, "Relative path to the Fiji macros which should be included in the container?")
	data.MacroName = read(data.MacroName, "Name of the Fiji macro file to execute when running the container?")

	mkdir(data.PluginDir)
	mkdir(data.MacroDir)

	macroPath := data.MacroDir + "/" + data.MacroName
	if Utils.FileExists(macroPath) {
		Utils.PrintSuccess("Found macro file at %s", macroPath)
	} else {
		Utils.PrintFatal("Could not find macro file at %s", macroPath)
	}

	dockerFilePath := "Dockerfile"
	if Utils.FileExists(dockerFilePath) {

		Utils.PrintError("Dockerfile already exists")

	} else {
		if f, err := os.OpenFile(dockerFilePath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0644); err == nil {
			defer f.Close()

			err2 := tpl.Execute(f, data)
			if err2 != nil {
				Utils.PrintFatal("Error creating Dockerfile: %s", err2)
			}

			Utils.PrintSuccess("Created Dockerfile")

		} else {
			Utils.PrintFatal("Error creating Dockerfile: %s", err)
		}
	}
}

func mkdir(dir string) {
	if Utils.Mkdir(dir) {
		Utils.PrintSuccess("Created directory %s", dir)
	} else {
		Utils.PrintSuccess("Found existing directory: %s", dir)
	}
}

func read(value string, query string) string {
	newValue, err := ui.Ask(query, &input.Options{
		Default:  value,
		Required: true,
		Loop:     true,
	})
	if err != nil {
		log.Fatalln(err)
	}
	return newValue
}
