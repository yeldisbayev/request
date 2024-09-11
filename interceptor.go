package req

import "net/http"

type Interceptor func(http.RoundTripper) http.RoundTripper
