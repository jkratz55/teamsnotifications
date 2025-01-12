package teams

type Option func(*Client)

func WithHttpClient(httpClient HttpClient) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}
