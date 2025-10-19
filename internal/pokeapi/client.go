package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	httpClient http.Client
	BaseURL    string
}

func NewClient(timeout time.Duration) *Client {
	return &Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		BaseURL: "https://pokeapi.co/api/v2/",
	}
}

func (c *Client) GetLocationAreas(url string) (*LocationAreaResponse, error) {
	fullURL := url
	if fullURL == "" {
		fullURL = c.BaseURL + "location-area/"
	}

	res, err := c.httpClient.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var data LocationAreaResponse
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}
