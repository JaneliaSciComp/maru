package utils

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

// ConfFile is the maru configuration file for the project
const ConfFile = "maru.yaml"

// MaruConfig is the parsed configuration file in memory
type MaruConfig struct {
	MaruVersion string `yaml:"maru_version"`
	Name        string
	Version     string
	Remotes     []string          `yaml:"remotes,omitempty"`
	BuildArgs   map[string]string `yaml:"build_args,omitempty"`

	TemplateArgs struct {
		Flavor string

		Build struct {
			RepoUrl string `yaml:"repo_url"`
			Command string `yaml:"command,omitempty"`
		} `yaml:"build,omitempty"`

		Executable struct {
			RelativeExePath string `yaml:"exe_path"`
		} `yaml:"executable,omitempty"`

		PythonConda struct {
			PythonVersion      string `yaml:"python_version"`
			Dependencies       string `yaml:"dependencies"`
			RelativeScriptPath string `yaml:"script_path"`
		} `yaml:"python_conda,omitempty"`

		JavaMaven struct {
			JDKVersion string `yaml:"jdk_version"`
			MainClass  string `yaml:"main_class"`
		} `yaml:"java_maven,omitempty"`

		FijiMacro struct {
			PluginDir string `yaml:"plugin_dir"`
			MacroDir  string `yaml:"macro_dir"`
			MacroName string `yaml:"macro_name"`
		} `yaml:"fiji_macro,omitempty"`

		MatlabCompiled struct {
		} `yaml:"matlab_compiled,omitempty"`
	} `yaml:"template_args,omitempty"`
}

// NewMaruConfig is the constructor for a MaruConfig
func NewMaruConfig(name string, version string) *MaruConfig {
	c := &MaruConfig{}
	c.Name = name
	c.Version = version
	return c
}

// GetBuildArg returns the value of BuildArgs with the given key. Applies string interpolation to the value,
// e.g. $version becomes the value of Version.
func (c *MaruConfig) GetBuildArg(key string) string {
	return strings.Replace(c.BuildArgs[key], "$version", c.Version, 1)
}

// SetBuildArg sets the given key/value pair in BuildArgs
func (c *MaruConfig) SetBuildArg(key string, value string) {
	if c.BuildArgs == nil {
		c.BuildArgs = make(map[string]string)
	}
	c.BuildArgs[key] = value
}

// GetRepoTag returns the value of GIT_TAG in BuildArgs, after applying string interpolation.
func (c *MaruConfig) GetRepoTag() string {
	return c.GetBuildArg("GIT_TAG")
}

// GetVersion returns the value of Version, after applying string interpolation,
// e.g. $git_tag becomes the value of GIT_TAG in BuildArgs.
func (c *MaruConfig) GetVersion() string {
	return strings.Replace(c.Version, "$git_tag", c.BuildArgs["GIT_TAG"], 1)
}

// GetNameVersion returns the versioned name of the container, e.g. name:version
func (c *MaruConfig) GetNameVersion() string {
	return c.Name + ":" + c.GetVersion()
}

// GetNameLatest returns the name of the container tagged with latest, e.g. name:latest
func (c *MaruConfig) GetNameLatest() string {
	return c.Name + ":latest"
}

// GetDockerTag returns the namespaced tag for the given remote, e.g. remote/name:version
func (c *MaruConfig) GetDockerTag(remote string) string {
	return remote + "/" + c.GetNameVersion()
}

// GetBuildCommand returns the command to use for building the code, prepended with line continuation
func (c *MaruConfig) GetBuildCommand() string {
	if c.TemplateArgs.Build.Command == "" {
		return ""
	}
	return "\\\n    && " + c.TemplateArgs.Build.Command
}

// HasRemotes returns true if the Remotes array is not empty
func (c *MaruConfig) HasRemotes() bool {
	return c.Remotes != nil && len(c.Remotes) > 0
}

// GetTemplateArgsChecksum calculates a checksum for the current values stored in the TemplateArgs
func (c *MaruConfig) GetTemplateArgsChecksum() string {
	h := sha256.New()
	s := fmt.Sprintf("%v", c.TemplateArgs)
	h.Write([]byte(s))
	sum := h.Sum(nil)
	return fmt.Sprintf("%x", sum)
}

// WriteProjectConfig writes the given project configuration to the working directory
func WriteProjectConfig(c *MaruConfig) {

	// Always overwrite the Maru version with the current version
	c.MaruVersion = MaruVersion

	raw, err := yaml.Marshal(&c)
	if err != nil {
		PrintFatal("Error creating project config: %s", err)
	}

	PrintDebug("Writing to %s...", ConfFile)
	err = ioutil.WriteFile(ConfFile, raw, 0644)
	if err != nil {
		PrintFatal("Error writing project config file: %s", err)
	}
}

// ReadProjectConfig reads the current project configuration from the working directory. Returns nil if no file exists.
func ReadProjectConfig() *MaruConfig {

	PrintDebug("Checking for %s...", ConfFile)
	if !FileExists(ConfFile) {
		return nil
	}

	var c = &MaruConfig{}

	PrintDebug("Reading from %s...", ConfFile)
	yamlFile, err := ioutil.ReadFile(ConfFile)
	if err != nil {
		PrintFatal("Error reading config file: %s", err)
	}

	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		PrintFatal("Error reading config file: %s", err)
	}

	return c
}

// ReadMandatoryProjectConfig reads the current project configuration from the working directory.
// Prints an errror message and quits if no configuration exists in the working directory.
func ReadMandatoryProjectConfig() *MaruConfig {
	var config = ReadProjectConfig()
	if config == nil {
		PrintFatal("Current directory does not contain a Maru project configuration")
	}
	return config
}
