package utils

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const ConfFile = "jape.yaml"

type JapeConfig struct {
	Name    string
	Version string
	Flavor  string
	Config  struct {
		Build                  struct {
			RepoUrl            string `yaml:"repo_url"`
			RepoTag            string `yaml:"repo_tag"`
			Command            string
		} `yaml:"build,omitempty"`

		Executable             struct {
			RelativeExePath    string `yaml:"relative_exe_path"`
		} `yaml:"executable,omitempty"`

		PythonConda            struct {
			PythonVersion      string `yaml:"python_version"`
			Dependencies       string `yaml:"dependencies"`
			RelativeScriptPath string `yaml:"relative_script_path"`
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

func NewJapeConfig(flavor string, name string, version string) *JapeConfig {
	c := &JapeConfig{}
	c.Name = name
	c.Flavor = flavor
	c.Version = version
	return c
}

// Writes the given project configuration to the working directory
func WriteProjectConfig(c *JapeConfig) {

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
func ReadProjectConfig() *JapeConfig {

	if !FileExists(ConfFile) {
		return nil
	}

	var c = &JapeConfig{}

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

