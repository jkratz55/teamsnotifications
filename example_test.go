package teams_test

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	teams "github.com/jkratz55/teamsnotifications"
)

func Example() {
	card := teams.AdaptiveCard{
		Body: []teams.AdaptiveCardElement{
			teams.AdaptiveCardTextBlock{
				Text:   "Deployment: checkout-service",
				Weight: teams.WeightBolder,
				Size:   teams.SizeLarge,
				Color:  teams.ColorGood,
			},
			teams.AdaptiveCardFactSet{
				Facts: []teams.AdaptiveCardFact{
					{Title: "Environment", Value: "prod"},
					{Title: "Version", Value: "v2.18.4"},
				},
			},
		},
		Actions: []teams.AdaptiveCardAction{
			teams.AdaptiveCardOpenURLAction{Title: "View release", URL: "https://example.com/releases/v2.18.4"},
		},
		MSTeams: &teams.AdaptiveCardMSTeams{Width: teams.AdaptiveCardWidthFull},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// If the webhook's flow reads cards from an attachments array, post
	// teams.NewAdaptiveCardMessage(card) instead.
	err := teams.PostPayload(ctx, os.Getenv("TEAMS_WEBHOOK_URL"), card)
	var statusErr *teams.StatusError
	if errors.As(err, &statusErr) && statusErr.StatusCode == http.StatusTooManyRequests {
		log.Print("throttled by Teams, retry later")
	} else if err != nil {
		log.Print(err)
	}
}

func ExampleWithHTTPClient() {
	client := teams.New(os.Getenv("TEAMS_WEBHOOK_URL"),
		teams.WithHTTPClient(&http.Client{Timeout: 5 * time.Second}))

	card := teams.AdaptiveCard{
		Body: []teams.AdaptiveCardElement{teams.AdaptiveCardTextBlock{Text: "Hello from Go", Wrap: true}},
	}
	if err := client.PostPayload(context.Background(), card); err != nil {
		log.Print(err)
	}
}
