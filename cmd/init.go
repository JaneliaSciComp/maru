package cmd

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/posener/gitfs"
	"github.com/posener/gitfs/fsutil"
	"github.com/spf13/cobra"
	Utils "jape/utils"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// Add debug mode environment variable. When running with `LOCAL_DEBUG=.`, the
// local git repository will be used instead of the remote github.
var localDebug = os.Getenv("LOCAL_DEBUG")

var initCmd = &cobra.Command{
	Use:   "init <flavor>",
	Short: "Initialize a new Jade project of the given flavor in the current directory",
	Long: `This command initializes a new Jade project for containerizing the code 
in the current directory. The flavor determines how the code is packaged. 
Valid flavors include: fiji_macro, python_conda, java_maven, matlab_compiled
`,
	Run: func(cmd *cobra.Command, args []string) {

		Utils.PrintInfo("Configure Jape Project")

		const fijiMacro = "fiji_macro"
		const pythonConda = "python_conda"
		const javaMaven = "java_maven"
		const matlabCompiled = "matlab_compiled"
		var isNewProject = false

		var config = Utils.ReadProjectConfig()
		if config == nil {
			isNewProject = true
			config = Utils.NewJapeConfig(fijiMacro, "")
		}

		flavor := config.Flavor
		if len(args)==0 {
			prompt := &survey.Select{
				Message: "What flavor of scientific software do you want to containerize?",
				Options: []string{
					pythonConda,
					javaMaven,
					fijiMacro,
					matlabCompiled,
				},
				Default: flavor,
			}
			ask(prompt, &flavor)
		} else {
			flavor = args[0]
		}

		config.Flavor = flavor
		config.Config.Repository.Url = askForString("Git URL:", config.Config.Repository.Url)
		config.Config.Repository.Tag = askForString("Git tag:", config.Config.Repository.Tag)

		u, err := url.Parse(config.Config.Repository.Url)
		if err != nil {
			Utils.PrintFatal("Problem parsing Git URL: %s",err)
		}

		if u.Scheme != "https" {
			Utils.PrintFatal("URL must begin with https")
		}

		if u.Host == "" {
			Utils.PrintFatal("URL must contain valid hostname")
		}

		if config.Name=="" {
			basename := path.Base(u.Path)
			config.Name = strings.ToLower(strings.TrimSuffix(basename, filepath.Ext(basename)))
		}

		config.Name = askForString("Container name:", config.Name)

		if strings.HasPrefix(flavor, pythonConda) {
			initProjectPython(config, isNewProject)

		} else if strings.HasPrefix(flavor, javaMaven) {
			initProjectJavaMaven(config, isNewProject)

		} else if strings.HasPrefix(flavor, fijiMacro) {
			initProjectFiji(config, isNewProject)

		} else {
			Utils.PrintFatal("Flavor is currently not supported: %s", flavor)
			os.Exit(1)
		}

		Utils.WriteProjectConfig(config)

		Utils.PrintInfo("Jape project was successfully initialized.")
		Utils.PrintInfo("You can edit the jape.yaml file any time to update the project configuration.")
		Utils.PrintInfo("Next run `jape build` to build and tag the container.")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func initProjectFiji(config *Utils.JapeConfig, isNewProject bool) {

	pc := &config.Config.FijiMacro

	if isNewProject {
		pc.MacroDir = "fiji_plugins"
		pc.MacroDir = "fiji_macros"
		pc.MacroName = "macro"
	}

	pc.PluginDir = askForString("Relative path to Fiji plugins:", pc.PluginDir)
	createDirectory(config.Config.FijiMacro.PluginDir)

	pc.MacroDir = askForString("Relative path to Fiji macros:", pc.MacroDir)
	createDirectory(config.Config.FijiMacro.MacroDir)

	pc.MacroName = askForString("Name of the Fiji macro file to run:", pc.MacroName)

	macroPath := pc.MacroDir + "/" + pc.MacroName
	if Utils.FileExists(macroPath) {
		Utils.PrintSuccess("Found macro file at %s", macroPath)
	} else {
		Utils.PrintFatal("Could not find macro file at %s", macroPath)
	}

	generateDockerfile("fiji.got", config)
}

func initProjectPython(config *Utils.JapeConfig, isNewProject bool) {

	pc := &config.Config.PythonConda

	if isNewProject {
		pc.PythonVersion = "3.6"
		pc.Dependencies = ""
		pc.RelativeScriptPath = "main.py"
	}

	prompt := &survey.Select{
		Message: "Python version:",
		Options: []string{"2.7", "3.6", "3.7"},
		Default: pc.PythonVersion,
	}
	ask(prompt, &pc.PythonVersion)

	dependenciesText := pc.Dependencies
	mlPrompt := &survey.Multiline{
		Message: "Dependencies to install with Conda (e.g. h5py=2.8.0)",
		Default: dependenciesText,
	}
	ask(mlPrompt, &dependenciesText)
	pc.Dependencies = regexp.MustCompile(`\s+`).ReplaceAllString(dependenciesText, " ")

	pc.RelativeScriptPath = askForString("Relative path to main script:", pc.RelativeScriptPath)

	generateDockerfile("python.got", config)
}

func initProjectJavaMaven(config *Utils.JapeConfig, isNewProject bool) {

	pc := &config.Config.JavaMaven

	if isNewProject {
		config.Config.Build.Command = "mvn package"
		pc.MainClass = "org.myapp.MyClass"
	}

	config.Config.Build.Command = askForString("Build command:", config.Config.Build.Command)
	pc.MainClass = askForString("Main class:", pc.MainClass)

	generateDockerfile("java_maven.got", config)
}

func ask(prompt survey.Prompt, response interface{}, opts ...survey.AskOpt) {
	err := survey.AskOne(prompt, response)
	if err == terminal.InterruptErr {
		fmt.Println("interrupted")
		os.Exit(0)
	} else if err != nil {
		Utils.PrintFatal("%s", err)
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
		replace := true
		prompt := &survey.Confirm{
			Message: "Found existing Dockerfile. Replace?",
			Default: replace,
		}
		ask(prompt, &replace)
		if !replace {
			Utils.PrintFatal("Project initialization aborted")
		}
	}

	fs, err := gitfs.New(context.Background(),
		"github.com/JaneliaSciComp/jape/templates", gitfs.OptLocal(localDebug))
	if err != nil {
		Utils.PrintFatal("Failed creating gitfs: %s", err)
	}

	tmpls, err := fsutil.TmplParse(fs, nil, "/"+templateName)
	if err != nil {
		Utils.PrintFatal("Failed parsing templates: %s", err)
	}

	if f, err := os.OpenFile(dockerFilePath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644); err == nil {
		defer f.Close()

		err2 := tmpls.ExecuteTemplate(f, templateName, data)
		if err2 != nil {
			Utils.PrintFatal("Failed to create Dockerfile: %s", err2)
		}

		Utils.PrintSuccess("Created Dockerfile")

	} else {
		Utils.PrintFatal("Failed to create Dockerfile: %s", err)
	}
}