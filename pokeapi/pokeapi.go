package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jdingus/Pokedex/internal/pokecache"
)

var cache *pokecache.Cache

func init() {
	cache = pokecache.NewCache(5 * time.Minute)
}

type LocationResponse struct {
	Count int `json:"count"`
	Next	*string `json:"next"`
	Previous *string `json:"previous"`
	Results []LocationArea `json:"results"`
}

type LocationAreaDetail struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type LocationArea struct {
	Name string `json:"name"`
	URL string `json:"url"`
}

func FetchLocationAreas(url string) (LocationResponse, error) {
	var locationData LocationResponse
	if cachedData, found := cache.Get(url); found {
		err := json.Unmarshal(cachedData, &locationData)
		if err != nil {
			return LocationResponse{}, fmt.Errorf("failed to decode cached data: %w", err)
		}
		return locationData, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return LocationResponse{}, fmt.Errorf("failed to make GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return LocationResponse{}, fmt.Errorf("request to %v returned status code %d", url, resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationResponse{}, fmt.Errorf("error reading response body: %w", err)
	}

	cache.Add(url, bodyBytes)

	err = json.Unmarshal(bodyBytes, &locationData)
	if err != nil {
		return LocationResponse{}, fmt.Errorf("error decoding API response: %w", err)
	}

	return locationData, nil
}