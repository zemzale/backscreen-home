package primitives

import (
	"net/http"
	"time"
)

type HTTPClientOption func(*http.Client)

func WithTimeout(timeout time.Duration) HTTPClientOption {
	return func(c *http.Client) {
		c.Timeout = timeout
	}
}

// NewHTTPClient with sensible defaults
func NewHTTPClient(opts ...HTTPClientOption) *http.Client {
	client := &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			MaxIdleConns:          10,
			IdleConnTimeout:       time.Second * 60,
			MaxConnsPerHost:       10,
			ResponseHeaderTimeout: time.Second * 10,
		},
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}
