package request

import "net/http"

type Interceptor func(http.RoundTripper) http.RoundTripper
