package pokeapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

const baseUrl = "https://pokeapi.co/api/v2/location-area"

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

var errBadStatus = errors.New("pokeapi: unexpected http status code")

func MakeLocationsRequest(configUrl *string) (*LocationAreas, error) {
	return MakeLocationsRequestCtx(context.Background(), configUrl)
}

func MakeLocationsRequestCtx(ctx context.Context,
	configUrl *string) (*LocationAreas, error) {
	targetUrl := baseUrl
	if configUrl != nil && len(*configUrl) > 0 {
		targetUrl = *configUrl
	}
	resp, err := makeRequest(ctx, targetUrl, "GET")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errBadStatus
	}

	var locations LocationAreas
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&locations)
	if err != nil {
		return nil, err
	}
	return &locations, nil
}

func makeRequest(ctx context.Context, url string,
	method string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	resp, err := httpClient.Do(req)
	return resp, err
}
