package library

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"notifier/config"
	"testing"
	"time"
)

func TestClient_SendNotification(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		statusCode  int
		withTimeout bool
		wantErr     bool
	}{
		{
			name:        "with successfully response",
			statusCode:  http.StatusCreated,
			ctx:         context.Background(),
			withTimeout: false,
			wantErr:     false,
		},
		{
			name:        "with timeout response",
			ctx:         context.Background(),
			withTimeout: true,
			wantErr:     true,
		},
		{
			name:        "with bad response",
			statusCode:  http.StatusBadRequest,
			ctx:         context.Background(),
			withTimeout: false,
			wantErr:     true,
		},
		{
			name:        "with canceled context",
			statusCode:  http.StatusOK,
			ctx:         cancelledContext(),
			withTimeout: false,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpClient := newHTTPClient(tt.statusCode)
			testURL := buildURL(chooseClient(httpClient, tt.withTimeout).URL)
			client := NewClient(testURL, *defaultHTTPConfig())
			err := client.SendNotification(tt.ctx, "message")
			if err != nil && tt.wantErr == false {
				t.Errorf("unexpected error %v", err)
			}
		})
	}
}

type httpClient struct {
	statusCode int
}

func newHTTPClient(statusCode int) httpClient {
	return httpClient{statusCode: statusCode}
}

func (h httpClient) buildClient() *httptest.Server {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(h.statusCode)
	}))
	return srv
}

func (h httpClient) buildWithTimeoutClient() *httptest.Server {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(20 * time.Millisecond)
	}))
	http.DefaultTransport.(*http.Transport).ResponseHeaderTimeout = 10 * time.Millisecond
	return srv
}

func buildURL(u string) *url.URL {
	uri, _ := url.ParseRequestURI(u)

	return uri
}

func chooseClient(httpClient httpClient, withTimeout bool) *httptest.Server {
	if withTimeout {
		return httpClient.buildWithTimeoutClient()
	}
	return httpClient.buildClient()
}

func defaultHTTPConfig() *config.HTTPClient {
	return &config.HTTPClient{
		HTTPRequestTimeout:  1 * time.Second,
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
	}
}
