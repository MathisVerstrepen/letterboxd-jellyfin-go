package jellyfin

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	f "diikstra.fr/letterboxd-jellyfin-go/fetch"
)

type User struct {
	Name string
	Id   string
}

const JellyfinUrl = "https://stream.diikstra.fr/"

func GetUsers() []User {
	body := f.Fetcher(f.FetcherParams{
		Method: "GET",
		Url:    JellyfinUrl + "Users",
		Body:   nil,
		Headers: f.Header{
			"content-type": "application/json; charset=utf-8",
		},
		Params: f.Param{
			"ApiKey": os.Getenv("JELLYFIN_API_KEY"),
		},
	})

	var users []User
	json.Unmarshal(body, &users)

	return users
}

func GetUserId(userName string) (string, error) {
	users := GetUsers()

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

func GetUserViews(userId string, userCollectionId string) ([]UserView, error) {
	body := f.Fetcher(f.FetcherParams{
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

func removeSeenMoviesFromUserCollection(userId string, userCollectionId string) int {
	userViews, err := GetUserViews(userId, userCollectionId)
	numberOfMoviesRemoved := 0

	if err != nil {
		log.Printf("Failed to get user %s views in collection %s", userId, userCollectionId)
		return -1
	}

	for _, movie := range userViews {
		if movie.UserData.Played {
			log.Printf("Deleting %s of user %s from collection %s\n", movie.Name, userId, userCollectionId)
			body := f.Fetcher(f.FetcherParams{
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
			})

			log.Println(string(body))

			numberOfMoviesRemoved += 1
		}
	}

	return numberOfMoviesRemoved
}
