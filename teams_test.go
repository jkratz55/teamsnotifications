package teams

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

// posters runs a test against both the Client method and the package function.
var posters = []struct {
	name string
	post func(ctx context.Context, webhook string, payload any) error
}{
	{
		name: "client method",
		post: func(ctx context.Context, webhook string, payload any) error {
			return New(webhook).PostPayload(ctx, payload)
		},
	},
	{
		name: "package function",
		post: PostPayload,
	},
}

type httpClientFunc func(*http.Request) (*http.Response, error)

func (f httpClientFunc) Do(req *http.Request) (*http.Response, error) { return f(req) }

func TestPostPayloadSendsAdaptiveCardMessage(t *testing.T) {
	const expectedBody = `{
		"type": "message",
		"attachments": [{
			"contentType": "application/vnd.microsoft.card.adaptive",
			"content": {
				"type": "AdaptiveCard",
				"version": "1.4",
				"body": [
					{"type": "TextBlock", "text": "Test Event From Terminal: checkout-service", "weight": "Bolder", "size": "Large", "color": "Good"},
					{"type": "TextBlock", "text": "checkout-service deployed to production", "wrap": true},
					{"type": "FactSet", "facts": [
						{"title": "Type", "value": "deployment"},
						{"title": "Environment", "value": "prod"},
						{"title": "Name", "value": "checkout-service"},
						{"title": "Timestamp (US EST)", "value": "2026-10-07 11:30:00 EDT"},
						{"title": "Timestamp (BLR)", "value": "2026-10-07 21:00:00 IST"},
						{"title": "Timestamp (UTC)", "value": "2026-10-07 15:30:00 UTC"},
						{"title": "Region", "value": "eastus2"},
						{"title": "Version", "value": "v2.18.4"},
						{"title": "Change Request", "value": "CHG-104299"}
					]}
				],
				"actions": [{"type": "Action.OpenUrl", "title": "View release", "url": "https://example.com/releases/v2.18.4"}],
				"msteams": {"width": "Full"},
				"$schema": "http://adaptivecards.io/schemas/adaptive-card.json"
			}
		}]
	}`

	for _, poster := range posters {
		t.Run(poster.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("request method = %q, want %q", r.Method, http.MethodPost)
				}
				if got := r.Header.Get("Content-Type"); got != "application/json" {
					t.Errorf("Content-Type = %q, want %q", got, "application/json")
				}

				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Errorf("read request body: %v", err)
				}
				assertJSONEqual(t, string(body), expectedBody)
				w.WriteHeader(http.StatusAccepted)
			}))
			defer server.Close()

			err := poster.post(context.Background(), server.URL, NewAdaptiveCardMessage(testAdaptiveCard()))
			if err != nil {
				t.Fatalf("PostPayload() error = %v", err)
			}
		})
	}
}

func TestPostPayloadReturnsStatusErrorForUnsuccessfulResponse(t *testing.T) {
	for _, poster := range posters {
		t.Run(poster.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusBadGateway)
				_, _ = io.WriteString(w, "webhook unavailable")
			}))
			defer server.Close()

			err := poster.post(context.Background(), server.URL, NewAdaptiveCardMessage(testAdaptiveCard()))
			var statusErr *StatusError
			if !errors.As(err, &statusErr) {
				t.Fatalf("PostPayload() error = %v, want *StatusError", err)
			}
			if statusErr.StatusCode != http.StatusBadGateway {
				t.Errorf("StatusCode = %d, want %d", statusErr.StatusCode, http.StatusBadGateway)
			}
			if statusErr.Body != "webhook unavailable" {
				t.Errorf("Body = %q, want %q", statusErr.Body, "webhook unavailable")
			}
			if !strings.Contains(err.Error(), "502") || !strings.Contains(err.Error(), "webhook unavailable") {
				t.Errorf("PostPayload() error = %q, want status and response body", err)
			}
		})
	}
}

func TestPostPayloadTruncatesErrorBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, strings.Repeat("x", 1<<20))
	}))
	defer server.Close()

	err := PostPayload(context.Background(), server.URL, nil)
	var statusErr *StatusError
	if !errors.As(err, &statusErr) {
		t.Fatalf("PostPayload() error = %v, want *StatusError", err)
	}
	if len(statusErr.Body) != maxErrorBodySize {
		t.Errorf("len(Body) = %d, want %d", len(statusErr.Body), maxErrorBodySize)
	}
}

func TestPostPayloadRedactsWebhookFromErrors(t *testing.T) {
	const secret = "SUPERSECRET"

	closed := httptest.NewServer(http.NotFoundHandler())
	closed.Close()

	tests := []struct {
		name    string
		webhook string
		opts    []Option
	}{
		{
			name:    "transport error",
			webhook: closed.URL + "/workflows/abc/triggers/manual/paths/invoke?api-version=1&sig=" + secret,
		},
		{
			name:    "secret in path",
			webhook: closed.URL + "/webhookb2/" + secret + "/IncomingWebhook/abc",
		},
		{
			name:    "invalid URL",
			webhook: "http://example.com/%zz?sig=" + secret,
		},
		{
			name:    "custom client wrapping url.Error",
			webhook: "https://example.com/workflows/abc?sig=" + secret,
			opts: []Option{WithHTTPClient(httpClientFunc(func(req *http.Request) (*http.Response, error) {
				return nil, fmt.Errorf("custom client: %w", &url.Error{Op: "Post", URL: req.URL.String(), Err: errors.New("boom")})
			}))},
		},
		{
			name:    "custom client with URL in message",
			webhook: "https://example.com/workflows/abc?sig=" + secret,
			opts: []Option{WithHTTPClient(httpClientFunc(func(req *http.Request) (*http.Response, error) {
				return nil, fmt.Errorf("failed to post to %s", req.URL)
			}))},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := New(tt.webhook, tt.opts...).PostPayload(context.Background(), nil)
			if err == nil {
				t.Fatal("PostPayload() error = nil, want error")
			}
			if strings.Contains(err.Error(), secret) {
				t.Errorf("PostPayload() error = %q, leaks the webhook secret", err)
			}
			if !strings.Contains(err.Error(), "[redacted]") {
				t.Errorf("PostPayload() error = %q, want redaction marker", err)
			}
		})
	}
}

func TestPostPayloadPreservesContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := PostPayload(ctx, server.URL+"?sig=secret", nil)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("PostPayload() error = %v, want context.Canceled", err)
	}
}

func TestPostPayloadReturnsMarshalError(t *testing.T) {
	err := PostPayload(context.Background(), "http://example.com", make(chan int))
	var unsupported *json.UnsupportedTypeError
	if !errors.As(err, &unsupported) {
		t.Errorf("PostPayload() error = %v, want *json.UnsupportedTypeError", err)
	}
}

func TestWithHTTPClient(t *testing.T) {
	var called bool
	client := httpClientFunc(func(*http.Request) (*http.Response, error) {
		called = true
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(""))}, nil
	})

	err := New("https://example.com/webhook", WithHTTPClient(client)).PostPayload(context.Background(), nil)
	if err != nil {
		t.Fatalf("PostPayload() error = %v", err)
	}
	if !called {
		t.Error("custom HTTP client was not used")
	}
}

func TestNewUsesDefaultHTTPClientWithTimeout(t *testing.T) {
	for name, c := range map[string]*Client{
		"no options":        New("https://example.com/webhook"),
		"nil client option": New("https://example.com/webhook", WithHTTPClient(nil)),
	} {
		t.Run(name, func(t *testing.T) {
			if c.httpClient != defaultHTTPClient {
				t.Errorf("httpClient = %v, want default client", c.httpClient)
			}
		})
	}
	if defaultHTTPClient.Timeout <= 0 {
		t.Errorf("default client Timeout = %v, want > 0", defaultHTTPClient.Timeout)
	}
}

func assertJSONEqual(t *testing.T, got, want string) {
	t.Helper()
	var gotValue, wantValue any
	if err := json.Unmarshal([]byte(got), &gotValue); err != nil {
		t.Fatalf("unmarshal got: %v\n%s", err, got)
	}
	if err := json.Unmarshal([]byte(want), &wantValue); err != nil {
		t.Fatalf("unmarshal want: %v\n%s", err, want)
	}
	if !reflect.DeepEqual(gotValue, wantValue) {
		t.Errorf("JSON mismatch\n got: %s\nwant: %s", got, want)
	}
}

func testAdaptiveCard() AdaptiveCard {
	return AdaptiveCard{
		Body: []AdaptiveCardElement{
			AdaptiveCardTextBlock{
				Text:   "Test Event From Terminal: checkout-service",
				Weight: WeightBolder,
				Size:   SizeLarge,
				Color:  ColorGood,
			},
			AdaptiveCardTextBlock{
				Text: "checkout-service deployed to production",
				Wrap: true,
			},
			AdaptiveCardFactSet{
				Facts: []AdaptiveCardFact{
					{Title: "Type", Value: "deployment"},
					{Title: "Environment", Value: "prod"},
					{Title: "Name", Value: "checkout-service"},
					{Title: "Timestamp (US EST)", Value: "2026-10-07 11:30:00 EDT"},
					{Title: "Timestamp (BLR)", Value: "2026-10-07 21:00:00 IST"},
					{Title: "Timestamp (UTC)", Value: "2026-10-07 15:30:00 UTC"},
					{Title: "Region", Value: "eastus2"},
					{Title: "Version", Value: "v2.18.4"},
					{Title: "Change Request", Value: "CHG-104299"},
				},
			},
		},
		Actions: []AdaptiveCardAction{
			AdaptiveCardOpenURLAction{Title: "View release", URL: "https://example.com/releases/v2.18.4"},
		},
		MSTeams: &AdaptiveCardMSTeams{Width: AdaptiveCardWidthFull},
	}
}
