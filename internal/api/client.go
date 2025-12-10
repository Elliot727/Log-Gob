package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	BaseURL string
	APIKey  string
	HTTP    *http.Client
}

func New(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		HTTP: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) request(path string) (*http.Request, error) {
	req, err := http.NewRequest("GET", c.BaseURL+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	return req, nil
}

func (c *Client) Get(path string, out interface{}) error {
	req, err := c.request(path)
	if err != nil {
		return err
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return readErr
	}

	if resp.StatusCode != http.StatusOK {

		return fmt.Errorf("status %d: %s", resp.StatusCode, string(body))

	}

	if len(body) == 0 {
		return errors.New("empty response body")
	}

	return json.Unmarshal(body, out)
}
