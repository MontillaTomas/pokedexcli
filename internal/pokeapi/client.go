package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MontillaTomas/pokedexcli/internal/pokecache"
)

type Client struct {
	httpClient http.Client
	BaseURL    string
	cache      *pokecache.Cache
}

func NewClient(timeout time.Duration) *Client {
	return &Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		BaseURL: "https://pokeapi.co/api/v2/",
		cache:   pokecache.NewCache(5 * time.Second),
	}
}

func (c *Client) GetLocationAreas(url string) (*LocationArea, error) {
	fullURL := url
	if fullURL == "" {
		fullURL = c.BaseURL + "location-area/"
	}

	// Check cache first
	if c.cache != nil {
		if cachedData, found := c.cache.Get(fullURL); found {
			var data LocationArea
			if err := json.Unmarshal(cachedData, &data); err == nil {
				return &data, nil
			}
		}
	}

	// Fetch from API
	res, err := c.httpClient.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var data LocationArea
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}

	// Store in cache
	if c.cache != nil {
		rawData, err := json.Marshal(data)
		if err == nil {
			c.cache.Add(fullURL, rawData)
		}
	}

	return &data, nil
}

func (c *Client) GetLocationAreaPokemons(locationAreaName string) (*LocationAreaDetails, error) {
	fullURL := c.BaseURL + "location-area/" + locationAreaName + "/"

	// Check cache first
	if c.cache != nil {
		if cachedData, found := c.cache.Get(fullURL); found {
			var data LocationAreaDetails
			if err := json.Unmarshal(cachedData, &data); err == nil {
				return &data, nil
			}
		}
	}

	// Fetch from API
	res, err := c.httpClient.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("location area '%s' not found", locationAreaName)
		}
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var data LocationAreaDetails
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}

	// Store in cache
	if c.cache != nil {
		rawData, err := json.Marshal(data)
		if err == nil {
			c.cache.Add(fullURL, rawData)
		}
	}

	return &data, nil
}

func (c *Client) GetPokemon(pokemonName string) (*Pokemon, error) {
	fullURL := c.BaseURL + "pokemon/" + pokemonName + "/"

	// Check cache first
	if c.cache != nil {
		if cachedData, found := c.cache.Get(fullURL); found {
			var data Pokemon
			if err := json.Unmarshal(cachedData, &data); err == nil {
				return &data, nil
			}
		}
	}

	// Fetch from API
	res, err := c.httpClient.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("pokemon '%s' not found", pokemonName)
		}
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var data Pokemon
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}

	// Store in cache
	if c.cache != nil {
		rawData, err := json.Marshal(data)
		if err == nil {
			c.cache.Add(fullURL, rawData)
		}
	}

	return &data, nil
}
