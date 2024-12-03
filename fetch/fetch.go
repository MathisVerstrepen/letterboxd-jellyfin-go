package fetch

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"slices"

	"golang.org/x/net/proxy"
)

type Header map[string]string
type Param map[string]string

type FetcherParams struct {
	Method       string
	Url          string
	Body         any
	Headers      Header
	Params       Param
	UseProxy     bool
	WantErrCodes []int
}

type FetcherClient interface {
	FetchData(fp FetcherParams) []byte
}

type Fetcher struct {
	ProxyUrl  string
	ProxyUser string
	ProxyPass string
}

func (f Fetcher) FetchData(fp FetcherParams) []byte {
	client := &http.Client{}

	if fp.UseProxy {
		dialer, err := proxy.SOCKS5("tcp", f.ProxyUrl, &proxy.Auth{
			User:     f.ProxyUser,
			Password: f.ProxyPass,
		}, proxy.Direct)
		if err != nil {
			log.Println("Failed to initialize proxy.")
			panic(err)
		}

		dialContext := func(ctx context.Context, network, address string) (net.Conn, error) {
			return dialer.Dial(network, address)
		}
		transport := &http.Transport{DialContext: dialContext,
			DisableKeepAlives: true}
		client = &http.Client{Transport: transport}
	}

	baseUrl, err := url.Parse(fp.Url)
	if err != nil {
		log.Println("Failed to parse url.")
		panic(err)
	}

	params := url.Values{}
	for paramKey, paramValue := range fp.Params {
		params.Add(paramKey, paramValue)
	}
	baseUrl.RawQuery = params.Encode()

	var bodyBuffer *bytes.Buffer = &bytes.Buffer{}
	if fp.Body != nil {
		jsonBytes, err := json.Marshal(fp.Body)
		if err != nil {
			log.Println("Failed to encode req body in bytes.")
			panic(err)
		}
		bodyBuffer = bytes.NewBuffer(jsonBytes)
	}

	req, err := http.NewRequest(fp.Method, baseUrl.String(), bodyBuffer)
	if err != nil {
		log.Println("Failed to initialize request.")
		panic(err)
	}

	for headerKey, headerValue := range fp.Headers {
		req.Header.Set(headerKey, headerValue)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to make request.")
		panic(err)
	}

	if fp.WantErrCodes == nil && resp.StatusCode != 200 {
		log.Printf("Got status code %d instead of wanted 200\nUrl : %s", resp.StatusCode, fp.Url)
		panic("Failed to get 200 status code.")
	} else if fp.WantErrCodes != nil && !slices.Contains(fp.WantErrCodes, resp.StatusCode) {
		log.Printf("Got status code %d instead of wanted %d\nUrl : %s", resp.StatusCode, fp.WantErrCodes, fp.Url)
		panic("Failed to get wanted status code.")
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read body from request response.")
		panic(err)
	}

	return body
}
