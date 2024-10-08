package req

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"time"
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

	Body() io.Reader

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

func (r *request) URL() *url.URL {
	if r.httpReq != nil {
		return r.httpReq.URL
	}

	return nil
}

func (r *request) Header() http.Header {
	if r.httpReq != nil {
		return r.httpReq.Header
	}

	return r.header

}

func (r *request) Body() io.Reader {
	if r.httpReq != nil {
		return r.httpReq.Body
	}

	return nil

}

func (r *request) WithHeader(
	key string,
	values ...string,
) Request {
	for _, value := range values {
		r.header.Add(key, value)
	}

	return r

}

func (r *request) WithHeaders(
	values map[string][]string,
) Request {
	maps.Copy(r.header, values)

	return r

}

func (r *request) WithContentType(
	value string,
) Request {
	r.header.Set(
		ContentType,
		value,
	)

	return r

}

func (r *request) WithJSONContentType() Request {
	r.header.Set(
		ContentType,
		ApplicationJSON,
	)

	return r

}

func (r *request) WithXMLContentType() Request {
	r.header.Set(
		ContentType,
		ApplicationXML,
	)

	return r

}

func (r *request) WithFormContentType() Request {
	r.header.Set(
		ContentType,
		ApplicationFormUrlencoded,
	)

	return r

}

func (r *request) WithMultipartFormContentType() Request {
	r.header.Set(
		ContentType,
		MultipartFormData,
	)

	return r

}

func (r *request) WithAuth(
	values ...string,
) Request {
	for _, value := range values {
		r.header.Add(Authorization, value)
	}

	return r

}

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

func (r *request) WithQueries(
	values map[string][]string,
) Request {
	maps.Copy(r.query, values)

	return r

}

func (r *request) WithTimeout(
	timeout time.Duration,
) Request {
	r.timeout = timeout

	return r

}
