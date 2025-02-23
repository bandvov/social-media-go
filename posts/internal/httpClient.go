package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

// HTTPClient wraps an HTTP client with timeout settings.
type HTTPClient struct {
	Client   *http.Client
	BasePath string
}

// NewHTTPClient initializes a new HTTP client.
func NewHTTPClient(basePath string) ClientInterface {

	return &HTTPClient{
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
		BasePath: basePath,
	}
}

// Do sends an HTTP request and returns an HTTP response.
func (h *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	if h.BasePath != "" {
		req.URL.Path = h.BasePath + req.URL.Path
	}
	return h.Client.Do(req)
}

// GetJSON performs a GET request and decodes the response into the provided struct.
func (h *HTTPClient) GetJSON(req *http.Request, target interface{}) error {
	resp, err := h.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}
