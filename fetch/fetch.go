package fetch

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Header map[string]string
type Param map[string]string

type FetcherParams struct {
	Method  string
	Url     string
	Body    any
	Headers Header
	Params  Param
}

func Fetcher(fp FetcherParams) []byte {
	client := &http.Client{}

	baseUrl, err := url.Parse(fp.Url)
	if err != nil {
		log.Fatalf("Failed to parse url.\nErr : %s", err)
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
			log.Fatalf("Failed to encode req body in bytes.\nErr : %s", err)
		}
		bodyBuffer = bytes.NewBuffer(jsonBytes)
	}

	req, err := http.NewRequest(fp.Method, baseUrl.String(), bodyBuffer)
	if err != nil {
		log.Fatalf("Failed to initialize request.\nErr : %s", err)
	}

	for headerKey, headerValue := range fp.Headers {
		req.Header.Set(headerKey, headerValue)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to make request.\nErr : %s", err)
	}

	if resp.StatusCode != 200 {
		log.Printf("Warn: Got status code %d instead of standard 200", resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read body from request response.\nErr : %s", err)
	}

	return body
}
