package config

import (
	"github.com/everactive/dmscore/config/keys"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

var defaultValues = map[string]interface{}{
	keys.MQTTHealthTopic:                            "devices/health/+",
	keys.MQTTPubTopic:                               "devices/pub/+",
	keys.MQTTCertificatesPath:                       "/srv/devicetwin-certs",
	keys.MQTTClientCertificateFilename:              "server.crt",
	keys.MQTTClientKeyFilename:                      "server.key",
	keys.MQTTRootCAFilename:                         "ca.crt",
	keys.MQTTClientIDPrefix:                         "devicetwin",
	keys.MQTTURL:                                    "mqtt",
	keys.MQTTPort:                                   "8883",
	keys.ServiceScheme:                              "http",
	keys.ServicePort:                                "8010",
	keys.ServiceHost:                                "localhost:8080",
	keys.DatabaseDriver:                             "postgres",
	keys.DatabaseName:                               "management",
	keys.GetDeviceTwinKey(keys.DatabaseDriver):      "postgres",
	keys.GetDeviceTwinKey(keys.DatabaseName):        "devicetwin",
	keys.GetDeviceTwinKey(keys.ServicePort):         "8030",
	keys.GetDeviceTwinKey(keys.MigrationsSourceURL): "/migrations/devicetwin",
	keys.GetIdentityKey(keys.DatabaseDriver):        "postgres",
	keys.GetIdentityKey(keys.DatabaseName):          "identity",
	keys.GetIdentityKey(keys.ServicePortInternal):   "8041",
	keys.GetIdentityKey(keys.ServicePortEnroll):     "8040",
	keys.GetIdentityKey(keys.MigrationsSourceURL):   "/migrations/identity",
	keys.GetIdentityKey(keys.CertificatesPath):      "/srv/identity-certs",
	keys.DefaultServiceHeartbeat:                    "60s",
	keys.RequiredSnapsInstallServiceCheckInterval:   "30s",
	keys.RefreshSnapListOnAnyChange:                 false,
	keys.RequiredSnapsCheckInterval:                 "100ms",
}

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
