package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

const confFilePath = "config.json"

type UserDate struct {
	Username             string
	LatestWatchlistMovie string
	CollectionId         string
	JellyfinUserName     string
}

type Configuration struct {
	Users           []UserDate
	Proxy           string
	CollectionIds   map[string]string
	RadarrRootPaths map[string]string
}

func LoadConfiguration() Configuration {
	file, err := os.Open(filepath.Join(basepath, confFilePath))
	if err != nil {
		log.Println("Fail to open config file")
		log.Fatal(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Println("Fail to decode config file")
		log.Fatal(err)
	}

	return configuration
}

func PersistChanges(configuration Configuration) {
	json, err := json.Marshal(configuration)

	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile(filepath.Join(basepath, confFilePath), json, 0777)
}
