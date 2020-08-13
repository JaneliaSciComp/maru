package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2/terminal"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	Utils "jape/utils"
)

var initCmd = &cobra.Command{
	Use:   "init <flavor>",
	Short: "Initialize a new Jade project of the given flavor in the current directory",
	Long: `This command initializes a new Jade project for containerizing the code 
in the current directory. The flavor determines how the code is packaged. 
Valid flavors include: fiji_macro, python_conda, java_maven, matlab_compiled
`,
	Run: func(cmd *cobra.Command, args []string) {

		Utils.PrintInfo("Initialize Jape Project")

		const fijiMacro = "fiji_macro"
		const pythonConda = "python_conda"
		const javaMaven = "java_maven"
		const matlabCompiled = "matlab_compiled"

		flavor := pythonConda
		if len(args)==0 {
			prompt := &survey.Select{
				Message: "What flavor of scientific software do you want to containerize?",
				Options: []string{
					pythonConda+" - Python project packaged with Conda",
					javaMaven+" - Java project with Maven build",
					fijiMacro+" - Fiji macro",
					matlabCompiled+" - MATLAB script, compiled",
				},
				Default: flavor,
			}
			ask(prompt, &flavor)
		} else {
			flavor = args[0]
		}

		if strings.HasPrefix(flavor, pythonConda) {
			initProjectPython()
		} else if strings.HasPrefix(flavor, javaMaven) {
			initProjectJavaMaven()
		} else if strings.HasPrefix(flavor, fijiMacro) {
			initProjectFiji()
		} else {
			Utils.PrintFatal("Flavor is currently not supported: %s", flavor)
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

func initProjectFiji() {

	pluginDir := askForString("Relative path to Fiji plugins:", "fiji_plugins")
	createDirectory(pluginDir)

	macroDir := askForString("Relative path to Fiji macros:", "fiji_macros")
	createDirectory(macroDir)

	macroName := askForString("Name of the Fiji macro file to run:", "macro.ijm")

	macroPath := macroDir + "/" + macroName
	if Utils.FileExists(macroPath) {
		Utils.PrintSuccess("Found macro file at %s", macroPath)
	} else {
		Utils.PrintFatal("Could not find macro file at %s", macroPath)
	}

	data := struct {
		PluginDir string
		MacroDir  string
		MacroName string
	}{
		pluginDir,
		macroDir,
		macroName,
	}
	generateDockerfile("templates/fiji.got", data)
}

func initProjectPython() {

	pythonVersion := "3.6"
	prompt := &survey.Select{
		Message: "Python version:",
		Options: []string{"2.7", "3.6", "3.7"},
		Default: pythonVersion,
	}
	ask(prompt, &pythonVersion)

	dependenciesText := ""
	mlPrompt := &survey.Multiline{
		Message: "Dependencies to install with Conda (e.g. h5py=2.8.0)",
	}
	ask(mlPrompt, &dependenciesText)

	re := regexp.MustCompile(`\n`)
	dependencies := re.ReplaceAllString(dependenciesText, " ")

	relativeScriptPath := askForString("Relative path to main script:", "main.py")

	data := struct {
		PythonVersion string
		Dependencies  string
		RelativeScriptPath string
	}{
		pythonVersion,
		dependencies,
		relativeScriptPath,
	}
	generateDockerfile("templates/python.got", data)
}

func initProjectJavaMaven() {

	gitUrl := askForString("URL to git repo:", "")
	buildCommand := askForString("Build command:", "mvn package")
	mainClass := askForString("Main class:", "org.janelia.app.MyClass")

	data := struct {
		GitUrl string
		BuildCommand  string
		MainClass string
	}{
		gitUrl,
		buildCommand,
		mainClass,
	}
	generateDockerfile("templates/java_maven.got", data)
}

func ask(prompt survey.Prompt, response interface{}, opts ...survey.AskOpt) {
	err := survey.AskOne(prompt, response)
	if err == terminal.InterruptErr {
		fmt.Println("interrupted")
		os.Exit(0)
	} else if err != nil {
		panic(err)
	}
}

func askForString(message string, defaultValue string) string {
	value := defaultValue
	prompt := &survey.Input{
		Message: message,
		Default: value,
	}
	ask(prompt, &value)
	return value
}

func createDirectory(dir string) {
	if Utils.Mkdir(dir) {
		Utils.PrintSuccess("Created directory %s", dir)
	} else {
		Utils.PrintSuccess("Found existing directory: %s", dir)
	}
}

func generateDockerfile(templateName string, data interface{}) {
	const dockerFilePath = "Dockerfile"

	if Utils.FileExists(dockerFilePath) {
		replace := false
		prompt := &survey.Confirm{
			Message: "Found existing Dockerfile. Replace?",
			Default: replace,
		}
		ask(prompt, &replace)
		if !replace {
			Utils.PrintFatal("Project initialization aborted")
		}
	}
	tmpl := template.Must(template.ParseFiles(templateName))
	if f, err := os.OpenFile(dockerFilePath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644); err == nil {
		defer f.Close()
		err2 := tmpl.Execute(f, data)
		if err2 != nil {
			Utils.PrintFatal("Error creating Dockerfile: %s", err2)
		}

		Utils.PrintSuccess("Created Dockerfile")

	} else {
		Utils.PrintFatal("Error creating Dockerfile: %s", err)
	}
}