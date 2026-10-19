# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```sh
go test -race -cover ./...                         # all tests
go test -run TestPostPayloadRedactsWebhookFromErrors ./...   # one test (regex; subtests: 'TestName/sub_name')
go vet ./...                                       # includes stdversion: catches stdlib APIs newer than go.mod's go 1.18
golangci-lint run --enable-only=staticcheck,errcheck,govet,ineffassign,unused,revive,gosec ./...
```

## Constraints

- **Zero dependencies** is a stated goal (README). Standard library only, in tests too.
- **Minimum Go is 1.18** (`go.mod`). Don't use newer stdlib APIs such as `errors.Join`, `min`/`max`, `slices`, or `maps`. `go vet` flags them.
- The package name is `teams` but the module path ends in `teamsnotifications`. External code and `example_test.go` import it as `teams "github.com/jkratz55/teamsnotifications"`.
- The module is pre-1.0, so breaking changes are allowed but must be recorded in `CHANGELOG.md` under the next version heading (`## vX.Y.Z`).

## Architecture

There's one public operation: `Client.PostPayload` (`teams.go`) marshals any value to JSON and POSTs it. The package-level `PostPayload` is a thin wrapper over `New(webhook).PostPayload`, so make behavior changes in the method only.

**Teams payload shape.** Teams webhooks (Workflows) only render Adaptive Cards wrapped in an envelope: `{"type":"message","attachments":[{"contentType":"application/vnd.microsoft.card.adaptive","content":<card>}]}`. `NewAdaptiveCardMessage` builds it. Posting a bare `AdaptiveCard` gets a 202 response but never appears in the channel, so don't "simplify" the envelope away.

**Card types (`types.go`) follow one pattern** that every new element or action should copy:
- No `Type` field. `MarshalJSON` injects `"type"` by converting to a local method-less alias type (avoiding recursion) and embedding it in an anonymous struct after a `Type` field, so `"type"` comes first.
- Nil slices are replaced with empty slices in `MarshalJSON` so they encode as `[]`, not `null`. `AdaptiveCard` also defaults `Version` and `Schema`.
- `AdaptiveCardElement` and `AdaptiveCardAction` are sealed interfaces (unexported marker methods). `RawAdaptiveCardElement` and `RawAdaptiveCardAction` are the escape hatch for unmodeled JSON. The package intentionally does not cover the full Adaptive Card schema.
- Enum-like fields use typed string constants (`TextWeight`, `TextSize`, `TextColor`), so untyped string literals still compile.

**Error handling invariants:**
- Webhook URLs are credentials: Workflows puts a `sig=` query param in them, and legacy connectors put the secret in the path. Every error that could carry the URL goes through `redact()` (`redact.go`). It rewrites the message down to scheme and host, and keeps the original error reachable via `Unwrap` so `errors.Is` and `errors.As` still work. Route any new error path that involves the request or URL through it.
- A non-2xx response returns `*StatusError`. Only the first `maxErrorBodySize` (4 KiB) of the body is read, and the rest is drained up to `maxDrainSize` so the connection can be reused.
- The default client is the package-level `defaultHTTPClient` with a 30-second timeout. Never use `http.DefaultClient`, which has no timeout.

## Tests

- `teams_test.go` runs HTTP behavior against both entry points via the shared `posters` table. Add new request-level tests to that table.
- Request bodies are compared with `assertJSONEqual` (semantic). `types_test.go` compares exact marshaled strings on purpose, to pin field order and defaults.
- `httpClientFunc` adapts a function to `HTTPClient` for tests that fake transport behavior.
