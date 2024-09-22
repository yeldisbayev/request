package req

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	type args struct {
		options []func(*client)
	}

	type want struct {
		client *client
	}

	type test struct {
		name string
		args args
		want want
	}

	tests := []test{
		{
			name: "Without any options",
			args: args{},
			want: want{
				client: &client{
					httpClient: &http.Client{
						Transport: http.DefaultTransport,
					},
					timeout:                   DefaultTimeout,
					idleConnectionTimeout:     DefaultIdleConnectionTimeout,
					maxIdleConnections:        DefaultMaxIdleConnections,
					maxConnectionsPerHost:     DefaultMaxConnectionsPerHost,
					maxIdleConnectionsPerHost: DefaultMaxIdleConnectionsPerHost,
					forceAttemptHTTP2:         DefaultForceAttemptHTTP2,
				},
			},
		},
		{
			name: "WithTimeout option",
			args: args{
				options: []func(*client){
					WithTimeout(time.Second),
				},
			},
			want: want{
				client: &client{
					httpClient: &http.Client{
						Transport: http.DefaultTransport,
					},
					timeout:                   time.Second,
					idleConnectionTimeout:     DefaultIdleConnectionTimeout,
					maxIdleConnections:        DefaultMaxIdleConnections,
					maxConnectionsPerHost:     DefaultMaxConnectionsPerHost,
					maxIdleConnectionsPerHost: DefaultMaxIdleConnectionsPerHost,
					forceAttemptHTTP2:         DefaultForceAttemptHTTP2,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := NewClient(tc.args.options...)

			assert.Equal(t, tc.want.client, client)

		})
	}
}

func TestClient_Request(t *testing.T) {
	type want struct {
		req Request
	}

	type test struct {
		name string
		want want
	}

	tests := []test{
		{
			name: "Client_Request",
			want: want{
				req: &request{
					client: &client{
						httpClient: http.DefaultClient,
					},
					header: make(http.Header),
					query:  make(url.Values),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &client{
				httpClient: http.DefaultClient,
			}

			assert.Equal(t, tc.want.req, client.Request())

		})
	}
}

func TestWithTimeout(t *testing.T) {
	type args struct {
		timeout time.Duration
	}

	type want struct {
		timeout time.Duration
	}

	type test struct {
		name string
		args args
		want want
	}

	tests := []test{
		{
			name: "WithTimeout",
			args: args{
				timeout: time.Second,
			},
			want: want{
				timeout: time.Second,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &client{}

			WithTimeout(tc.args.timeout)(c)

			assert.Equal(t, tc.want.timeout, c.timeout)

		})
	}

}

func TestWithIdleConnectionTimeout(t *testing.T) {
	type args struct {
		idleConnectionTimeout time.Duration
	}

	type want struct {
		idleConnectionTimeout time.Duration
	}

	type test struct {
		name string
		args args
		want want
	}

	tests := []test{
		{
			name: "WithIdleConnectionTimeout",
			args: args{
				idleConnectionTimeout: time.Second,
			},
			want: want{
				idleConnectionTimeout: time.Second,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &client{}

			WithIdleConnectionTimeout(tc.args.idleConnectionTimeout)(c)

			assert.Equal(t, tc.want.idleConnectionTimeout, c.idleConnectionTimeout)

		})
	}

}

func TestWithMaxIdleConnections(t *testing.T) {
	type args struct {
		maxIdleConnections int
	}

	type want struct {
		maxIdleConnections int
	}

	type test struct {
		name string
		args args
		want want
	}

	tests := []test{
		{
			name: "WithMaxIdleConnections",
			args: args{
				maxIdleConnections: 10,
			},
			want: want{
				maxIdleConnections: 10,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &client{}

			WithMaxIdleConnections(tc.args.maxIdleConnections)(c)

			assert.Equal(t, tc.want.maxIdleConnections, c.maxIdleConnections)

		})
	}

}

func TestWithMaxConnectionsPerHost(t *testing.T) {
	type args struct {
		maxConnectionsPerHost int
	}

	type want struct {
		maxConnectionsPerHost int
	}

	type test struct {
		name string
		args args
		want want
	}

	tests := []test{
		{
			name: "WithMaxConnectionsPerHost",
			args: args{
				maxConnectionsPerHost: 10,
			},
			want: want{
				maxConnectionsPerHost: 10,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &client{}

			WithMaxConnectionsPerHost(tc.args.maxConnectionsPerHost)(c)

			assert.Equal(t, tc.want.maxConnectionsPerHost, c.maxConnectionsPerHost)

		})
	}

}

func TestWithMaxIdleConnectionsPerHost(t *testing.T) {
	type args struct {
		maxIdleConnectionsPerHost int
	}

	type want struct {
		maxIdleConnectionsPerHost int
	}

	type test struct {
		name string
		args args
		want want
	}

	tests := []test{
		{
			name: "WithMaxConnectionsPerHost",
			args: args{
				maxIdleConnectionsPerHost: 10,
			},
			want: want{
				maxIdleConnectionsPerHost: 10,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &client{}

			WithMaxIdleConnectionsPerHost(tc.args.maxIdleConnectionsPerHost)(c)

			assert.Equal(t, tc.want.maxIdleConnectionsPerHost, c.maxIdleConnectionsPerHost)

		})
	}

}

func TestWithForceAttemptHTTP2(t *testing.T) {
	type args struct {
		forceAttemptHTTP2 bool
	}

	type want struct {
		forceAttemptHTTP2 bool
	}

	type test struct {
		name string
		args args
		want want
	}

	tests := []test{
		{
			name: "WithMaxConnectionsPerHost",
			args: args{
				forceAttemptHTTP2: true,
			},
			want: want{
				forceAttemptHTTP2: true,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &client{}

			WithForceAttemptHTTP2(tc.args.forceAttemptHTTP2)(c)

			assert.Equal(t, tc.want.forceAttemptHTTP2, c.forceAttemptHTTP2)

		})
	}

}

func TestWithInterceptors(t *testing.T) {
	type args struct {
		interceptors []Interceptor
	}

	type want struct {
		roundTripper http.RoundTripper
	}

	type test struct {
		name string
		args args
		want want
	}

	tests := []test{
		{
			name: "WithInterceptors",
			args: args{
				interceptors: []Interceptor{
					func(tripper http.RoundTripper) http.RoundTripper {
						return RoundTripper(
							func(req *http.Request) (*http.Response, error) {
								return tripper.RoundTrip(req)
							},
						)
					},
				},
			},
			want: want{
				roundTripper: func(tripper http.RoundTripper) http.RoundTripper {
					return RoundTripper(
						func(req *http.Request) (*http.Response, error) {
							return tripper.RoundTrip(req)
						},
					)
				}(http.DefaultTransport),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &client{
				httpClient: &http.Client{
					Transport: nil,
				},
			}

			WithInterceptors(tc.args.interceptors...)(c)

			assert.NotNil(t, c.httpClient.Transport)

		})
	}
}
