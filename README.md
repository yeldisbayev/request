## What is Req

Req is golang enhanced http client

## Installation

```
$ go get -u github.com/yeldisbayev/req
```

## Example usage

```go
package main

import (
	"encoding/json"
	"log"
)

type Todo struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"userID"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func main() {
	
	client := req.NewClient()

	res, err := client.Request().Get(context.Background(), "https://jsonplaceholder.typicode.com/todos/1")
	if err != nil {
		log.Fatalf("Request error: %v", err)
	}

	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(res.Body)

	var todo Todo

	if res.Success() {
		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&todo); err != nil {
			log.Fatalf("Decode error: %v", err)
		}
	}

	log.Printf(
		"todo = (id: %d, userId: %d, title: %s, completed: %t)",
		todo.ID,
		todo.UserID,
		todo.Title,
		todo.Completed,
	)

}

```

## Client

### Options

There is Options to set up Client

#### WithTimeout

Sets timeout for all client requests. Timeout implemented without using http.Client's Timeout property, but with context. Client timeout has lesser priority than Request timeout property. If not provided, [DefaultTimeout](https://github.com/yeldisbayev/req/blob/89f395aca69a4a1ddb28fe08951ceff807238c2f/client.go#L10C2-L10C16) is used.

```go
timeout := time.Second * 5

client := req.NewClient(req.WithTimeout(timeout))
```

#### WithIdleConnectionTimeout

Controls amount of time an idle (keep-alive) connection will remain idle before closing itself. If not provided, [DefaultIdleConnectionTimeout](https://github.com/yeldisbayev/req/blob/89f395aca69a4a1ddb28fe08951ceff807238c2f/client.go#L11C2-L11C30) is used.

```go
idleConnectionTimeout := time.Second * 30

client := req.NewClient(req.WithIdleConnectionTimeout(idleConnectionTimeout))
```

#### WithMaxIdleConnections

Controls the maximum number of idle (keep-alive) connections across all hosts. If not provided, [DefaultMaxIdleConnections](https://github.com/yeldisbayev/req/blob/89f395aca69a4a1ddb28fe08951ceff807238c2f/client.go#L12) is used.

```go
maxIdleConnections := 20

client := req.NewClient(req.WithMaxIdleConnections(maxIdleConnections))
```

#### WithMaxConnectionsPerHost

Optionally limits the total number of connections per host, including connections in the dialing, active, and idle states. On limit violation, dials will block. If not provided, [DefaultMaxConnectionsPerHost](https://github.com/yeldisbayev/req/blob/89f395aca69a4a1ddb28fe08951ceff807238c2f/client.go#L13) is used.

```go
maxConnectionsPerHost := 30

client := req.NewClient(req.WithMaxConnectionsPerHost(maxConnectionsPerHost))
```

#### WithMaxIdleConnectionsPerHost

Controls the maximum idle (keep-alive) connections to keep per-host. If not provided, [DefaultMaxIdleConnectionsPerHost](https://github.com/yeldisbayev/req/blob/89f395aca69a4a1ddb28fe08951ceff807238c2f/client.go#L14) is used.

```go
maxIdleConnectionsPerHost := 30

client := req.NewClient(req.WithMaxIdleConnectionsPerHost(maxIdleConnectionsPerHost))
```

#### WithForceAttemptHTTP2

Controls whether HTTP/2 is enabled when a non-zero Dial, DialTLS, or DialContext func or TLSClientConfig is provided. By default, use of any those fields conservatively disables HTTP/2. To use a custom dialer or TLS config and still attempt HTTP/2 upgrades, set this to true. If not provided, [DefaultForceAttemptHTTP2](https://github.com/yeldisbayev/req/blob/89f395aca69a4a1ddb28fe08951ceff807238c2f/client.go#L15) is used.

```go
forceAttemptHTTP2 := true

client := req.NewClient(req.WithForceAttemptHTTP2(forceAttemptHTTP2))
```

#### WithInterceptors

Wraps Client with given [interceptors](https://github.com/yeldisbayev/req/blob/48f91285a13c6e2ed3afd768bc3692996af9e62b/interceptor.go#L5)

```go
retry := req.Retry()

client := req.NewClient(req.WithInterceptors(retry))
```

## Interceptor

Interceptor wraps http Transport and calls before or after due to client usage. In order to create custom one [Interceptor](https://github.com/yeldisbayev/req/blob/48f91285a13c6e2ed3afd768bc3692996af9e62b/interceptor.go#L5) function implementation is needed. There is also built in [Retry](https://github.com/yeldisbayev/req/blob/4ec32c09e979df025d0ba4967e5ea52e9f2d5cdf/interceptor_retry.go#L26C6-L26C11) interceptor and its should be at the end in interceptors chain.

```go
package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"req"
	"time"
)

func FirstInterceptor(rt http.RoundTripper) http.RoundTripper {
	return req.RoundTripper( // Wrap custom interceptor with req.RoundTripper
		func(req *http.Request) (res *http.Response, err error) {
			// Before request started
			log.Println("FirstInterceptor started")

			res, err = rt.RoundTrip(req)
            
			// After request completed
			log.Println("FirstInterceptor finished")

			return res, err

		},
	)
}

func SecondInterceptor(rt http.RoundTripper) http.RoundTripper {
	return req.RoundTripper( // Wrap custom interceptor with req.RoundTripper
		func(req *http.Request) (res *http.Response, err error) {
			// Before request started
			log.Println("SecondInterceptor started")
			startTime := time.Now()

			res, err = rt.RoundTrip(req)

			// After request completed
			log.Println("SecondInterceptor", "Request Duration", time.Since(startTime))
			log.Println("SecondInterceptor finished")

			return res, err

		},
	)
}

type Todo struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"userID"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func main() {
	client := req.NewClient(
		req.WithInterceptors(FirstInterceptor, SecondInterceptor), 
	)

	res, err := client.Request().Get(context.Background(), "https://jsonplaceholder.typicode.com/todos/1")
	if err != nil {
		log.Fatalf("Request error: %v", err)
	}

	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(res.Body)

	var todo Todo

	if res.Success() {
		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&todo); err != nil {
			log.Fatalf("Decode error: %v", err)
		}
	}

	log.Println(todo)

}
```

```text
 FirstInterceptor started
 
 SecondInterceptor started
 
 SecondInterceptor Request Duration 697.908708ms
 
 SecondInterceptor finished
 
 FirstInterceptor finished
 
 {1 1 delectus aut autem false}
```

## License

MIT License

Copyright (c) 2024 Duisen Yeldisbayev

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.