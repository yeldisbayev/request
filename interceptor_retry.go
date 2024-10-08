package request

import (
	"context"
	"io"
	"math"
	"net/http"
	"slices"
	"time"
)

const maxRetries = 3

var defaultStatusCodes = []int{
	http.StatusRequestTimeout,
	http.StatusTooEarly,
	http.StatusTooManyRequests,
	http.StatusBadGateway,
	http.StatusServiceUnavailable,
	http.StatusGatewayTimeout,
}

// Retry interceptor retry request on request failure
// or on defined Response status codes.
// By default Retry uses defaultStatusCodes
func Retry(statusCodes ...int) Interceptor {
	retryStatusCodes := defaultStatusCodes

	if len(statusCodes) != 0 {
		retryStatusCodes = statusCodes
	}

	return func(tripper http.RoundTripper) http.RoundTripper {
		return RoundTripper(
			func(req *http.Request) (res *http.Response, err error) {
				var body io.ReadCloser
				if req.Body != nil {
					body, err = req.GetBody()
					if err != nil {
						return res, err
					}
				}

				res, err = tripper.RoundTrip(req)
				retries := 0
				for shouldRetry(res, err, retryStatusCodes) && retries < maxRetries {
					if retries != 0 {
						sleepWithContext(
							req.Context(),
							delay(retries),
						)
					}

					drainBody(res)

					if req.Body != nil {
						req.Body = body
					}

					res, err = tripper.RoundTrip(req)
					retries++

				}

				return res, err

			},
		)
	}
}

// delay calculates Retry duration
func delay(retries int) time.Duration {
	return time.Duration(math.Pow(2, float64(retries))) * time.Second
}

// shouldRetry determine conditions to Retry interceptor
// by including Response status code
func shouldRetry(res *http.Response, err error, statusCodes []int) bool {
	if err != nil {
		return true
	}

	if res != nil && slices.Contains(statusCodes, res.StatusCode) {
		return true
	}

	return false

}

// sleepWithContext delays Retry interceptor
// considering its context and duration
func sleepWithContext(ctx context.Context, d time.Duration) {
	timer := time.NewTimer(d)

	select {
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
	case <-timer.C:
	}

}

// drainBody drain http.Response body
// to reuse same connection
func drainBody(res *http.Response) {
	if res != nil && res.Body != nil {
		_, _ = io.Copy(io.Discard, res.Body)
		_ = res.Body.Close()
	}

}
