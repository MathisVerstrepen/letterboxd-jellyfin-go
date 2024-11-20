package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

const confFilePath = "config.json"

type UserData struct {
	Username             string
	LatestWatchlistMovie string
	CollectionId         string
	JellyfinUserName     string
	LastFullSync         time.Time
}

type Configuration struct {
	Users           []UserData
	ProxyUrl        string
	ProxyUser       string
	ProxyPass       string
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

	configuration.ProxyUrl = os.Getenv("PROXY_URL")
	configuration.ProxyUser = os.Getenv("PROXY_USER")
	configuration.ProxyPass = os.Getenv("PROXY_PASS")

	return configuration
}

func IsLocked() bool {
	if _, err := os.Stat(filepath.Join(basepath, "app.lock")); err == nil {
		return true
	}

	os.Create(filepath.Join(basepath, "app.lock"))
	return false
}

func Unlock() {
	os.Remove(filepath.Join(basepath, "app.lock"))
}

func PersistChanges(configuration Configuration) {
	json, err := json.Marshal(configuration)

	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile(filepath.Join(basepath, confFilePath), json, 0777)
}
