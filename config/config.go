package config

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	NAME       = "go-ma-actor"
	VERSION    = "v0.0.4"
	ENV_PREFIX = "GO_MA_ACTOR"
)

var configFile string = ""

func init() {

	// Look in the current directory, the home directory and /etc for the config file.
	// In that order.
	viper.SetConfigName(NAME)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.ma")
	viper.AddConfigPath("/etc/ma")

	viper.SetEnvPrefix(ENV_PREFIX)
	viper.AutomaticEnv()

	// Allow to set config file via command line flag.
	pflag.StringVarP(&configFile, "config", "c", "", "Config file to use.")

	pflag.BoolP("version", "v", false, "Print version and exit.")
	viper.BindPFlag("version", pflag.Lookup("version"))

}

// This should be called after pflag.Parse() in main.
// The name parameter is the name of the config file to search for without the extension.
func Init(configName string) error {

	if configFile != "" {
		log.Infof("Using config file: %s", configFile)
		viper.SetConfigFile(configFile)
	} else if configName != "" {
		viper.SetConfigName(configName)
	}

	err := viper.ReadInConfig()
	if err != nil {
		log.Warnf("No config file found: %s", err)
	}

	if viper.GetBool("version") {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	// This will exit when done. It will also publish if applicable.
	if viper.GetBool("generate") {
		log.Info("Generating new keyset and node identity")
		handleGenerateOrExit()
		os.Exit(0)
	}

	return nil

}
