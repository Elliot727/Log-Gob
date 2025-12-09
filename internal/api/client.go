// Package api provides an HTTP client for making API requests with authentication.
package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents an API client with base URL, API key, and HTTP client.
type Client struct {
	BaseURL string
	APIKey  string
	HTTP    *http.Client
}

// New creates a new API client with the given base URL and API key.
func New(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		HTTP: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// request creates an HTTP request with proper headers for API authentication.
func (c *Client) request(path string) (*http.Request, error) {
	req, err := http.NewRequest("GET", c.BaseURL+path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	return req, nil
}

// Get makes a GET request to the specified path and unmarshals the response into the provided interface.
func (c *Client) Get(path string, response interface{}) error {
	req, err := c.request(path)
	if err != nil {
		return err
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, response)
}
