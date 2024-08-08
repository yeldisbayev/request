package req

import "net/http"

type Client interface {
	Request() Request
}

type client struct {
	httpClient *http.Client
}

func NewClient(
	httpClient *http.Client,
) Client {
	return &client{
		httpClient: httpClient,
	}

}

func (c *client) Request() Request {
	return &request{
		header: make(http.Header),
		client: c,
	}

}
