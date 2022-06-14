package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

var defaultValues = map[string]interface{}{}

const (
	envPrefix = "DMS"
)

// LoadConfig loads all of the generic file and environment settings
func LoadConfig(configFilePath string) {
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	if len(configFilePath) > 0 {
		viper.SetConfigFile(configFilePath)
	} else {
		viper.SetConfigFile("config.yaml") // name of config file (without extension)
		viper.AddConfigPath(".")           // path to look for the config file in
	}

	// set defaults first
	for key, val := range defaultValues {
		viper.SetDefault(key, val)
	}

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Warn("Config file not found, using defaults")
	}

	for _, key := range viper.AllKeys() {
		log.Tracef("%s = %+v", key, viper.Get(key))
	}
}
