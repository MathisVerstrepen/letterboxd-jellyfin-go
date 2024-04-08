package letterboxd

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/proxy"

	gs "diikstra.fr/letterboxd-jellyfin-go/gosoup"
)

var ErrProxy = fmt.Errorf("fail to initialize proxy")
var ErrFetch = fmt.Errorf("fail to fetch page")
var ErrParse = fmt.Errorf("fail to parse")

const numMoviesWatchlistPage = 28
const letterboxdUrl = "https://letterboxd.com/"

type Fetcher struct {
	ProxyUrl string
}

func (f Fetcher) letterboxdGetFetcher(endpoint string) (*html.Node, error) {
	dialer, err := proxy.SOCKS5("tcp", f.ProxyUrl, nil, proxy.Direct)
	if err != nil {
		return nil, ErrProxy
	}

	dialContext := func(ctx context.Context, network, address string) (net.Conn, error) {
		return dialer.Dial(network, address)
	}
	transport := &http.Transport{DialContext: dialContext,
		DisableKeepAlives: true}
	cl := &http.Client{Transport: transport}

	resp, err := cl.Get(endpoint)

	if err != nil || resp.StatusCode != 200 {
		fmt.Println(err)
		return nil, ErrFetch
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, ErrParse
	}

	parsedBody, err := html.Parse(strings.NewReader(string(body)))

	if err != nil {
		return nil, ErrParse
	}

	return parsedBody, nil
}

func (f Fetcher) letterboxdGetFetcherWithRetry(endpoint string) (*html.Node, error) {
	numFetch := 0
	var err error
	for numFetch < 3 {
		node, err := f.letterboxdGetFetcher(endpoint)

		if err == nil {
			return node, nil
		}
		numFetch += 1
		fmt.Println("fetch failed, retrying...")
	}
	fmt.Println("fetch failed after 3 retries, aborting...")
	return nil, err
}

func (f Fetcher) getTmdbIdFromSlug(dataTargetLink string) (string, error) {
	node, err := f.letterboxdGetFetcherWithRetry(letterboxdUrl + dataTargetLink)

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

func (f Fetcher) GetNewestUserWatchlist(userName string, latestFetched *string) ([]string, error) {
	pageIndex := 1
	var tmdbIds []string

	for pageIndex > 0 {
		fmt.Printf("Fetching page %d\n", pageIndex)
		node, err := f.letterboxdGetFetcherWithRetry(letterboxdUrl + userName + "/watchlist/page/" + fmt.Sprint(pageIndex))

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
			tmdbId, err := f.getTmdbIdFromSlug(dataTargetLink[1:])
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
		}
		pageIndex += 1
	}

	if len(tmdbIds) > 0 {
		*latestFetched = tmdbIds[0]
	}

	return tmdbIds, nil
}
