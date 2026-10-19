package teams

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// defaultTimeout bounds requests sent with the default HTTP client so a
	// stalled connection can't block forever when the context has no deadline.
	defaultTimeout = 30 * time.Second

	// maxErrorBodySize caps how much of an unsuccessful response body is read
	// into a StatusError.
	maxErrorBodySize = 4 << 10

	// maxDrainSize caps how much of a response body is discarded so the
	// connection can be reused. Anything larger is abandoned instead.
	maxDrainSize = 64 << 10
)

var defaultHTTPClient = &http.Client{Timeout: defaultTimeout}

// HTTPClient is the interface used to send webhook requests. *http.Client
// satisfies it.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// StatusError is returned when the webhook responds with a non-2xx status code.
// Use errors.As to inspect it, for example to back off when Teams throttles
// requests with a 429.
type StatusError struct {
	StatusCode int
	// Body holds at most the first 4 KiB of the response body.
	Body string
}

func (e *StatusError) Error() string {
	return fmt.Sprintf("teams: POST webhook returned unsuccessful status code: %d, body: %s",
		e.StatusCode, e.Body)
}

// Client is a type that is capable of posting messages/notifications to a channel
// in Microsoft Teams using webhooks.
type Client struct {
	webhook    string
	httpClient HTTPClient
}

// New initializes a new Client. Unless WithHTTPClient is given, requests are
// sent with an http.Client that times out after 30 seconds.
func New(webhook string, opts ...Option) *Client {
	c := &Client{
		webhook:    webhook,
		httpClient: defaultHTTPClient,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// PostPayload marshals any payload value and sends it to the MS Teams webhook.
// Post an AdaptiveCard directly or wrapped with NewAdaptiveCardMessage,
// depending on what the webhook expects.
//
// A 2xx response means the webhook accepted the request. Workflows webhooks
// respond with 202 Accepted before the workflow runs, so a nil error does not
// guarantee the message was posted to the channel. Any other status code
// returns a *StatusError.
//
// Webhook URLs embed their credentials, so the URL is redacted from returned
// errors.
func (c *Client) PostPayload(ctx context.Context, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("teams: marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.webhook, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("teams: create request: %w", redact(err, c.webhook))
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("teams: send request: %w", redact(err, c.webhook, req.URL.String()))
	}
	defer func() {
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, maxDrainSize))
		_ = resp.Body.Close()
	}()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return nil
	}

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, maxErrorBodySize))
	statusErr := &StatusError{StatusCode: resp.StatusCode, Body: string(respBody)}
	if err != nil {
		return fmt.Errorf("%w (read response body: %v)", statusErr, err)
	}
	return statusErr
}

// PostPayload marshals any payload value and sends it to the MS Teams webhook
// using an http.Client that times out after 30 seconds. It behaves the same as
// Client.PostPayload.
func PostPayload(ctx context.Context, webhook string, payload any) error {
	return New(webhook).PostPayload(ctx, payload)
}
