package teams

import (
	"errors"
	"net/url"
	"sort"
	"strings"
)

// redactedError replaces the message of an error that contained a webhook URL
// while keeping the original available to errors.Is and errors.As.
type redactedError struct {
	msg string
	err error
}

func (e *redactedError) Error() string { return e.msg }

func (e *redactedError) Unwrap() error { return e.err }

// redact removes the given URLs, and the URL of any *url.Error in err's chain,
// from err's message. Teams webhook URLs carry their credentials in the query
// (Workflows' sig parameter) or the path (legacy connectors), and net/http
// includes the full URL in its errors, so only the scheme and host are kept.
func redact(err error, urls ...string) error {
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		urls = append(urls, urlErr.URL)
	}
	// Replace longer URLs first so a URL that prefixes another can't leave the
	// remainder, such as the query string, behind.
	sort.Slice(urls, func(i, j int) bool { return len(urls[i]) > len(urls[j]) })

	msg := err.Error()
	redacted := msg
	for _, u := range urls {
		if u != "" {
			redacted = strings.ReplaceAll(redacted, u, redactURL(u))
		}
	}
	if redacted == msg {
		return err
	}
	return &redactedError{msg: redacted, err: err}
}

func redactURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" {
		return "[redacted]"
	}
	return u.Scheme + "://" + u.Host + "/[redacted]"
}
