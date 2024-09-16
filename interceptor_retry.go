package req

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
	http.StatusTooManyRequests,
	http.StatusServiceUnavailable,
	http.StatusGatewayTimeout,
}

func Retry(retryOnStatusCodes ...int) Interceptor {
	statusCodes := defaultStatusCodes

	if len(retryOnStatusCodes) != 0 {
		statusCodes = retryOnStatusCodes
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

				defer func() {
					retries := 0
					for shouldRetry(res, err, statusCodes) && retries < maxRetries {
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
				}()

				return tripper.RoundTrip(req)

			},
		)
	}
}

func delay(retries int) time.Duration {
	return time.Duration(math.Pow(2, float64(retries))) * time.Second
}

func shouldRetry(res *http.Response, err error, statusCodes []int) bool {
	if err != nil {
		return true
	}

	if res != nil && slices.Contains(statusCodes, res.StatusCode) {
		return true
	}

	return false

}

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

func drainBody(res *http.Response) {
	if res != nil && res.Body != nil {
		_, _ = io.Copy(io.Discard, res.Body)
		_ = res.Body.Close()
	}

}
