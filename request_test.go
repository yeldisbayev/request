package req

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestRequest_URL(t *testing.T) {
	type want struct {
		url *url.URL
	}

	type depends struct {
		httpRequest *http.Request
	}

	type test struct {
		name    string
		want    want
		depends depends
	}

	tests := []test{
		{
			name: "Nil request",
			want: want{
				url: nil,
			},
			depends: depends{
				httpRequest: nil,
			},
		},
		{
			name: "Nil url",
			want: want{
				url: nil,
			},
			depends: depends{
				httpRequest: &http.Request{
					URL: nil,
				},
			},
		},
		{
			name: "Not nil url",
			want: want{
				url: &url.URL{
					Scheme: "https",
				},
			},
			depends: depends{
				httpRequest: &http.Request{
					URL: &url.URL{
						Scheme: "https",
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := request{
				httpReq: tc.depends.httpRequest,
			}

			assert.Equal(t, tc.want.url, req.URL())

		})
	}

}

func TestRequest_Body(t *testing.T) {
	type want struct {
		body io.ReadCloser
	}

	type depends struct {
		httpRequest *http.Request
	}

	type test struct {
		name    string
		want    want
		depends depends
	}

	tests := []test{
		{
			name: "Nil request",
			want: want{
				body: nil,
			},
			depends: depends{
				httpRequest: nil,
			},
		},
		{
			name: "Nil body",
			want: want{
				body: nil,
			},
			depends: depends{
				httpRequest: &http.Request{
					URL: nil,
				},
			},
		},
		{
			name: "Not nil body",
			want: want{
				body: io.NopCloser(bytes.NewBuffer([]byte(`Sample`))),
			},
			depends: depends{
				httpRequest: &http.Request{
					Body: io.NopCloser(bytes.NewBuffer([]byte(`Sample`))),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := request{
				httpReq: tc.depends.httpRequest,
			}

			assert.Equal(t, tc.want.body, req.Body())

		})
	}

}

func TestRequest_Header(t *testing.T) {
	type want struct {
		header http.Header
	}

	type depends struct {
		httpRequest *http.Request
	}

	type test struct {
		name    string
		want    want
		depends depends
	}

	tests := []test{
		{
			name: "Nil request",
			want: want{
				header: nil,
			},
			depends: depends{
				httpRequest: nil,
			},
		},
		{
			name: "Nil header",
			want: want{
				header: nil,
			},
			depends: depends{
				httpRequest: &http.Request{
					URL: nil,
				},
			},
		},
		{
			name: "Not nil header",
			want: want{
				header: http.Header{
					http.CanonicalHeaderKey("Content-Type"): {"text/plain; charset=utf-8"},
				},
			},
			depends: depends{
				httpRequest: &http.Request{
					Header: http.Header{
						http.CanonicalHeaderKey("Content-Type"): {"text/plain; charset=utf-8"},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := request{
				httpReq: tc.depends.httpRequest,
			}

			assert.Equal(t, tc.want.header, req.Header())

		})
	}

}

func TestRequest_WithHeader(t *testing.T) {
	type args struct {
		key    string
		values []string
	}

	type want struct {
		req *request
	}

	type depends struct {
		headers http.Header
	}

	type test struct {
		name    string
		args    args
		want    want
		depends depends
	}

	tests := []test{
		{
			name: "Without collision",
			args: args{
				key:    http.CanonicalHeaderKey("key"),
				values: []string{"value 1", "value 2", "value 3"},
			},
			want: want{
				req: &request{
					header: http.Header{
						http.CanonicalHeaderKey("key"): {"value 1", "value 2", "value 3"},
					},
				},
			},
			depends: depends{
				headers: make(http.Header),
			},
		},
		{
			name: "With collision",
			args: args{
				key:    http.CanonicalHeaderKey("key"),
				values: []string{"value 4", "value 5", "value 6"},
			},
			want: want{
				req: &request{
					header: http.Header{
						http.CanonicalHeaderKey("key"): {"value 1", "value 2", "value 3", "value 4", "value 5", "value 6"},
					},
				},
			},
			depends: depends{
				headers: http.Header{
					http.CanonicalHeaderKey("key"): {"value 1", "value 2", "value 3"},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := &request{
				header: tc.depends.headers,
			}

			assert.Equal(
				t,
				tc.want.req,
				req.WithHeader(
					tc.args.key,
					tc.args.values...,
				),
			)

		})
	}

}

func TestRequest_WithHeaders(t *testing.T) {
	type args struct {
		headers http.Header
	}

	type want struct {
		req *request
	}

	type depends struct {
		headers http.Header
	}

	type test struct {
		name    string
		args    args
		want    want
		depends depends
	}

	tests := []test{
		{
			name: "Without collision",
			args: args{
				headers: http.Header{
					http.CanonicalHeaderKey("key 2"): {"value 4", "value 5", "value 6"},
				},
			},
			want: want{
				req: &request{
					header: http.Header{
						http.CanonicalHeaderKey("key"):   {"value 1", "value 2", "value 3"},
						http.CanonicalHeaderKey("key 2"): {"value 4", "value 5", "value 6"},
					},
				},
			},
			depends: depends{
				headers: http.Header{
					http.CanonicalHeaderKey("key"): {"value 1", "value 2", "value 3"},
				},
			},
		},
		{
			name: "With collision",
			args: args{
				headers: http.Header{
					http.CanonicalHeaderKey("key"): {"value 4", "value 5", "value 6"},
				},
			},
			want: want{
				req: &request{
					header: http.Header{
						http.CanonicalHeaderKey("key"): {"value 4", "value 5", "value 6"},
					},
				},
			},
			depends: depends{
				headers: http.Header{
					http.CanonicalHeaderKey("key"): {"value 1", "value 2", "value 3"},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := &request{
				header: tc.depends.headers,
			}

			assert.Equal(
				t,
				tc.want.req,
				req.WithHeaders(tc.args.headers),
			)

		})
	}

}

func TestRequest_WithContentType(t *testing.T) {
	type args struct {
		value string
	}

	type want struct {
		req *request
	}

	type depends struct {
		headers http.Header
	}

	type test struct {
		name    string
		args    args
		want    want
		depends depends
	}

	tests := []test{
		{
			name: "Without collision",
			args: args{
				value: "application/json; charset=utf-8",
			},
			want: want{
				req: &request{
					header: http.Header{
						http.CanonicalHeaderKey(ContentType): {"application/json; charset=utf-8"},
					},
				},
			},
			depends: depends{
				headers: make(http.Header),
			},
		},
		{
			name: "With collision",
			args: args{
				value: "application/json; charset=utf-8",
			},
			want: want{
				req: &request{
					header: http.Header{
						http.CanonicalHeaderKey(ContentType): {"application/json; charset=utf-8"},
					},
				},
			},
			depends: depends{
				headers: http.Header{
					http.CanonicalHeaderKey(ContentType): {"application/xml; charset=utf-8"},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := &request{
				header: tc.depends.headers,
			}

			assert.Equal(
				t,
				tc.want.req,
				req.WithContentType(tc.args.value),
			)

		})
	}

}

func TestRequest_WithJSONContentType(t *testing.T) {
	type want struct {
		req *request
	}

	type depends struct {
		headers http.Header
	}

	type test struct {
		name    string
		want    want
		depends depends
	}

	tests := []test{
		{
			name: "Without collision",
			want: want{
				req: &request{
					header: http.Header{
						http.CanonicalHeaderKey(ContentType): {"application/json; charset=utf-8"},
					},
				},
			},
			depends: depends{
				headers: make(http.Header),
			},
		},
		{
			name: "With collision",
			want: want{
				req: &request{
					header: http.Header{
						http.CanonicalHeaderKey(ContentType): {"application/json; charset=utf-8"},
					},
				},
			},
			depends: depends{
				headers: http.Header{
					http.CanonicalHeaderKey(ContentType): {"application/xml; charset=utf-8"},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := &request{
				header: tc.depends.headers,
			}

			assert.Equal(
				t,
				tc.want.req,
				req.WithJSONContentType(),
			)

		})
	}

}

func TestRequest_WithXMLContentType(t *testing.T) {
	type want struct {
		req *request
	}

	type depends struct {
		headers http.Header
	}

	type test struct {
		name    string
		want    want
		depends depends
	}

	tests := []test{
		{
			name: "Without collision",
			want: want{
				req: &request{
					header: http.Header{
						http.CanonicalHeaderKey(ContentType): {"application/xml; charset=utf-8"},
					},
				},
			},
			depends: depends{
				headers: make(http.Header),
			},
		},
		{
			name: "With collision",
			want: want{
				req: &request{
					header: http.Header{
						http.CanonicalHeaderKey(ContentType): {"application/xml; charset=utf-8"},
					},
				},
			},
			depends: depends{
				headers: http.Header{
					http.CanonicalHeaderKey(ContentType): {"application/json; charset=utf-8"},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := &request{
				header: tc.depends.headers,
			}

			assert.Equal(
				t,
				tc.want.req,
				req.WithXMLContentType(),
			)

		})
	}

}

func TestRequest_WithFormContentType(t *testing.T) {
	type want struct {
		req *request
	}

	type depends struct {
		headers http.Header
	}

	type test struct {
		name    string
		want    want
		depends depends
	}

	tests := []test{
		{
			name: "Without collision",
			want: want{
				req: &request{
					header: http.Header{
						http.CanonicalHeaderKey(ContentType): {"application/x-www-form-urlencoded"},
					},
				},
			},
			depends: depends{
				headers: make(http.Header),
			},
		},
		{
			name: "With collision",
			want: want{
				req: &request{
					header: http.Header{
						http.CanonicalHeaderKey(ContentType): {"application/x-www-form-urlencoded"},
					},
				},
			},
			depends: depends{
				headers: http.Header{
					http.CanonicalHeaderKey(ContentType): {"multipart/form-data"},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := &request{
				header: tc.depends.headers,
			}

			assert.Equal(
				t,
				tc.want.req,
				req.WithFormContentType(),
			)

		})
	}

}

func TestRequest_WithMultipartFormContentType(t *testing.T) {
	type want struct {
		req *request
	}

	type depends struct {
		headers http.Header
	}

	type test struct {
		name    string
		want    want
		depends depends
	}

	tests := []test{
		{
			name: "Without collision",
			want: want{
				req: &request{
					header: http.Header{
						http.CanonicalHeaderKey(ContentType): {"multipart/form-data"},
					},
				},
			},
			depends: depends{
				headers: make(http.Header),
			},
		},
		{
			name: "With collision",
			want: want{
				req: &request{
					header: http.Header{
						http.CanonicalHeaderKey(ContentType): {"multipart/form-data"},
					},
				},
			},
			depends: depends{
				headers: http.Header{
					http.CanonicalHeaderKey(ContentType): {"application/x-www-form-urlencoded; charset=utf-8"},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := &request{
				header: tc.depends.headers,
			}

			assert.Equal(
				t,
				tc.want.req,
				req.WithMultipartFormContentType(),
			)

		})
	}

}

func TestRequest_WithAuth(t *testing.T) {
	type args struct {
		values []string
	}

	type want struct {
		req *request
	}

	type depends struct {
		headers http.Header
	}

	type test struct {
		name    string
		args    args
		want    want
		depends depends
	}

	tests := []test{
		{
			name: "Without collision",
			args: args{
				values: []string{"value 1", "value 2", "value 3"},
			},
			want: want{
				req: &request{
					header: http.Header{
						http.CanonicalHeaderKey(Authorization): {"value 1", "value 2", "value 3"},
					},
				},
			},
			depends: depends{
				headers: make(http.Header),
			},
		},
		{
			name: "With collision",
			args: args{
				values: []string{"value 4", "value 5", "value 6"},
			},
			want: want{
				req: &request{
					header: http.Header{
						http.CanonicalHeaderKey(Authorization): {"value 1", "value 2", "value 3", "value 4", "value 5", "value 6"},
					},
				},
			},
			depends: depends{
				headers: http.Header{
					http.CanonicalHeaderKey(Authorization): {"value 1", "value 2", "value 3"},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := &request{
				header: tc.depends.headers,
			}

			assert.Equal(
				t,
				tc.want.req,
				req.WithAuth(tc.args.values...),
			)

		})
	}

}

func TestRequest_WithBasicAuth(t *testing.T) {
	type args struct {
		username string
		password string
	}

	type want struct {
		req *request
	}

	type depends struct {
		headers http.Header
	}

	type test struct {
		name    string
		args    args
		want    want
		depends depends
	}

	tests := []test{
		{
			name: "Without collision",
			args: args{
				username: "username",
				password: "password",
			},
			want: want{
				req: &request{
					header: http.Header{
						http.CanonicalHeaderKey(Authorization): {
							fmt.Sprintf(
								"%s %s",
								Basic,
								base64.StdEncoding.EncodeToString(
									[]byte("username:password"),
								),
							),
						},
					},
				},
			},
			depends: depends{
				headers: make(http.Header),
			},
		},
		{
			name: "With collision",
			args: args{
				username: "username",
				password: "password",
			},
			want: want{
				req: &request{
					header: http.Header{
						http.CanonicalHeaderKey(Authorization): {
							fmt.Sprintf(
								"%s %s",
								Bearer,
								"abcdefghijklmnopqrstuvwxyzABCDEFGHUJKLMNOPQRSTUVWXYZ0123456789",
							),
							fmt.Sprintf(
								"%s %s",
								Basic,
								base64.StdEncoding.EncodeToString(
									[]byte("username:password"),
								),
							),
						},
					},
				},
			},
			depends: depends{
				headers: http.Header{
					http.CanonicalHeaderKey(Authorization): {
						fmt.Sprintf(
							"%s %s",
							Bearer,
							"abcdefghijklmnopqrstuvwxyzABCDEFGHUJKLMNOPQRSTUVWXYZ0123456789",
						),
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := &request{
				header: tc.depends.headers,
			}

			assert.Equal(
				t,
				tc.want.req,
				req.WithBasicAuth(
					tc.args.username,
					tc.args.password,
				),
			)
		})
	}

}

func TestRequest_WithBearerAuth(t *testing.T) {
	type args struct {
		value string
	}

	type want struct {
		req *request
	}

	type depends struct {
		headers http.Header
	}

	type test struct {
		name    string
		args    args
		want    want
		depends depends
	}

	tests := []test{
		{
			name: "Without collision",
			args: args{
				value: "abcdefghijklmnopqrstuvwxyzABCDEFGHUJKLMNOPQRSTUVWXYZ0123456789",
			},
			want: want{
				req: &request{
					header: http.Header{
						http.CanonicalHeaderKey(Authorization): {
							fmt.Sprintf(
								"%s %s",
								Bearer,
								"abcdefghijklmnopqrstuvwxyzABCDEFGHUJKLMNOPQRSTUVWXYZ0123456789",
							),
						},
					},
				},
			},
			depends: depends{
				headers: make(http.Header),
			},
		},
		{
			name: "With collision",
			args: args{
				value: "abcdefghijklmnopqrstuvwxyzABCDEFGHUJKLMNOPQRSTUVWXYZ0123456789",
			},
			want: want{
				req: &request{
					header: http.Header{
						http.CanonicalHeaderKey(Authorization): {
							fmt.Sprintf(
								"%s %s",
								Basic,
								base64.StdEncoding.EncodeToString(
									[]byte("username:password"),
								),
							),
							fmt.Sprintf(
								"%s %s",
								Bearer,
								"abcdefghijklmnopqrstuvwxyzABCDEFGHUJKLMNOPQRSTUVWXYZ0123456789",
							),
						},
					},
				},
			},
			depends: depends{
				headers: http.Header{
					http.CanonicalHeaderKey(Authorization): {
						fmt.Sprintf(
							"%s %s",
							Basic,
							base64.StdEncoding.EncodeToString(
								[]byte("username:password"),
							),
						),
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := &request{
				header: tc.depends.headers,
			}

			assert.Equal(
				t,
				tc.want.req,
				req.WithBearerAuth(
					tc.args.value,
				),
			)
		})
	}

}

func TestRequest_WithJWTAuth(t *testing.T) {
	type args struct {
		value string
	}

	type want struct {
		req *request
	}

	type depends struct {
		headers http.Header
	}

	type test struct {
		name    string
		args    args
		want    want
		depends depends
	}

	tests := []test{
		{
			name: "Without collision",
			args: args{
				value: "abcdefghijklmnopqrstuvwxyzABCDEFGHUJKLMNOPQRSTUVWXYZ0123456789",
			},
			want: want{
				req: &request{
					header: http.Header{
						http.CanonicalHeaderKey(Authorization): {
							fmt.Sprintf(
								"%s %s",
								JWT,
								"abcdefghijklmnopqrstuvwxyzABCDEFGHUJKLMNOPQRSTUVWXYZ0123456789",
							),
						},
					},
				},
			},
			depends: depends{
				headers: make(http.Header),
			},
		},
		{
			name: "With collision",
			args: args{
				value: "abcdefghijklmnopqrstuvwxyzABCDEFGHUJKLMNOPQRSTUVWXYZ0123456789",
			},
			want: want{
				req: &request{
					header: http.Header{
						http.CanonicalHeaderKey(Authorization): {
							fmt.Sprintf(
								"%s %s",
								Basic,
								base64.StdEncoding.EncodeToString(
									[]byte("username:password"),
								),
							),
							fmt.Sprintf(
								"%s %s",
								JWT,
								"abcdefghijklmnopqrstuvwxyzABCDEFGHUJKLMNOPQRSTUVWXYZ0123456789",
							),
						},
					},
				},
			},
			depends: depends{
				headers: http.Header{
					http.CanonicalHeaderKey(Authorization): {
						fmt.Sprintf(
							"%s %s",
							Basic,
							base64.StdEncoding.EncodeToString(
								[]byte("username:password"),
							),
						),
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := &request{
				header: tc.depends.headers,
			}

			assert.Equal(
				t,
				tc.want.req,
				req.WithJWTAuth(
					tc.args.value,
				),
			)
		})
	}

}

func TestRequest_WithQuery(t *testing.T) {
	type args struct {
		key    string
		values []any
	}

	type want struct {
		req *request
	}

	type depends struct {
		queries url.Values
	}

	type test struct {
		name    string
		args    args
		want    want
		depends depends
	}

	tests := []test{
		{
			name: "Without collision",
			args: args{
				key:    "key",
				values: []any{"value 1", "value 2", "value 3", 4, true, false},
			},
			want: want{
				req: &request{
					query: url.Values{
						"key": {"value 1", "value 2", "value 3", "4", "true", "false"},
					},
				},
			},
			depends: depends{
				queries: make(url.Values),
			},
		},
		{
			name: "With collision",
			args: args{
				key:    "key",
				values: []any{"value 4", "value 5", "value 6", 7, true, false},
			},
			want: want{
				req: &request{
					query: url.Values{
						"key": {"value 1", "value 2", "value 3", "value 4", "value 5", "value 6", "7", "true", "false"},
					},
				},
			},
			depends: depends{
				queries: url.Values{
					"key": {"value 1", "value 2", "value 3"},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := &request{
				query: tc.depends.queries,
			}

			assert.Equal(
				t,
				tc.want.req,
				req.WithQuery(
					tc.args.key,
					tc.args.values...,
				),
			)

		})
	}

}

func TestRequest_WithQueries(t *testing.T) {
	type args struct {
		queries map[string][]string
	}

	type want struct {
		req *request
	}

	type depends struct {
		queries url.Values
	}

	type test struct {
		name    string
		args    args
		want    want
		depends depends
	}

	tests := []test{
		{
			name: "Without collision",
			args: args{
				queries: map[string][]string{
					"key 2": {"value 4", "value 5", "value 6"},
				},
			},
			want: want{
				req: &request{
					query: url.Values{
						"key":   {"value 1", "value 2", "value 3"},
						"key 2": {"value 4", "value 5", "value 6"},
					},
				},
			},
			depends: depends{
				queries: url.Values{
					"key": {"value 1", "value 2", "value 3"},
				},
			},
		},
		{
			name: "With collision",
			args: args{
				queries: http.Header{
					"key": {"value 4", "value 5", "value 6"},
				},
			},
			want: want{
				req: &request{
					query: url.Values{
						"key": {"value 4", "value 5", "value 6"},
					},
				},
			},
			depends: depends{
				queries: url.Values{
					"key": {"value 1", "value 2", "value 3"},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := &request{
				query: tc.depends.queries,
			}

			assert.Equal(
				t,
				tc.want.req,
				req.WithQueries(tc.args.queries),
			)

		})
	}

}

func TestRequest_WithTimeout(t *testing.T) {
	type args struct {
		timeout time.Duration
	}

	type want struct {
		req *request
	}

	type test struct {
		name string
		args args
		want want
	}

	tests := []test{
		{
			name: "With timeout",
			args: args{
				timeout: time.Second,
			},
			want: want{
				req: &request{
					timeout: time.Second,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := &request{}

			assert.Equal(
				t,
				tc.want.req,
				req.WithTimeout(tc.args.timeout),
			)

		})
	}
}
