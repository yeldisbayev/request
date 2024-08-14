package req

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
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
}

type request struct {
	req    *http.Request
	client *client
	header http.Header
	query  url.Values
}

func (r *request) do(
	ctx context.Context,
	url string,
	body io.Reader,
) (resp *Response, err error) {
	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
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
		"Content-Type",
		value,
	)

	return r

}

func (r *request) WithJSONContentType() Request {
	r.header.Set(
		"Content-Type",
		"application/json",
	)

	return r

}

func (r *request) WithXMLContentType() Request {
	r.header.Set(
		"Content-Type",
		"application/xml",
	)

	return r

}

func (r *request) WithMultipartFormDataContentType() Request {
	r.header.Set(
		"Content-Type",
		"multipart/form-data",
	)

	return r

}

func (r *request) WithFormURLEncodedContentType() Request {
	r.header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
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

	r.header.Add("Authorization", auth)

	return r

}

func (r *request) WithBearerAuthorization(
	value string,
) Request {
	r.header.Add(
		"Authorization",
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
		"Authorization",
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
