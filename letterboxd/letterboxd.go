package letterboxd

import (
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/net/html"

	f "diikstra.fr/letterboxd-jellyfin-go/fetch"
	gs "diikstra.fr/letterboxd-jellyfin-go/gosoup"
)

var ErrParse = fmt.Errorf("fail to parse")

const numMoviesWatchlistPage = 28
const letterboxdUrl = "https://letterboxd.com/"

type LetterboxdScrapper struct {
	Client f.FetcherClient
}

func (ls LetterboxdScrapper) letterboxdGetFetcher(endpoint string) (*html.Node, error) {

	body := ls.Client.FetchData(f.FetcherParams{
		Method:   "GET",
		Url:      endpoint,
		UseProxy: true,
	})

	parsedBody, err := html.Parse(strings.NewReader(string(body)))

	if err != nil {
		return nil, ErrParse
	}

	return parsedBody, nil
}

func (ls LetterboxdScrapper) letterboxdGetFetcherWithRetry(endpoint string) (*html.Node, error) {
	numFetch := 0
	var err error
	for numFetch < 3 {
		node, err := ls.letterboxdGetFetcher(endpoint)

		if err == nil {
			return node, nil
		}
		numFetch += 1
		fmt.Println("fetch failed, retrying...")

		time.Sleep(2 * time.Second)
	}
	fmt.Println("fetch failed after 3 retries, aborting...")
	return nil, err
}

func (ls LetterboxdScrapper) getTmdbIdFromSlug(dataTargetLink string) (string, error) {
	node, err := ls.letterboxdGetFetcherWithRetry(letterboxdUrl + dataTargetLink)

	if err != nil {
		log.Println(err)
		return "", err
	}

	body := gs.GetNodeByClass(node, &gs.HtmlSelector{
		ClassNames: "film",
		Tag:        "body",
		Multiple:   false,
	})

	return gs.GetAttribute(body[0], "data-tmdb-id"), nil
}

func (ls LetterboxdScrapper) GetNewestUserWatchlist(userName string, latestFetched *string) ([]string, error) {
	pageIndex := 1
	var tmdbIds []string

	for pageIndex > 0 {
		fmt.Printf("Fetching page %d\n", pageIndex)
		node, err := ls.letterboxdGetFetcherWithRetry(letterboxdUrl + userName + "/watchlist/page/" + fmt.Sprint(pageIndex))

		if err != nil {
			log.Println(err)
			break
		}

		posters := gs.GetNodeByClass(node, &gs.HtmlSelector{
			ClassNames: "really-lazy-load poster film-poster",
			Tag:        "div",
			Multiple:   true,
		})

		if len(posters) < numMoviesWatchlistPage {
			pageIndex = -1
		}

		for _, poster := range posters {
			dataTargetLink := gs.GetAttribute(poster, "data-target-link")
			tmdbId, err := ls.getTmdbIdFromSlug(dataTargetLink[1:])
			fmt.Printf("%s -> %s\n", dataTargetLink, tmdbId)
			if err != nil {
				log.Println(err)
				continue
			}

			if *latestFetched == tmdbId {
				pageIndex = -1
				break
			}

			tmdbIds = append(tmdbIds, tmdbId)

			time.Sleep(1 * time.Second)
		}
		pageIndex += 1
	}

	if len(tmdbIds) > 0 {
		*latestFetched = tmdbIds[0]
	}

	return tmdbIds, nil
}

func (ls LetterboxdScrapper) GetFullUserWatchlist(userName string) ([]string, error) {
	pageIndex := 1
	var tmdbIds []string

	for pageIndex > 0 {
		fmt.Printf("Fetching page %d\n", pageIndex)
		node, err := ls.letterboxdGetFetcherWithRetry(letterboxdUrl + userName + "/watchlist/page/" + fmt.Sprint(pageIndex))

		if err != nil {
			log.Println(err)
			break
		}

		posters := gs.GetNodeByClass(node, &gs.HtmlSelector{
			ClassNames: "really-lazy-load poster film-poster",
			Tag:        "div",
			Multiple:   true,
		})

		if len(posters) < numMoviesWatchlistPage {
			pageIndex = -1
		}

		for _, poster := range posters {
			dataTargetLink := gs.GetAttribute(poster, "data-target-link")
			tmdbId, err := ls.getTmdbIdFromSlug(dataTargetLink[1:])
			fmt.Printf("%s -> %s\n", dataTargetLink, tmdbId)
			if err != nil {
				log.Println(err)
				continue
			}

			tmdbIds = append(tmdbIds, tmdbId)

			time.Sleep(1 * time.Second)
		}
		pageIndex += 1
	}

	return tmdbIds, nil
}
