package fetch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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
	FetchData(fp FetcherParams) ([]byte, error)
}

type Fetcher struct {
	ProxyUrl  string
	ProxyUser string
	ProxyPass string
}

func (f Fetcher) FetchData(fp FetcherParams) ([]byte, error) {
	client := &http.Client{}

	if fp.UseProxy {
		dialer, err := proxy.SOCKS5("tcp", f.ProxyUrl, &proxy.Auth{
			User:     f.ProxyUser,
			Password: f.ProxyPass,
		}, proxy.Direct)
		if err != nil {
			log.Println("Failed to initialize proxy.")
			return nil, err
		}

		dialContext := func(ctx context.Context, network, address string) (net.Conn, error) {
			return dialer.Dial(network, address)
		}
		transport := &http.Transport{
			DialContext:       dialContext,
			DisableKeepAlives: true,
		}
		client = &http.Client{Transport: transport}
	}

	baseUrl, err := url.Parse(fp.Url)
	if err != nil {
		log.Println("Failed to parse url.")
		return nil, err
	}

	params := url.Values{}
	for paramKey, paramValue := range fp.Params {
		params.Add(paramKey, paramValue)
	}
	baseUrl.RawQuery = params.Encode()

	bodyBuffer := &bytes.Buffer{}
	if fp.Body != nil {
		jsonBytes, err := json.Marshal(fp.Body)
		if err != nil {
			log.Println("Failed to encode req body in bytes.")
			return nil, err
		}
		bodyBuffer = bytes.NewBuffer(jsonBytes)
	}

	req, err := http.NewRequest(fp.Method, baseUrl.String(), bodyBuffer)
	if err != nil {
		log.Println("Failed to initialize request.")
		return nil, err
	}

	for headerKey, headerValue := range fp.Headers {
		req.Header.Set(headerKey, headerValue)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to make request.")
		return nil, err
	}
	defer resp.Body.Close()

	if fp.WantErrCodes == nil && resp.StatusCode != 200 {
		log.Printf("Got status code %d instead of wanted 200\nUrl : %s", resp.StatusCode, fp.Url)
		return nil, errors.New("failed to get 200 status code")
	} else if fp.WantErrCodes != nil && !slices.Contains(fp.WantErrCodes, resp.StatusCode) {
		log.Printf("Got status code %d instead of wanted %v\nUrl : %s", resp.StatusCode, fp.WantErrCodes, fp.Url)
		return nil, errors.New("failed to get wanted status code")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read body from request response.")
		return nil, err
	}

	return body, nil
}
