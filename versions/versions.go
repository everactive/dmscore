// Package versions provides a way to read component versions from disk, cache and provide
package versions

import (
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/everactive/dmscore/iot-management/config/configkey"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Version is the variables set by compilation needed for version information
var (
	Version = "VERSION_UNSET"
)

// String returns the version string
func String() string {
	return Version
}

const versionsPath = "/opt/component/versions"

var cachedVersions map[string]string
var lastVersionTime time.Time

// GetComponentVersions gets the component versions from disk or returns a cached copy that has a configurable duration
func GetComponentVersions() map[string]string {
	cacheDuration := viper.GetDuration(configkey.ComponentVersionsCacheDuration)
	if cachedVersions == nil || lastVersionTime.Add(cacheDuration).Before(time.Now()) {
		log.Tracef("Cached versions expired or does not exist, creating/renewing")
		lastVersionTime = time.Now()
		cachedVersions = make(map[string]string)
		cachedVersions["management"] = String()

		dirEntries, err := os.ReadDir(versionsPath)
		if err != nil {
			log.Error(err)
			return cachedVersions
		}

		for _, dirEntry := range dirEntries {
			if dirEntry.IsDir() {
				fileBytes, err := ioutil.ReadFile(path.Join(versionsPath, dirEntry.Name(), "version"))
				if err != nil {
					log.Error(err)
					continue
				}

				cachedVersions[dirEntry.Name()] = string(fileBytes)
			} else {
				log.Warnf("Unexpected file in versions directory: %s", dirEntry.Name())
			}
		}
	} else {
		log.Tracef("Using cached versions until %s", lastVersionTime)
	}

	return cachedVersions
}
