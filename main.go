package main

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"time"

	"github.com/joho/godotenv"

	"diikstra.fr/letterboxd-jellyfin-go/config"
	f "diikstra.fr/letterboxd-jellyfin-go/fetch"
	jf "diikstra.fr/letterboxd-jellyfin-go/jellyfin"
	lt "diikstra.fr/letterboxd-jellyfin-go/letterboxd"
	rd "diikstra.fr/letterboxd-jellyfin-go/radarr"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func main() {
	err := godotenv.Load(filepath.Join(basepath, ".env"))
	if err != nil {
		log.Fatalf("Error while loading env file.\nErr: %s", err)
	}

	if config.IsLocked() {
		log.Fatal("App is locked, wait for the current run to finish.")
	}
	defer config.Unlock()

	conf := config.LoadConfiguration()

	fetcher := f.Fetcher{
		ProxyUrl:  conf.ProxyUrl,
		ProxyUser: conf.ProxyUser,
		ProxyPass: conf.ProxyPass,
	}
	letterboxdScrapper := lt.LetterboxdScrapper{
		Client: fetcher,
	}

	allMovies := jf.GetAllMovies(fetcher)
	currentTime := time.Now()

	for index := range conf.Users {
		fmt.Println(conf.Users[index].Username)
		var tmdbIds []string

		if currentTime.Sub(conf.Users[index].LastFullSync) > 24*time.Hour {
			fmt.Println("Last full sync was less than 24 hours ago. Full syncing.")
			tmdbIds, err = letterboxdScrapper.GetFullUserWatchlist(conf.Users[index].Username)

			if err != nil {
				panic(err)
			}

			conf.Users[index].LastFullSync = currentTime
		} else {
			tmdbIds, err = letterboxdScrapper.GetNewestUserWatchlist(conf.Users[index].Username, &conf.Users[index].LatestWatchlistMovie)

			if err != nil {
				panic(err)
			}
		}

		radarrStates := rd.SendTmdbIDsToRadarr(fetcher, tmdbIds, &conf)

		userId, err := jf.GetUserId(fetcher, conf.Users[index].JellyfinUserName)
		if err != nil {
			panic(err)
		}
		jf.RemoveSeenMoviesFromUserCollection(fetcher, userId, conf.Users[index].CollectionId)
		jf.AddMoviesToCollection(fetcher, allMovies, radarrStates, userId, conf.Users[index].CollectionId)
	}

	config.PersistChanges(conf)
}
