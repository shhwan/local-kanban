package client

import (
	"net/http"
	"os"
	"time"
)

type BackendClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewBackendClient() *BackendClient {
	baseURL := os.Getenv("BACKEND_URL")
	if baseURL == "" {
		// Docker Compose内のサービス名をデフォルトとして使用
		baseURL = "http://backend:8080"
	}

	return &BackendClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}
