package pokeapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/oleg-rocks/pokedexcli/internal/pokecache"
)

const baseUrl = "https://pokeapi.co/api/v2/location-area"

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

var cache = pokecache.NewCache(30 * time.Second)

var errBadStatus = errors.New("pokeapi: unexpected http status code")

func MakeLocationsRequest(configUrl *string) (*LocationAreas, error) {
	return MakeLocationsRequestCtx(context.Background(), configUrl)
}

func MakeLocationAreaRequest(location string) (*LocationAreaResponse, error) {
	return MakeLocationAreaRequestCtx(context.Background(), location)
}

func MakeLocationsRequestCtx(ctx context.Context,
	configUrl *string) (*LocationAreas, error) {
	targetUrl := baseUrl
	if configUrl != nil && len(*configUrl) > 0 {
		targetUrl = *configUrl
	}
	cacheResp, ok := cache.Get(targetUrl)
	if ok {
		var locations LocationAreas
		err := json.Unmarshal(cacheResp, &locations)
		if err != nil {
			return nil, err
		}
		return &locations, nil
	} else {
		resp, err := makeRequest(ctx, targetUrl, "GET")
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, errBadStatus
		}

		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		cache.Add(targetUrl, bytes)

		var locations LocationAreas
		err = json.Unmarshal(bytes, &locations)
		if err != nil {
			return nil, err
		}
		return &locations, nil
	}
}

func MakeLocationAreaRequestCtx(ctx context.Context,
	location string) (*LocationAreaResponse, error) {
	targetUrl := baseUrl + "/" + location

	cacheResp, ok := cache.Get(targetUrl)
	if ok {
		var locationArea LocationAreaResponse
		err := json.Unmarshal(cacheResp, &locationArea)
		if err != nil {
			return nil, err
		}
		return &locationArea, nil
	} else {
		resp, err := makeRequest(ctx, targetUrl, "GET")
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, errBadStatus
		}

		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		cache.Add(targetUrl, bytes)

		var locationArea LocationAreaResponse
		err = json.Unmarshal(bytes, &locationArea)
		if err != nil {
			return nil, err
		}
		return &locationArea, nil
	}
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
