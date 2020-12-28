package cmd

import (
	Utils "maru/utils"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var EnvParam []string
var UserParam string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:              "maru",
	Short:            "Maru makes scientific containers fun and easy.",
	Long:             `Maru is a CLI utility for containerizing scientific applications and managing those containers.`,
	TraverseChildren: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// Error was already reported to the user
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global configuration
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.maru.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&Utils.Debug, "debug", "d", false, "print debug output")

	// Docker parameters
	rootCmd.PersistentFlags().StringArrayVarP(&EnvParam, "env", "e", nil, "Set environment variables for the running container, e.g. when using run or shell")
	rootCmd.PersistentFlags().StringVarP(&UserParam, "user", "u", "", "Set user for the running container, e.g. when using run or shell")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			Utils.PrintFatal("%s", err)
		}

		// Search config in home directory with name ".maru" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".maru")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		Utils.PrintMessage("Using config file:", viper.ConfigFileUsed())
	}
}
