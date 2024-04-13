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
			body := client.FetchData(f.FetcherParams{
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

			log.Println(string(body))

			numberOfMoviesRemoved += 1
		}
	}

	return numberOfMoviesRemoved
}
