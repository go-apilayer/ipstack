package ipstack

import "net/http"

// HTTPClient is the interface used to send HTTP requests. Users can provide their own implementation.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Option is a functional option to modify the underlying Client.
type Option func(*Client)

// OptionHTTPClient - provide a custom http client to the client.
func OptionHTTPClient(client HTTPClient) func(*Client) {
	return func(c *Client) {
		c.client = client
	}
}

// OptionDebug enable debugging for the client.
func OptionDebug(b bool) func(*Client) {
	return func(c *Client) {
		c.debug = b
	}
}
