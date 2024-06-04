package main

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"

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

	for index := range conf.Users {
		fmt.Println(conf.Users[index].Username)
		tmdbIds, _ := letterboxdScrapper.GetNewestUserWatchlist(conf.Users[index].Username, &conf.Users[index].LatestWatchlistMovie)
		fmt.Println(tmdbIds)

		radarrStates := rd.SendTmdbIDsToRadarr(fetcher, tmdbIds, &conf)

		userId, err := jf.GetUserId(fetcher, conf.Users[index].JellyfinUserName)
		if err != nil {
			log.Fatalf("No user matching in Jellyfin found for %s", conf.Users[index].JellyfinUserName)
		}
		jf.RemoveSeenMoviesFromUserCollection(fetcher, userId, conf.Users[index].CollectionId)
		jf.AddMoviesToCollection(fetcher, allMovies, radarrStates, userId, conf.Users[index].CollectionId)
	}

	config.PersistChanges(conf)
}
