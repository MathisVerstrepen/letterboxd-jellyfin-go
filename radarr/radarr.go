package radarr

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"

	"diikstra.fr/letterboxd-jellyfin-go/config"

	f "diikstra.fr/letterboxd-jellyfin-go/fetch"
)

const RadarrUrl = "http://localhost:7878/api/v3/"

type RadarrState string

type MovieLookupFile struct {
	RelativePath string `json:"relativePath"`
}

type RadarrMovieLookupResp struct {
	MovieFile MovieLookupFile `json:"movieFile"`
	Monitored bool            `json:"monitored"`
	Title     string          `json:"title"`
	TmdbId    int             `json:"tmdbId"`
	Year      int             `json:"year"`
	Genres    []string        `json:"genres"`
}

type RadarrStatus struct {
	HasFile        bool
	Monitored      bool
	Title          string
	TmdbId         string
	ProductionYear int
	IsAnimation    bool
}

func GetRadarrState(client f.FetcherClient, tmdbId string) (RadarrStatus, error) {
	body, err := client.FetchData(f.FetcherParams{
		Method: "GET",
		Url:    RadarrUrl + "movie/lookup",
		Body:   nil,
		Headers: f.Header{
			"X-Api-Key": os.Getenv("RADARR_API_KEY"),
		},
		Params: f.Param{
			"term": "tmdb:" + tmdbId,
		},
	})

	if err != nil {
		log.Printf("Failed to get movie from Radarr: %v", err)
		return RadarrStatus{}, err
	}

	parsedBody := []RadarrMovieLookupResp{}
	err = json.Unmarshal(body, &parsedBody)
	if err != nil {
		log.Println("Failed to parse JSON.")
		panic(err)
	}

	if len(parsedBody) == 0 {
		log.Printf("TMDB id %s return empty state\n", tmdbId)
		return RadarrStatus{}, errors.New("return empty state")
	}

	return RadarrStatus{
		HasFile:        parsedBody[0].MovieFile != MovieLookupFile{} && parsedBody[0].MovieFile.RelativePath != "",
		Monitored:      parsedBody[0].Monitored,
		Title:          parsedBody[0].Title,
		TmdbId:         fmt.Sprint(parsedBody[0].TmdbId),
		ProductionYear: parsedBody[0].Year,
		IsAnimation:    slices.Contains(parsedBody[0].Genres, "Animation"),
	}, nil
}

type RadarrAddBodyAddOptions struct {
	SearchForMovie bool `json:"searchForMovie"`
}

type RadarrAddBody struct {
	TmdbId           string                  `json:"tmdbId,omitempty"`
	Title            string                  `json:"title,omitempty"`
	Year             int                     `json:"year,omitempty"`
	QualityProfileId int                     `json:"qualityProfileId,omitempty"`
	Monitored        bool                    `json:"monitored,omitempty"`
	RootFolderPath   string                  `json:"rootFolderPath,omitempty"`
	AddOptions       RadarrAddBodyAddOptions `json:"addOptions,omitempty"`
}

func AddToRadarrDownload(client f.FetcherClient, movie RadarrStatus, conf *config.Configuration) {
	rootFolderPath := conf.RadarrRootPaths["movies"]
	if movie.IsAnimation {
		rootFolderPath = conf.RadarrRootPaths["anime_movies"]
	}

	reqBody := RadarrAddBody{
		TmdbId:           movie.TmdbId,
		Title:            movie.Title,
		Year:             movie.ProductionYear,
		QualityProfileId: 11,
		Monitored:        true,
		RootFolderPath:   rootFolderPath,
		AddOptions: RadarrAddBodyAddOptions{
			SearchForMovie: true,
		},
	}

	client.FetchData(f.FetcherParams{
		Method: "POST",
		Url:    RadarrUrl + "movie",
		Body:   reqBody,
		Headers: f.Header{
			"X-Api-Key":    os.Getenv("RADARR_API_KEY"),
			"Content-Type": "application/json",
		},
		Params:       f.Param{},
		WantErrCodes: []int{201, 400},
	})
}

func SendTmdbIDsToRadarr(client f.FetcherClient, tmdbIds []string, conf *config.Configuration) []RadarrStatus {
	var states []RadarrStatus

	for _, tmdbId := range tmdbIds {
		if tmdbId != "" {
			state, err := GetRadarrState(client, tmdbId)
			if err != nil {
				continue
			}
			AddToRadarrDownload(client, state, conf)
			states = append(states, state)
		}
	}

	return states
}
