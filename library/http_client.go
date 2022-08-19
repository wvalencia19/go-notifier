package library

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"notifier/config"
)

//go:generate mockgen -source=library/http_client.go -destination=library/mock/http_client.go
type Client interface {
	SendNotification(ctx context.Context, message string) error
}

type client struct {
	url    *url.URL
	client *http.Client
}

func NewClient(url *url.URL, httpConfig config.HTTPClient) Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = httpConfig.MaxIdleConns
	t.MaxIdleConnsPerHost = httpConfig.MaxIdleConnsPerHost

	return &client{
		url: url,
		client: &http.Client{
			Timeout:   httpConfig.HTTPRequestTimeout,
			Transport: t,
		},
	}
}

func (c *client) SendNotification(ctx context.Context, m string) error {
	message := map[string]string{"message": m}
	jsonBody, _ := json.Marshal(message)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url.String(), bytes.NewBuffer(jsonBody))

	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}

	statusCode := res.StatusCode
	if statusCode < 200 || statusCode > 299 {
		return fmt.Errorf("sending notification with status code: %d", statusCode)
	}
	return nil
}
