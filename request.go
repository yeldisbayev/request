package request

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"time"
)

var (
	ErrNoBody = errors.New("no body")
)

const (
	ContentType               = "Content-Type"
	ApplicationJSON           = "application/json"
	ApplicationXML            = "application/xml"
	ApplicationFormUrlencoded = "application/x-www-form-urlencoded"
	MultipartFormData         = "multipart/form-data"

	Authorization = "Authorization"
	Basic         = "Basic"
	Bearer        = "Bearer"
	JWT           = "JWT"
)

type Request interface {
	Get(
		ctx context.Context,
		url string,
	) (resp *Response, err error)

	Head(
		ctx context.Context,
		url string,
	) (resp *Response, err error)

	Post(
		ctx context.Context,
		url string,
		body io.Reader,
	) (resp *Response, err error)

	Put(
		ctx context.Context,
		url string,
		body io.Reader,
	) (resp *Response, err error)

	Delete(
		ctx context.Context,
		url string,
	) (resp *Response, err error)

	Connect(
		ctx context.Context,
		url string,
	) (resp *Response, err error)

	Options(
		ctx context.Context,
		url string,
	) (resp *Response, err error)

	Trace(
		ctx context.Context,
		url string,
	) (resp *Response, err error)

	Patch(
		ctx context.Context,
		url string,
	) (resp *Response, err error)

	URL() *url.URL

	Header() http.Header

	Body() (io.Reader, error)

	WithHeader(
		key string,
		values ...string,
	) Request

	WithHeaders(
		values map[string][]string,
	) Request

	WithContentType(
		value string,
	) Request

	WithJSONContentType() Request

	WithXMLContentType() Request

	WithMultipartFormContentType() Request

	WithFormContentType() Request

	WithAuth(
		values ...string,
	) Request

	WithBasicAuth(
		username,
		password string,
	) Request

	WithBearerAuth(
		value string,
	) Request

	WithJWTAuth(
		value string,
	) Request

	WithQuery(
		name string,
		values ...any,
	) Request

	WithQueries(
		values map[string][]string,
	) Request

	WithTimeout(
		timeout time.Duration,
	) Request
}

type request struct {
	httpReq *http.Request
	client  *client
	header  http.Header
	query   url.Values
	timeout time.Duration
}

func (r *request) do(
	ctx context.Context,
	method string,
	url string,
	body io.Reader,
) (resp *Response, err error) {
	timeout := r.timeout
	if timeout == 0 {
		timeout = r.client.timeout
	}

	ctxWithTimeout, cancel := context.WithTimeout(
		ctx,
		timeout,
	)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctxWithTimeout,
		method,
		url,
		body,
	)
	if err != nil {
		return nil, err
	}

	req.Header = r.header
	req.URL.RawQuery = r.query.Encode()

	r.httpReq = req

	res, err := r.client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return &Response{
		Response: res,
	}, err

}

// Get method does GET HTTP request.
func (r *request) Get(
	ctx context.Context,
	url string,
) (resp *Response, err error) {
	return r.do(
		ctx,
		http.MethodGet,
		url,
		nil,
	)

}

// Head method does HEAD HTTP request.
func (r *request) Head(
	ctx context.Context,
	url string,
) (resp *Response, err error) {
	return r.do(
		ctx,
		http.MethodHead,
		url,
		nil,
	)

}

// Post method does POST HTTP request.
func (r *request) Post(
	ctx context.Context,
	url string,
	body io.Reader,
) (resp *Response, err error) {
	return r.do(
		ctx,
		http.MethodPost,
		url,
		body,
	)

}

// Put method does PUT HTTP request.
func (r *request) Put(
	ctx context.Context,
	url string,
	body io.Reader,
) (resp *Response, err error) {
	return r.do(
		ctx,
		http.MethodPut,
		url,
		body,
	)

}

// Delete method does DELETE HTTP request.
func (r *request) Delete(
	ctx context.Context,
	url string,
) (resp *Response, err error) {
	return r.do(
		ctx,
		http.MethodDelete,
		url,
		nil,
	)

}

// Connect method does CONNECT HTTP request.
func (r *request) Connect(
	ctx context.Context,
	url string,
) (resp *Response, err error) {
	return r.do(
		ctx,
		http.MethodConnect,
		url,
		nil,
	)

}

// Options method does OPTIONS HTTP request.
func (r *request) Options(
	ctx context.Context,
	url string,
) (resp *Response, err error) {
	return r.do(
		ctx,
		http.MethodOptions,
		url,
		nil,
	)

}

// Trace method does TRACE HTTP request.
func (r *request) Trace(
	ctx context.Context,
	url string,
) (resp *Response, err error) {
	return r.do(
		ctx,
		http.MethodTrace,
		url,
		nil,
	)

}

// Patch method does PATCH HTTP request.
func (r *request) Patch(
	ctx context.Context,
	url string,
) (resp *Response, err error) {
	return r.do(
		ctx,
		http.MethodPatch,
		url,
		nil,
	)

}

// URL returns request URL.
func (r *request) URL() *url.URL {
	if r.httpReq != nil {
		return r.httpReq.URL
	}

	return nil
}

// Header returns request HEADER.
func (r *request) Header() http.Header {
	if r.httpReq != nil {
		return r.httpReq.Header
	}

	return r.header

}

// Body returns request BODY copy.
func (r *request) Body() (io.Reader, error) {
	if r.httpReq != nil {
		return r.httpReq.GetBody()
	}

	return nil, ErrNoBody
}

// WithHeader adds given HEADER values by key.
func (r *request) WithHeader(
	key string,
	values ...string,
) Request {
	for _, value := range values {
		r.header.Add(key, value)
	}

	return r

}

// WithHeaders sets HEADER.
func (r *request) WithHeaders(
	values map[string][]string,
) Request {
	maps.Copy(r.header, values)

	return r

}

// WithContentType sets content type HEADER.
func (r *request) WithContentType(
	value string,
) Request {
	r.header.Set(
		ContentType,
		value,
	)

	return r

}

// WithJSONContentType sets application/json content type HEADER.
func (r *request) WithJSONContentType() Request {
	r.header.Set(
		ContentType,
		ApplicationJSON,
	)

	return r

}

// WithXMLContentType sets application/xml content type HEADER.
func (r *request) WithXMLContentType() Request {
	r.header.Set(
		ContentType,
		ApplicationXML,
	)

	return r

}

// WithFormContentType sets application/x-www-form-urlencoded content type HEADER.
func (r *request) WithFormContentType() Request {
	r.header.Set(
		ContentType,
		ApplicationFormUrlencoded,
	)

	return r

}

// WithMultipartFormContentType sets multipart/form-data content type HEADER.
func (r *request) WithMultipartFormContentType() Request {
	r.header.Set(
		ContentType,
		MultipartFormData,
	)

	return r

}

// WithAuth adds given values to authorization HEADER.
func (r *request) WithAuth(
	values ...string,
) Request {
	for _, value := range values {
		r.header.Add(Authorization, value)
	}

	return r

}

// WithBasicAuth adds Basic authorization HEADER.
func (r *request) WithBasicAuth(
	username,
	password string,
) Request {
	auth := fmt.Sprintf(
		"%s %s",
		Basic,
		base64.StdEncoding.EncodeToString(
			[]byte(username+":"+password),
		),
	)

	r.header.Add(Authorization, auth)

	return r

}

// WithBearerAuth adds Bearers authorization HEADER.
func (r *request) WithBearerAuth(
	value string,
) Request {
	r.header.Add(
		Authorization,
		fmt.Sprintf(
			"%s %s",
			Bearer,
			value,
		),
	)

	return r

}

// WithJWTAuth adds JWT authorization HEADER.
func (r *request) WithJWTAuth(
	value string,
) Request {
	r.header.Add(
		Authorization,
		fmt.Sprintf(
			"%s %s",
			JWT,
			value,
		),
	)

	return r

}

// WithQuery adds given query parameter values by name.
func (r *request) WithQuery(
	name string,
	values ...any,
) Request {
	for _, value := range values {
		r.query.Add(
			name,
			fmt.Sprintf(
				"%v",
				value,
			),
		)
	}

	return r

}

// WithQueries adds given query parameters.
func (r *request) WithQueries(
	values map[string][]string,
) Request {
	maps.Copy(r.query, values)

	return r

}

// WithTimeout sets request timeout and implemented with context.Context.
// Request timeout has higher priority than Client's timeout
func (r *request) WithTimeout(
	timeout time.Duration,
) Request {
	r.timeout = timeout

	return r

}
