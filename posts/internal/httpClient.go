package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type ClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
	GetJSON(req *http.Request, target interface{}) error
}

// HTTPClient wraps an HTTP client with timeout settings.
type HTTPClient struct {
	Client    *http.Client
	BasePath  string
	Transport http.RoundTripper
}

// NewHTTPClient initializes a new HTTP client.
func NewHTTPClient(basePath string) *HTTPClient {

	return &HTTPClient{
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
		BasePath:  basePath,
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
}

// Do sends an HTTP request and returns an HTTP response.
func (h *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	if h.BasePath != "" {
		// Ensure no double slashes
		fullURL := strings.TrimRight(h.BasePath, "/") + "/" + strings.TrimLeft(req.URL.Path, "/")

		parsedURL, err := url.Parse(fullURL)
		if err != nil {
			return nil, fmt.Errorf("invalid URL: %v", err)
		}
		req.URL = parsedURL

		// Debug print
		fmt.Fprintln(os.Stdout, "Final URL:", req.URL.String())
	}

	return h.Client.Do(req)
}

// GetJSON performs a GET request and decodes the response into the provided struct.
func (h *HTTPClient) GetJSON(req *http.Request, target interface{}) error {
	fmt.Fprintf(os.Stdout, "%+v\n", req)
	resp, err := h.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}
