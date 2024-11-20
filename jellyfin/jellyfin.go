package jellyfin

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"strings"

	f "diikstra.fr/letterboxd-jellyfin-go/fetch"
	rd "diikstra.fr/letterboxd-jellyfin-go/radarr"
)

type User struct {
	Name string
	Id   string
}

const JellyfinUrl = "https://stream.diikstra.fr/"

func GetUsers(client f.FetcherClient) []User {
	body := client.FetchData(f.FetcherParams{
		Method: "GET",
		Url:    JellyfinUrl + "Users",
		Body:   nil,
		Headers: f.Header{
			"content-type": "application/json; charset=utf-8",
		},
		Params: f.Param{
			"ApiKey": os.Getenv("JELLYFIN_API_KEY"),
		},
		UseProxy: false,
	})

	var users []User
	json.Unmarshal(body, &users)

	return users
}

func GetUserId(client f.FetcherClient, userName string) (string, error) {
	users := GetUsers(client)

	var userId string
	for _, user := range users {
		if user.Name == userName {
			userId = user.Id
			break
		}
	}

	if userId == "" {
		return "", errors.New("no user matching found")
	}

	return userId, nil
}

type UserData struct {
	Played bool
}

type UserView struct {
	Name     string
	Id       string
	UserData UserData
}

type ReqUserViewWrapper struct {
	Items []UserView
}

func GetUserViews(client f.FetcherClient, userId string, userCollectionId string) ([]UserView, error) {
	body := client.FetchData(f.FetcherParams{
		Method: "GET",
		Url:    JellyfinUrl + "Items",
		Body:   nil,
		Headers: f.Header{
			"content-type": "application/json; charset=utf-8",
		},
		Params: f.Param{
			"ApiKey":           os.Getenv("JELLYFIN_API_KEY"),
			"ParentId":         userCollectionId,
			"Recursive":        "true",
			"IncludeItemTypes": "Movie",
			"enableUserData":   "true",
			"userId":           userId,
		},
	})

	var userView ReqUserViewWrapper
	json.Unmarshal(body, &userView)

	return userView.Items, nil
}

func RemoveSeenMoviesFromUserCollection(client f.FetcherClient, userId string, userCollectionId string) int {
	userViews, err := GetUserViews(client, userId, userCollectionId)
	numberOfMoviesRemoved := 0

	if err != nil {
		log.Printf("Failed to get user %s views in collection %s", userId, userCollectionId)
		return -1
	}

	for _, movie := range userViews {
		if movie.UserData.Played {
			log.Printf("Deleting %s of user %s from collection %s\n", movie.Name, userId, userCollectionId)
			client.FetchData(f.FetcherParams{
				Method: "DELETE",
				Url:    JellyfinUrl + "Collections/" + userCollectionId + "/Items",
				Body:   nil,
				Headers: f.Header{
					"content-type": "application/json; charset=utf-8",
				},
				Params: f.Param{
					"ApiKey": os.Getenv("JELLYFIN_API_KEY"),
					"ids":    movie.Id,
				},
				WantErrCodes: []int{204},
			})

			numberOfMoviesRemoved += 1
		}
	}

	return numberOfMoviesRemoved
}

type MoviesItem struct {
	Name           string
	ProductionYear int
	Id             string
}

type Movies struct {
	Items []MoviesItem
}

func GetAllMovies(client f.FetcherClient) *[]MoviesItem {
	body := client.FetchData(f.FetcherParams{
		Method: "GET",
		Url:    JellyfinUrl + "Items",
		Body:   nil,
		Headers: f.Header{
			"content-type": "application/json; charset=utf-8",
		},
		Params: f.Param{
			"ApiKey":           os.Getenv("JELLYFIN_API_KEY"),
			"Recursive":        "true",
			"IncludeItemTypes": "Movie",
			"fields":           "MediaSources,People",
		},
		UseProxy: false,
	})

	var res Movies
	json.Unmarshal(body, &res)

	return &res.Items
}

func GetMovieJellyfinId(movies *[]MoviesItem, movie_name string, movie_year int) (string, error) {
	for _, movie := range *movies {
		if movie.Name == movie_name && movie.ProductionYear == movie_year {
			return movie.Id, nil
		}
	}

	return "", errors.New("unable to find movie in the Jellyfin library")
}

func AddMoviesToCollection(client f.FetcherClient, allMovies *[]MoviesItem, radarrStates []rd.RadarrStatus, userId string, userCollectionId string) {
	const batchSize = 20
	var ids []string

	for _, state := range radarrStates {
		jellyfinId, err := GetMovieJellyfinId(allMovies, state.Title, state.ProductionYear)
		if err == nil {
			ids = append(ids, jellyfinId)
		}
	}

	for i := 0; i < len(ids); i += batchSize {
		end := i + batchSize
		if end > len(ids) {
			end = len(ids)
		}

		batch := ids[i:end]
		client.FetchData(f.FetcherParams{
			Method: "POST",
			Url:    JellyfinUrl + "Collections/" + userCollectionId + "/Items",
			Body:   nil,
			Headers: f.Header{
				"content-type": "application/json; charset=utf-8",
			},
			Params: f.Param{
				"ApiKey": os.Getenv("JELLYFIN_API_KEY"),
				"ids":    joinIds(batch),
			},
			WantErrCodes: []int{204},
		})
	}
}

func joinIds(ids []string) string {
	return strings.Join(ids, ",")
}
