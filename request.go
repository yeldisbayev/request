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
	ContentType                  = "Content-Type"
	ContentTypeApplicationJSON   = "application/json"
	ContentTypeApplicationXML    = "application/xml"
	ContentTypeMultipartFormData = "multipart/form-data"
	ContentTypeFormUrlencoded    = "application/x-www-form-urlencoded"

	Authorization = "Authorization"
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

	WithHeader(
		key string,
		values ...string,
	) Request

	WithContentType(
		value string,
	) Request

	WithJSONContentType() Request

	WithXMLContentType() Request

	WithMultipartFormDataContentType() Request

	WithFormURLEncodedContentType() Request

	WithAuthorization(
		values ...string,
	) Request

	WithBasicAuthorization(
		username,
		password string,
	) Request

	WithBearerAuthorization(
		value string,
	) Request

	WithJWTAuthorization(
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
	req     *http.Request
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

	req, _ := http.NewRequestWithContext(
		ctxWithTimeout,
		method,
		url,
		body,
	)

	req.Header = r.header
	req.URL.RawQuery = r.query.Encode()

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

func (r *request) WithHeader(
	key string,
	values ...string,
) Request {
	for _, value := range values {
		r.req.Header.Add(key, value)
	}

	return r

}

func (r *request) WithAuthorization(
	values ...string,
) Request {
	for _, value := range values {
		r.req.Header.Add("Authorization", value)
	}

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
		ContentTypeApplicationJSON,
	)

	return r

}

func (r *request) WithXMLContentType() Request {
	r.header.Set(
		ContentType,
		ContentTypeApplicationXML,
	)

	return r

}

func (r *request) WithMultipartFormDataContentType() Request {
	r.header.Set(
		ContentType,
		ContentTypeMultipartFormData,
	)

	return r

}

func (r *request) WithFormURLEncodedContentType() Request {
	r.header.Set(
		ContentType,
		ContentTypeFormUrlencoded,
	)

	return r

}

func (r *request) WithBasicAuthorization(
	username,
	password string,
) Request {
	auth := fmt.Sprintf(
		"Basic %s",
		base64.StdEncoding.EncodeToString(
			[]byte(username+":"+password),
		),
	)

	r.header.Add(Authorization, auth)

	return r

}

func (r *request) WithBearerAuthorization(
	value string,
) Request {
	r.header.Add(
		Authorization,
		fmt.Sprintf(
			"Bearer %s",
			value,
		),
	)

	return r

}

func (r *request) WithJWTAuthorization(
	value string,
) Request {
	r.header.Add(
		Authorization,
		fmt.Sprintf(
			"JWT %s",
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
	query := url.Values(values)
	maps.Copy(r.query, query)

	return r

}

func (r *request) WithTimeout(
	timeout time.Duration,
) Request {
	r.timeout = timeout

	return r

}
