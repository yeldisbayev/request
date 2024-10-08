package request

import "net/http"

type RoundTripper func(*http.Request) (*http.Response, error)

func (rt RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt(req)
}
