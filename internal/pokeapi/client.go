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

const baseUrl = "https://pokeapi.co/api/v2/"
const locationBaseUrl = "https://pokeapi.co/api/v2/location-area"

var cache = pokecache.NewCache(30 * time.Second)

var errBadStatus = errors.New("pokeapi: unexpected http status code")

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func MakeLocationsRequest(configUrl *string) (*LocationAreas, error) {
	return makeLocationsRequestCtx(context.Background(), configUrl)
}

func MakeLocationAreaRequest(location string) (*LocationAreaResponse, error) {
	return makeLocationAreaRequestCtx(context.Background(), location)
}

func MakePokemonInfoRequest(name string) (*PokemonInfoResponse, error) {
	return makePokemonInfoRequestCtx(context.Background(), name)
}

func makeLocationsRequestCtx(ctx context.Context,
	configUrl *string) (*LocationAreas, error) {
	targetUrl := locationBaseUrl
	if configUrl != nil && len(*configUrl) > 0 {
		targetUrl = *configUrl
	}

	cachedData, ok := cache.Get(targetUrl)
	if ok {
		return decodedResponse[LocationAreas](cachedData)
	} else {
		resp, err := makeRequest(ctx, targetUrl, "GET")
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		data, err := getBytesData(resp)
		cache.Add(targetUrl, data)
		return decodedResponse[LocationAreas](data)
	}
}

func makeLocationAreaRequestCtx(ctx context.Context,
	location string) (*LocationAreaResponse, error) {
	targetUrl := locationBaseUrl + "/" + location

	cachedData, ok := cache.Get(targetUrl)
	if ok {
		return decodedResponse[LocationAreaResponse](cachedData)
	} else {
		resp, err := makeRequest(ctx, targetUrl, "GET")
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		data, err := getBytesData(resp)
		cache.Add(targetUrl, data)
		return decodedResponse[LocationAreaResponse](data)
	}
}

func makePokemonInfoRequestCtx(ctx context.Context,
	name string) (*PokemonInfoResponse, error) {
	targetUrl := baseUrl + "pokemon/" + name

	cachedData, ok := cache.Get(targetUrl)
	if ok {
		return decodedResponse[PokemonInfoResponse](cachedData)
	} else {
		resp, err := makeRequest(ctx, targetUrl, "GET")
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		data, err := getBytesData(resp)
		cache.Add(targetUrl, data)
		return decodedResponse[PokemonInfoResponse](data)
	}
}

func decodedResponse[T any](data []byte) (*T, error) {
	var result T
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func getBytesData(resp *http.Response) ([]byte, error) {
	if resp.StatusCode != http.StatusOK {
		return nil, errBadStatus
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
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
