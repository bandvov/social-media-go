package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type HttpClient struct {
	BaseURL   string
	Timeout   time.Duration
	Transport http.RoundTripper
}

func NewClient(baseURL string) *HttpClient {
	return &HttpClient{
		BaseURL:   baseURL,
		Timeout:   5 * time.Second,
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
}

func (c *HttpClient) GetJSON(endpoint string, result interface{}) error {
	client := &http.Client{Timeout: c.Timeout}
	resp, err := client.Get(fmt.Sprintf("%s%s", c.BaseURL, endpoint))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: status %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(result)
}
