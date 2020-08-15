package utils

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const confFile = "jape.yaml"

type JapeConfig struct {
	Name    string
	Version string
	Flavor  string
	Config  struct {

		Repository  struct {
			Url                string
			Tag                string
		} `yaml:"repo"`

		Build       struct {
			Command            string
		} `yaml:"build"`

		PythonConda struct {
			PythonVersion      string `yaml:"python_version"`
			Dependencies       string `yaml:"dependencies"`
			RelativeScriptPath string `yaml:"relative_script_path"`
		} `yaml:"python_conda,omitempty"`

		JavaMaven   struct {
			MainClass          string `yaml:"main_class"`
		} `yaml:"java_maven,omitempty"`

		FijiMacro   struct {
			PluginDir          string `yaml:"plugin_dir"`
			MacroDir           string `yaml:"macro_dir"`
			MacroName          string `yaml:"macro_name"`
		} `yaml:"fiji_macro,omitempty"`

	} `yaml:"config"`
}

func NewJapeConfig(flavor string, name string) *JapeConfig {
	c := &JapeConfig{}
	c.Name = name
	c.Flavor = flavor
	c.Version = "1.0.0"
	c.Config.Repository.Url = "https://github.com/example/repo.git"
	c.Config.Repository.Tag = "master"
	c.Config.Build.Command = "true" // no-op by default
	return c
}

// Writes the given project configuration to the working directory
func WriteProjectConfig(c *JapeConfig) {

	raw, err := yaml.Marshal(&c)
	if err != nil {
		PrintFatal("Error creating project config: %s", err)
	}

	err = ioutil.WriteFile(confFile, raw, 0644)
	if err != nil {
		PrintFatal("Error writing project config file: %s", err)
	}
}

// Reads the current project configuration from the working directory. Returns nil if no file exists.
func ReadProjectConfig() *JapeConfig {

	if !FileExists(confFile) {
		return nil
	}

	var c = &JapeConfig{}

	yamlFile, err := ioutil.ReadFile(confFile)
	if err != nil {
		PrintFatal("Error reading config file: %s", err)
	}

	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		PrintFatal("Error reading config file: %s", err)
	}

	return c
}

