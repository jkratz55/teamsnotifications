package teams

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPostPayloadSendsAdaptiveCard(t *testing.T) {
	const expectedBody = `{"type":"AdaptiveCard","version":"1.4","body":[{"type":"TextBlock","text":"testEventFromTerminal:checkout-service","weight":"Bolder","size":"Large","color":"good"},{"type":"TextBlock","text":"checkout-servicedeployedtoproduction","wrap":true},{"type":"FactSet","facts":[{"title":"Type","value":"deployment"},{"title":"Environment","value":"prod"},{"title":"Name","value":"checkout-service"},{"title":"Timestamp(USEST)","value":"2026-10-0711:30:00EDT"},{"title":"Timestamp(BLR)","value":"2026-10-0721:00:00IST"},{"title":"Timestamp(UTC)","value":"2026-10-0715:30:00UTC"},{"title":"Region","value":"eastus2"},{"title":"Version","value":"v2.18.4"},{"title":"ChangeRequest","value":"CHG-104299"}]}],"$schema":"http://adaptivecards.io/schemas/adaptive-card.json"}`

	posters := []struct {
		name string
		post func(string, any) error
	}{
		{
			name: "client method",
			post: func(webhook string, payload any) error {
				return New(webhook).PostPayload(context.Background(), payload)
			},
		},
		{
			name: "package function",
			post: func(webhook string, payload any) error {
				return PostPayload(context.Background(), webhook, payload)
			},
		},
	}

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
				if got := string(body); got != expectedBody {
					t.Errorf("request body = %s, want %s", got, expectedBody)
				}
				w.WriteHeader(http.StatusNoContent)
			}))
			defer server.Close()

			if err := poster.post(server.URL, testAdaptiveCard()); err != nil {
				t.Fatalf("PostPayload() error = %v", err)
			}
		})
	}
}

func TestPostPayloadReturnsErrorForUnsuccessfulResponse(t *testing.T) {
	posters := []struct {
		name string
		post func(string, any) error
	}{
		{
			name: "client method",
			post: func(webhook string, payload any) error {
				return New(webhook).PostPayload(context.Background(), payload)
			},
		},
		{
			name: "package function",
			post: func(webhook string, payload any) error {
				return PostPayload(context.Background(), webhook, payload)
			},
		},
	}

	for _, poster := range posters {
		t.Run(poster.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusBadGateway)
				_, _ = io.WriteString(w, "webhook unavailable")
			}))
			defer server.Close()

			err := poster.post(server.URL, testAdaptiveCard())
			if err == nil {
				t.Fatal("PostPayload() error = nil, want unsuccessful status error")
			}
			if !strings.Contains(err.Error(), "502") || !strings.Contains(err.Error(), "webhook unavailable") {
				t.Errorf("PostPayload() error = %q, want status and response body", err)
			}
		})
	}
}

func testAdaptiveCard() AdaptiveCard {
	return AdaptiveCard{
		Type:    "AdaptiveCard",
		Version: "1.4",
		Body: []any{
			AdaptiveCardTextBlock{
				Type:   "TextBlock",
				Text:   "testEventFromTerminal:checkout-service",
				Weight: "Bolder",
				Size:   "Large",
				Color:  "good",
			},
			AdaptiveCardTextBlock{
				Type: "TextBlock",
				Text: "checkout-servicedeployedtoproduction",
				Wrap: true,
			},
			AdaptiveCardFactSet{
				Type: "FactSet",
				Facts: []AdaptiveCardFact{
					{Title: "Type", Value: "deployment"},
					{Title: "Environment", Value: "prod"},
					{Title: "Name", Value: "checkout-service"},
					{Title: "Timestamp(USEST)", Value: "2026-10-0711:30:00EDT"},
					{Title: "Timestamp(BLR)", Value: "2026-10-0721:00:00IST"},
					{Title: "Timestamp(UTC)", Value: "2026-10-0715:30:00UTC"},
					{Title: "Region", Value: "eastus2"},
					{Title: "Version", Value: "v2.18.4"},
					{Title: "ChangeRequest", Value: "CHG-104299"},
				},
			},
		},
		Schema: "http://adaptivecards.io/schemas/adaptive-card.json",
	}
}
