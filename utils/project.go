package utils

import (
	"crypto/sha256"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

const ConfFile = "maru.yaml"

type MaruConfig struct {

	ConfigChecksum string            `yaml:"-"`
	MaruVersion    string
	Name           string
	Version        string
	Remotes        []string          `yaml:"remotes,omitempty"`
	BuildArgs      map[string]string `yaml:"build_args,omitempty"`
	Flavor         string

	Config struct {
		Build                  struct {
			RepoUrl            string `yaml:"repo_url"`
			Command            string `yaml:"command,omitempty"`
		} `yaml:"build,omitempty"`

		Executable             struct {
			RelativeExePath    string `yaml:"exe_path"`
		} `yaml:"executable,omitempty"`

		PythonConda            struct {
			PythonVersion      string `yaml:"python_version"`
			Dependencies       string `yaml:"dependencies"`
			RelativeScriptPath string `yaml:"script_path"`
		} `yaml:"python_conda,omitempty"`

		JavaMaven              struct {
			MainClass          string `yaml:"main_class"`
		} `yaml:"java_maven,omitempty"`

		FijiMacro              struct {
			PluginDir          string `yaml:"plugin_dir"`
			MacroDir           string `yaml:"macro_dir"`
			MacroName          string `yaml:"macro_name"`
		} `yaml:"fiji_macro,omitempty"`

		MatlabCompiled struct {

		} `yaml:"matlab_compiled,omitempty"`

	} `yaml:"config,omitempty"`
}

func (c *MaruConfig) GetBuildArg(key string) string {
	// TODO: this needs a more general solution to variable interpolation
	s := strings.Replace(c.BuildArgs[key], "$version", c.Version, 1)
	return strings.Replace(s, "$git_tag", c.BuildArgs["GIT_TAG"], 1)
}

func (c *MaruConfig) GetRepoTag() string {
	return c.GetBuildArg("GIT_TAG")
}

func (c *MaruConfig) GetVersion() string {
	return strings.Replace(c.Version, "$git_tag", c.BuildArgs["GIT_TAG"], 1)
}

func (c *MaruConfig) GetNameVersion() string {
	return c.Name + ":" + c.GetVersion()
}

func (c *MaruConfig) GetNameLatest() string {
	return c.Name + ":latest"
}

func (c *MaruConfig) GetDockerTag(remote string) string {
	return remote + "/" + c.GetNameVersion()
}

func (c *MaruConfig) HasRemotes() bool {
	return c.Remotes != nil && len(c.Remotes)>0
}

func (c *MaruConfig) GetConfigChecksum() string {
	h := sha256.New()
	s := fmt.Sprintf("%v", c.Config)
	h.Write([]byte(s))
	sum := h.Sum(nil)
	return fmt.Sprintf("%x", sum)
}

func NewMaruConfig(flavor string, name string, version string) *MaruConfig {
	c := &MaruConfig{}
	c.Name = name
	c.Flavor = flavor
	c.Version = version
	return c
}

// Writes the given project configuration to the working directory
func WriteProjectConfig(c *MaruConfig) {

	// Always overwrite the Maru version with the current version
	c.MaruVersion = MaruVersion

	raw, err := yaml.Marshal(&c)
	if err != nil {
		PrintFatal("Error creating project config: %s", err)
	}

	err = ioutil.WriteFile(ConfFile, raw, 0644)
	if err != nil {
		PrintFatal("Error writing project config file: %s", err)
	}
}

// Reads the current project configuration from the working directory. Returns nil if no file exists.
func ReadProjectConfig() *MaruConfig {

	if !FileExists(ConfFile) {
		return nil
	}

	var c = &MaruConfig{}

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

func ReadMandatoryProjectConfig() *MaruConfig {
	var config = ReadProjectConfig()
	if config == nil {
		PrintFatal("Current directory does not contain a Maru project configuration")
	}
	return config
}