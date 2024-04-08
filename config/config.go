package config

import (
	"encoding/json"
	"log"
	"os"
)

const confFilePath = "config/config.json"

type UserDate struct {
	Username             string
	LatestWatchlistMovie string
}

type Configuration struct {
	Users           []UserDate
	Proxy           string
	CollectionIds   map[string]string
	RadarrRootPaths map[string]string
}

func LoadConfiguration() Configuration {
	file, _ := os.Open(confFilePath)
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Fatal(err)
	}

	return configuration
}

func PersistChanges(configuration Configuration) {
	json, err := json.Marshal(configuration)

	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile(confFilePath, json, 0777)
}
