// Package teams posts JSON payloads to Microsoft Teams through incoming
// webhooks, such as those created with the Teams Workflows app.
//
// Teams expects Adaptive Cards to be wrapped in a message envelope, which
// NewAdaptiveCardMessage builds:
//
//	card := teams.AdaptiveCard{
//		Body: []teams.AdaptiveCardElement{
//			teams.AdaptiveCardTextBlock{Text: "Deployment finished", Weight: teams.WeightBolder},
//		},
//	}
//	err := teams.PostPayload(ctx, webhookURL, teams.NewAdaptiveCardMessage(card))
//
// Webhook URLs contain credentials. Keep them out of source control and logs.
// Errors returned by this package have the URL redacted.
package teams
