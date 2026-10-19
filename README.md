# teamsnotifications

A very simple zero dependency library for posting JSON payloads to Microsoft Teams through webhooks. It includes a few basic types for Adaptive Cards, but does not aim to support every Adaptive Card feature. This covers basic use cases such as sending notifications in a CI/CD pipeline, from an application, Kubernetes controller, or some other automation component. The main reason I published this as a library is I often find myself needing this feature and copy-pasting the code between code bases.

## Installation

```sh
go get github.com/jkratz55/teamsnotifications
```

The package name is `teams`, so importing it with an explicit name keeps that obvious:

```go
import teams "github.com/jkratz55/teamsnotifications"
```

## Usage

Teams webhooks expect Adaptive Cards to be wrapped in a message envelope. Build the card, then wrap it with `NewAdaptiveCardMessage`:

```go
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

err := teams.PostPayload(ctx, webhookURL, teams.NewAdaptiveCardMessage(card))
```

The `type` property of each element is set for you, and `version` and `$schema` default to `1.4` and the Adaptive Card schema. For elements or actions this package doesn't model, such as `Image`, pass their JSON as a `RawAdaptiveCardElement` or `RawAdaptiveCardAction`. You can also post any other JSON-serializable payload your webhook expects.

The library provides two ways of posting payloads to Teams.

1. Creating a `Client` and using `PostPayload` on the `Client`
2. Using the `PostPayload` function in the package

There is very little difference between the two methods, so which one you use will often depend on personal preference. The `Client` is initialized with the webhook, so you do not need to pass the webhook as an argument when calling `PostPayload` on the `Client` type. The `Client` type also takes options if you want to customize the HTTP client with `WithHTTPClient`. Meanwhile, the package-level `PostPayload` function requires you to pass the webhook each time and always uses the default HTTP client, which times out after 30 seconds.

## Errors

A non-2xx response returns a `*StatusError` holding the status code and the first 4 KiB of the response body:

```go
var statusErr *teams.StatusError
if errors.As(err, &statusErr) && statusErr.StatusCode == http.StatusTooManyRequests {
	// back off and retry
}
```

Webhook URLs contain credentials, so they are redacted from returned errors. Keep them out of source control and logs as well.

## Things to know about Teams webhooks

* Office 365 Connector incoming webhooks are retired. Create webhooks with the Teams Workflows app instead ("Send webhook alerts to a channel" template).
* Workflows webhooks respond with `202 Accepted` before the workflow runs, so a nil error means the request was accepted, not that the message was posted. Check the workflow's run history if a message doesn't show up.
* Messages are limited to 28 KB.
* Teams throttles webhooks that receive more than four requests per second. Retry `429` responses with exponential backoff.
