package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"

	"diikstra.fr/letterboxd-jellyfin-go/config"
	lt "diikstra.fr/letterboxd-jellyfin-go/letterboxd"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error while loading env file.\nErr: %s", err)
	}

	conf := config.LoadConfiguration()
	fmt.Println(conf)

	f := lt.Fetcher{
		ProxyUrl: conf.Proxy,
	}

	for index := range conf.Users {
		fmt.Println(conf.Users[index].Username)
		tmdbIds, _ := f.GetNewestUserWatchlist(conf.Users[index].Username, &conf.Users[index].LatestWatchlistMovie)
		fmt.Println(tmdbIds)

	}

	fmt.Println(conf)
	config.PersistChanges(conf)
}
