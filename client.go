package req

import (
	"net/http"
	"net/url"
)

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
		client: c,
		header: make(http.Header),
		query:  make(url.Values),
	}

}
