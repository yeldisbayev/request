package req

import (
	"net/http"
)

type Response struct {
	*http.Response
}

// Success checks response status code for success.
func (r *Response) Success() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}
