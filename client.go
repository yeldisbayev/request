package req

import (
	"net/http"
	"net/url"
	"time"
)

type Client interface {
	Request() Request

	WithTimeout(
		timeout time.Duration,
	) Client
}

type client struct {
	httpClient *http.Client
	timeout    time.Duration
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

func (c *client) WithTimeout(
	timeout time.Duration,
) Client {
	c.timeout = timeout

	return c

}
