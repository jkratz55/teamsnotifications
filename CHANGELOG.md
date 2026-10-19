# Changelog

## v0.3.0

This release contains breaking API changes. See [Upgrading from v0.2.0](#upgrading-from-v020).

### Breaking changes

* Adaptive Cards must now be wrapped in a `Message` before posting, using `NewAdaptiveCardMessage`. Teams webhooks don't render a bare `AdaptiveCard`: the request is accepted but the card never appears in the channel.
* The `Type` fields have been removed from the Adaptive Card types. Each type now sets its `type` property when marshaled.
* `AdaptiveCard.Body` is now `[]AdaptiveCardElement` instead of `[]any`. Use `RawAdaptiveCardElement` for elements the package doesn't model.
* `AdaptiveCardTextBlock.Weight`, `Size` and `Color` are now the typed `TextWeight`, `TextSize` and `TextColor`. String literals still compile.
* Renamed `HttpClient` to `HTTPClient` and `WithHttpClient` to `WithHTTPClient`.
* Unsuccessful responses return a `*StatusError` instead of a plain error. Errors from encoding the payload, building the request and sending it are now prefixed with `teams: marshal payload:`, `teams: create request:` and `teams: send request:`.

### Added

* `Message`, `Attachment`, `NewAdaptiveCardMessage` and `AdaptiveCardContentType` for the Teams message envelope.
* `AdaptiveCardContainer`, `AdaptiveCardColumnSet` and `AdaptiveCardColumn` for layout.
* `AdaptiveCard.Actions` with `AdaptiveCardOpenURLAction`, and `RawAdaptiveCardAction` for other actions.
* `AdaptiveCard.MSTeams` (`AdaptiveCardMSTeams`) for Teams-specific properties, such as full-width cards with `AdaptiveCardWidthFull`.
* `StatusError`, exposing the status code and response body so callers can handle `429` throttling with `errors.As`.
* `AdaptiveCard.Version` and `Schema` default to `DefaultAdaptiveCardVersion` (`1.4`) and `AdaptiveCardSchema` when empty.
* Package documentation, runnable examples and a README usage guide.

### Fixed

* Webhook URLs are redacted from returned errors. Previously, transport errors included the full URL and the credentials in it, such as the Workflows `sig` parameter.
* Requests sent with the default HTTP client now time out after 30 seconds. They previously used `http.DefaultClient`, which never times out.
* Only the first 4 KiB of an unsuccessful response body is read into the error. Previously the whole body was read into memory and copied into the error message.
* Errors reading the response body are no longer silently ignored.
* `WithHTTPClient(nil)` is ignored instead of causing a panic on the first request.
* Nil `Body`, `Facts`, `Items` and `Columns` slices marshal as `[]` instead of `null`.

### Upgrading from v0.2.0

1. Remove `Type: "..."` from every card, text block and fact set.
2. Change `Body: []any{...}` to `Body: []teams.AdaptiveCardElement{...}`.
3. Wrap cards before posting: `PostPayload(ctx, card)` becomes `PostPayload(ctx, teams.NewAdaptiveCardMessage(card))`.
4. Rename `WithHttpClient` to `WithHTTPClient` and `HttpClient` to `HTTPClient`.
5. If you matched on error message text, use `errors.As(err, &statusErr)` with a `*teams.StatusError` instead.

```go
// v0.2.0
card := teams.AdaptiveCard{
	Type:    "AdaptiveCard",
	Version: "1.4",
	Schema:  "http://adaptivecards.io/schemas/adaptive-card.json",
	Body:    []any{teams.AdaptiveCardTextBlock{Type: "TextBlock", Text: "Deployed", Weight: "Bolder"}},
}
err := client.PostPayload(ctx, card)

// v0.3.0
card := teams.AdaptiveCard{
	Body: []teams.AdaptiveCardElement{teams.AdaptiveCardTextBlock{Text: "Deployed", Weight: teams.WeightBolder}},
}
err := client.PostPayload(ctx, teams.NewAdaptiveCardMessage(card))
```

## v0.2.0

* Replace `PostMessage` with `PostPayload` and replace the legacy message types with Adaptive Card types. This is a breaking API change.

## v0.1.0

* Initial release
