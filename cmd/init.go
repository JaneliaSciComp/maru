package cmd

import (
	"context"
	"github.com/AlecAivazis/survey/v2"
	"github.com/posener/gitfs"
	"github.com/posener/gitfs/fsutil"
	"github.com/spf13/cobra"
	Utils "maru/utils"
	"net/url"
	"os"
	"path"
	"regexp"
	"sort"
)

// When running with `LOCAL_DEBUG=.`, the local repository will be used instead of the remote github.
var localDebug = os.Getenv("LOCAL_DEBUG")

type initFunctionType func (*Utils.MaruConfig, bool)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize or update a Maru project in the current directory",
	Long: `This command initializes or updates a Maru project in the current directory. If a Dockerfile already exists
in the current directory, it can either be used to bootstrap a custom project or overwritten. If a maru.yaml file
exists in the current directory, the initialization questionnaire will run again using the default values from the 
maru.yaml file. 
`,
	Run: func(cmd *cobra.Command, args []string) {
		config := Init()
		printFinalInstructions(config)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func Init() *Utils.MaruConfig {

	Utils.PrintInfo("Configure Maru Project")

	flavorMap := map[string]initFunctionType {
		"executable":      initProjectExecutable,
		"python_conda":    initProjectPython,
		"java_maven":      initProjectJavaMaven,
		"fiji_macro":      initProjectFiji,
		"matlab_compiled": initProjectMatlab,
	}

	var isNewProject = false

	var config = Utils.ReadProjectConfig()
	if config == nil {

		if Utils.FileExists(Utils.DockerFilePath) {
			if Utils.AskForBool("Create new Maru project using existing Dockerfile?", true) {
				containerName := Utils.AskForString("Container name:", "")
				containerVersion := Utils.AskForString("Container version:", "1.0.0")
				config = Utils.NewMaruConfig("custom", containerName, containerVersion)
				Utils.WriteProjectConfig(config)
				printFinalInstructions(config)
				os.Exit(0)
			}
		}

		isNewProject = true
		config = Utils.NewMaruConfig("", "", "1.0.0")
		config.Config.Build.RepoUrl = "https://github.com/example/repo.git"
		config.BuildArgs["GIT_TAG"] = "$version"
		config.Config.Build.Command = ""
	}

	flavors := make([]string, 0, len(flavorMap))
	for k := range flavorMap {
		flavors = append(flavors, k)
	}
	sort.Strings(flavors)

	flavor := config.Flavor
	prompt := &survey.Select{
		Message: "Flavor of container to build:",
		Options: flavors,
		Default: flavor,
	}
	Utils.Ask(prompt, &flavor)
	config.Flavor = flavor

	// Validate flavor before going further
	initFunction := flavorMap[flavor]
	if initFunction==nil {
		Utils.PrintFatal("Flavor is currently not supported: %s", flavor)
		os.Exit(1)
	}

	Utils.PrintInfo("\nWhich git repository should be built inside the container when ^maru build^ is called?")
	config.Config.Build.RepoUrl = Utils.AskForString("Git URL:", config.Config.Build.RepoUrl)

	Utils.PrintInfo("\nWhich tag or branch should be built when ^maru build^ is called?")
	Utils.PrintMessage(
		`You can use ^master^ to build the master branch, but that's not recommended for creating reproducible containers.
The best practice is to tag your git commit with a version number, and then use that same version as the container tag.
`)
	config.BuildArgs["GIT_TAG"] = Utils.AskForString("Git tag:", config.BuildArgs["GIT_TAG"])

	u, err := url.Parse(config.Config.Build.RepoUrl)
	if err != nil {
		Utils.PrintFatal("Problem parsing Git URL: %s",err)
	}

	if u.Scheme != "https" {
		Utils.PrintFatal("URL must begin with https")
	}

	if u.Host == "" {
		Utils.PrintFatal("URL must contain valid hostname")
	}

	// Default container name is the name of the current working directory
	if config.Name == "" {
		cwd, err := os.Getwd()
		if err != nil {
			Utils.PrintFatal("%s", err)
		}
		cwdName := path.Base(cwd)
		config.Name = cwdName
	}

	Utils.PrintInfo("\nWhat is the name for this container?")
	Utils.PrintMessage(
		`The container name should only contain lowercase letters and underscores. 
As an example, the container ^scientificlinux/sl:7^ is composed of the following elements:
    remote/namespace:    scientificlinux
    container name:      sl
    container version:   7
`)
	config.Name = Utils.AskForString("Container name:", config.Name)

	Utils.PrintInfo("\nWhat is the current version for this container?")
	Utils.PrintMessage(`This will change over time, and can be easily updated with ^maru set version^.
`)
	config.Version = Utils.AskForString("Container version:", config.Version)

	// Invoke the init function for the chosen project flavor
	initFunction(config, isNewProject)

	Utils.WriteProjectConfig(config)
	Utils.PrintSuccess("Created %s", Utils.ConfFile)

	// Replace empty build with a no-op so that bash script still works
	if config.Config.Build.Command=="" {
		config.Config.Build.Command = "true"
	}

	// Set the checksum for printing into the template
	config.ConfigChecksum = config.GetConfigChecksum()

	generateDockerfile(config)
	return config
}

func initProjectExecutable(config *Utils.MaruConfig, isNewProject bool) {

	pc := &config.Config.Executable

	if isNewProject {
		// Default values
		config.Config.Build.Command = "make"
		pc.RelativeExePath = "bin/program"
	}

	config.Config.Build.Command = Utils.AskForString("Build command:", config.Config.Build.Command)
	pc.RelativeExePath = Utils.AskForString("Relative path to built executable:", pc.RelativeExePath)
}

func initProjectFiji(config *Utils.MaruConfig, isNewProject bool) {

	pc := &config.Config.FijiMacro

	if isNewProject {
		// Default values
		pc.MacroDir = "fiji_plugins"
		pc.MacroDir = "fiji_macros"
		pc.MacroName = "macro.ijm"
	}

	pc.PluginDir = Utils.AskForString("Relative path to Fiji plugins:", pc.PluginDir)
	pc.MacroDir = Utils.AskForString("Relative path to Fiji macros:", pc.MacroDir)
	pc.MacroName = Utils.AskForString("Name of the Fiji macro file to run:", pc.MacroName)
}

func initProjectPython(config *Utils.MaruConfig, isNewProject bool) {

	pc := &config.Config.PythonConda

	if isNewProject {
		// Default values
		pc.PythonVersion = "3.6"
		pc.Dependencies = ""
		pc.RelativeScriptPath = "main.py"
	}

	prompt := &survey.Select{
		Message: "Python version:",
		Options: []string{"2.7", "3.6", "3.7"},
		Default: pc.PythonVersion,
	}
	Utils.Ask(prompt, &pc.PythonVersion)

	dependenciesText := pc.Dependencies
	mlPrompt := &survey.Multiline{
		Message: "Dependencies to install with Conda (e.g. h5py=2.8.0)",
		Default: dependenciesText,
	}
	Utils.Ask(mlPrompt, &dependenciesText)
	pc.Dependencies = regexp.MustCompile(`\s+`).ReplaceAllString(dependenciesText, " ")

	pc.RelativeScriptPath = Utils.AskForString("Relative path to main script:", pc.RelativeScriptPath)
}

func initProjectJavaMaven(config *Utils.MaruConfig, isNewProject bool) {

	pc := &config.Config.JavaMaven

	if isNewProject {
		// Default values
		config.Config.Build.Command = "mvn package"
		pc.MainClass = "org.myapp.MyClass"
	}

	config.Config.Build.Command = Utils.AskForString("Build command:", config.Config.Build.Command)
	pc.MainClass = Utils.AskForString("Main class:", pc.MainClass)
}

func initProjectMatlab(config *Utils.MaruConfig, isNewProject bool) {
	Utils.PrintFatal("MATLAB is currently unsupported")
}

func generateDockerfile(config *Utils.MaruConfig) {

	if Utils.FileExists(Utils.DockerFilePath) {
		if !Utils.AskForBool("Found existing Dockerfile. Replace?", true) {
			Utils.PrintFatal("Project initialization aborted")
		}
	}

	templateName := config.Flavor+".got"

	fs, err := gitfs.New(context.Background(), "github.com/JaneliaSciComp/maru/templates", gitfs.OptLocal(localDebug))
	if err != nil {
		Utils.PrintFatal("Failed creating gitfs: %s", err)
	}

	tmpls, err := fsutil.TmplParse(fs, nil, "/"+templateName)
	if err != nil {
		Utils.PrintFatal("Failed parsing templates: %s", err)
	}

	if f, err := os.OpenFile(Utils.DockerFilePath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644); err == nil {
		defer f.Close()

		err2 := tmpls.ExecuteTemplate(f, templateName, config)
		if err2 != nil {
			Utils.PrintFatal("Failed to create Dockerfile: %s", err2)
		}

		Utils.PrintSuccess("Created Dockerfile")

	} else {
		Utils.PrintFatal("Failed to create Dockerfile: %s", err)
	}
}

func printFinalInstructions(config *Utils.MaruConfig) {
	Utils.PrintSuccess("Maru project %s was successfully initialized.", config.GetNameVersion())
	Utils.PrintInfo("You can edit the maru.yaml file any time to update the project configuration.")
	Utils.PrintInfo("Next run `maru build` to build and tag the container.")
}