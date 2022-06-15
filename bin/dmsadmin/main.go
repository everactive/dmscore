package main

import (
	"encoding/json"
	"fmt"
	"github.com/everactive/dmscore/cmd"
	"github.com/everactive/dmscore/config"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	log.SetLevel(log.TraceLevel)

	var configFilePath string
	if filePath, ok := os.LookupEnv("CONFIG_FILE_PATH"); ok {
		configFilePath = filePath
	}

	config.LoadConfig(configFilePath)

	ids, err := cmd.CreateIdentityDataStore()
	if err != nil {
		panic(err)
	}

	d, err := ids.DeviceGetEnrollmentByID("29aXsKUSaV5UCU9XAZdW2NP7Buf")
	if err != nil {
		panic(err)
	}

	bytes, err := json.Marshal(&d)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))
}
