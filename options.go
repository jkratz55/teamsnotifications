package teams

// Option configures a Client.
type Option func(*Client)

// WithHTTPClient sets the HTTP client used to send requests. A nil httpClient
// is ignored and the default client is kept.
func WithHTTPClient(httpClient HTTPClient) Option {
	return func(c *Client) {
		if httpClient != nil {
			c.httpClient = httpClient
		}
	}
}
