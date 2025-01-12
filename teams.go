package teams

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client is a type that is capable of posting messages/notifications to a channel
// in Microsoft Teams using webhooks.
type Client struct {
	webhook    string
	httpClient HttpClient
}

// New initializes a new Client
func New(webhook string, opts ...Option) *Client {
	c := &Client{
		webhook:    webhook,
		httpClient: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// PostMessage marshals the Message and sends to the MS Teams webhook. If the
// webhook responds with a non-successful HTTP status code a non-nil error value
// will be returned.
func (c *Client) PostMessage(ctx context.Context, msg Message) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.webhook, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("teams: POST webhook returned unsuccessful status code: %d, body: %s",
			resp.StatusCode, string(respBody))
	}

	return nil
}

// PostMessage marshals the Message and sends to the MS Teams webhook using the
// default http.Client. If the webhook responds with a non-successful HTTP status
// code a non-nil error value will be returned.
func PostMessage(ctx context.Context, webhook string, msg Message) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhook, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("teams: POST webhook returned unsuccessful status code: %d, body: %s",
			resp.StatusCode, string(respBody))
	}

	return nil
}
