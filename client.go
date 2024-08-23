package req

import (
	"net/http"
	"net/url"
	"time"
)

const (
	DefaultTimeout                   = 30 * time.Second
	DefaultIdleConnectionTimeout     = 90 * time.Second
	DefaultMaxIdleConnections        = 100
	DefaultMaxConnectionsPerHost     = 15
	DefaultMaxIdleConnectionsPerHost = 10
	DefaultForceAttemptHTTP2         = false
)

type Client interface {
	Request() Request
}

type client struct {
	httpClient                *http.Client
	timeout                   time.Duration
	idleConnectionTimeout     time.Duration
	maxIdleConnections        int
	maxIdleConnectionsPerHost int
	maxConnectionsPerHost     int
	forceAttemptHTTP2         bool
}

func NewClient(
	options ...func(*client),
) Client {
	httpClient := http.DefaultClient
	transport := httpClient.Transport.(*http.Transport)

	client := &client{
		timeout:                   DefaultTimeout,
		idleConnectionTimeout:     DefaultIdleConnectionTimeout,
		maxIdleConnections:        DefaultMaxIdleConnections,
		maxConnectionsPerHost:     DefaultMaxConnectionsPerHost,
		maxIdleConnectionsPerHost: DefaultMaxIdleConnectionsPerHost,
		forceAttemptHTTP2:         DefaultForceAttemptHTTP2,
	}

	for _, option := range options {
		option(client)
	}

	transport.MaxIdleConns = client.maxIdleConnectionsPerHost
	transport.MaxConnsPerHost = client.maxIdleConnections
	transport.MaxIdleConnsPerHost = client.maxIdleConnections
	transport.IdleConnTimeout = client.idleConnectionTimeout
	transport.ForceAttemptHTTP2 = client.forceAttemptHTTP2

	return client

}

func (c *client) Request() Request {
	return &request{
		client: c,
		header: make(http.Header),
		query:  make(url.Values),
	}

}

// WithTimeout sets timeout for all client requests.
// Timeout implemented without using http.Client's Timeout property,
// but with context. Client timeout has lesser priority than Request timeout property.
// If not provided, DefaultTimeout is used.
func WithTimeout(timeout time.Duration) func(*client) {
	return func(c *client) {
		c.timeout = timeout
	}

}

// WithIdleConnectionTimeout controls amount of time an idle
// (keep-alive) connection will remain idle before closing itself.
// If not provided, DefaultIdleConnectionTimeout is used.
func WithIdleConnectionTimeout(idleConnectionTimeout time.Duration) func(*client) {
	return func(c *client) {
		c.idleConnectionTimeout = idleConnectionTimeout
	}

}

// WithMaxIdleConnections controls the maximum number of idle (keep-alive)
// connections across all hosts. If not provided, DefaultMaxIdleConnections is used.
func WithMaxIdleConnections(maxIdleConnections int) func(*client) {
	return func(c *client) {
		c.maxIdleConnections = maxIdleConnections
	}

}

// WithMaxConnectionsPerHost optionally limits the total number of
// connections per host, including connections in the dialing,
// active, and idle states. On limit violation, dials will block.
// If not provided, DefaultMaxConnectionsPerHost is used.
func WithMaxConnectionsPerHost(maxConnectionsPerHost int) func(*client) {
	return func(c *client) {
		c.maxConnectionsPerHost = maxConnectionsPerHost
	}

}

// WithMaxOpenIdleConnectionsPerHost controls the maximum idle
// (keep-alive) connections to keep per-host.
// If not provided, DefaultMaxIdleConnectionsPerHost is used.
func WithMaxOpenIdleConnectionsPerHost(maxOpenIdleConnections int) func(*client) {
	return func(c *client) {
		c.maxIdleConnectionsPerHost = maxOpenIdleConnections
	}

}

// WithForceAttemptHTTP2 controls whether HTTP/2 is enabled when a non-zero
// Dial, DialTLS, or DialContext func or TLSClientConfig is provided.
// By default, use of any those fields conservatively disables HTTP/2.
// To use a custom dialer or TLS config and still attempt HTTP/2
// upgrades, set this to true. If not provided, DefaultForceAttemptHTTP2 is used.
func WithForceAttemptHTTP2(forceAttemptHTTP2 bool) func(*client) {
	return func(c *client) {
		c.forceAttemptHTTP2 = forceAttemptHTTP2
	}

}
