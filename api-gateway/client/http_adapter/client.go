package http_adapter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
}

// Response represents a generic HTTP response
type Response struct {
	StatusCode int                    `json:"status_code"`
	Body       map[string]interface{} `json:"body"`
	Headers    map[string]string      `json:"headers"`
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Get makes a GET request to the specified URL
func (c *Client) Get(url string) (*Response, error) {
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse JSON response
	var bodyMap map[string]interface{}
	if len(body) > 0 {
		if err := json.Unmarshal(body, &bodyMap); err != nil {
			// If JSON parsing fails, store as string
			bodyMap = map[string]interface{}{
				"raw": string(body),
			}
		}
	}

	// Extract headers
	headers := make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Body:       bodyMap,
		Headers:    headers,
	}, nil
}
