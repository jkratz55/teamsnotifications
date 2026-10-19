// Package teams posts JSON payloads to Microsoft Teams through incoming
// webhooks, such as those created with the Teams Workflows app.
//
// The payload format depends on the webhook. A Workflows webhook runs a flow,
// and the flow decides how to read the request body. A flow that passes the
// body straight to a "Post card" action expects a bare AdaptiveCard. A flow
// that reads cards from an attachments array, the format Microsoft documents
// for incoming webhooks, expects the card wrapped in a Message, which
// NewAdaptiveCardMessage builds:
//
//	card := teams.AdaptiveCard{
//		Body: []teams.AdaptiveCardElement{
//			teams.AdaptiveCardTextBlock{Text: "Deployment finished", Weight: teams.WeightBolder},
//		},
//	}
//	err := teams.PostPayload(ctx, webhookURL, card)                              // bare card
//	err = teams.PostPayload(ctx, webhookURL, teams.NewAdaptiveCardMessage(card)) // wrapped
//
// Webhook URLs contain credentials. Keep them out of source control and logs.
// Errors returned by this package have the URL redacted.
package teams
