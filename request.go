package req

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
)

type Request interface {
	Get(
		ctx context.Context,
		url string,
	) (resp *http.Response, err error)

	Head(
		ctx context.Context,
		url string,
	) (resp *http.Response, err error)

	Post(
		ctx context.Context,
		url string,
		body io.Reader,
	) (resp *http.Response, err error)

	Put(
		ctx context.Context,
		url string,
		body io.Reader,
	) (resp *http.Response, err error)

	Delete(
		ctx context.Context,
		url string,
	) (resp *http.Response, err error)

	Connect(
		ctx context.Context,
		url string,
	) (resp *http.Response, err error)

	Options(
		ctx context.Context,
		url string,
	) (resp *http.Response, err error)

	Trace(
		ctx context.Context,
		url string,
	) (resp *http.Response, err error)

	Patch(
		ctx context.Context,
		url string,
	) (resp *http.Response, err error)

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

	WithFormURLContentType() Request

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
}

type request struct {
	req    *http.Request
	client *client
	header http.Header
}

func (r *request) Get(
	ctx context.Context,
	url string,
) (resp *http.Response, err error) {
	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)

	req.Header = r.header

	return r.client.httpClient.Do(req)

}

func (r *request) Head(
	ctx context.Context,
	url string,
) (resp *http.Response, err error) {
	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodHead,
		url,
		nil,
	)

	req.Header = r.header

	return r.client.httpClient.Do(req)

}

func (r *request) Post(
	ctx context.Context,
	url string,
	body io.Reader,
) (resp *http.Response, err error) {
	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		body,
	)

	req.Header = r.header

	return r.client.httpClient.Do(req)

}

func (r *request) Put(
	ctx context.Context,
	url string,
	body io.Reader,
) (resp *http.Response, err error) {
	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		url,
		body,
	)

	req.Header = r.header

	return r.client.httpClient.Do(req)

}

func (r *request) Delete(
	ctx context.Context,
	url string,
) (resp *http.Response, err error) {
	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		nil,
	)

	req.Header = r.header

	return r.client.httpClient.Do(req)

}

func (r *request) Connect(
	ctx context.Context,
	url string,
) (resp *http.Response, err error) {
	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodConnect,
		url,
		nil,
	)

	req.Header = r.header

	return r.client.httpClient.Do(req)

}

func (r *request) Options(
	ctx context.Context,
	url string,
) (resp *http.Response, err error) {
	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodOptions,
		url,
		nil,
	)

	req.Header = r.header

	return r.client.httpClient.Do(req)

}

func (r *request) Trace(
	ctx context.Context,
	url string,
) (resp *http.Response, err error) {
	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodTrace,
		url,
		nil,
	)

	req.Header = r.header

	return r.client.httpClient.Do(req)

}

func (r *request) Patch(
	ctx context.Context,
	url string,
) (resp *http.Response, err error) {
	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		url,
		nil,
	)

	req.Header = r.header

	return r.client.httpClient.Do(req)

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

func (r *request) WithFormURLContentType() Request {
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
	r.header.Set(
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
	r.header.Set(
		"Authorization",
		fmt.Sprintf(
			"JWT %s",
			value,
		),
	)

	return r

}
