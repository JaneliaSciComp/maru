package utils

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const confFile = "jape.yaml"

type JapeConfig struct {
	Flavor string
	Name string
	Version string
	GitUrl string
	GitTag string
}

func NewJapeConfig(flavor string, name string, gitUrl string) *JapeConfig {
	return &JapeConfig{flavor, name, "1.0.0", gitUrl, "master"}
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

// Reads the current project configuration from the working directory
func ReadProjectConfig() *JapeConfig {

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

