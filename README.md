# go-re

[![CI](https://github.com/da440dil/go-re/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/da440dil/go-re/actions/workflows/go.yml)
[![Coverage Status](https://coveralls.io/repos/github/da440dil/go-re/badge.svg?branch=main)](https://coveralls.io/github/da440dil/go-re?branch=main)
[![Go Reference](https://pkg.go.dev/badge/github.com/da440dil/go-re.svg)](https://pkg.go.dev/github.com/da440dil/go-re)
[![Go Report Card](https://goreportcard.com/badge/github.com/da440dil/go-re)](https://goreportcard.com/report/github.com/da440dil/go-re)

Abstraction for retry strategies.

## [Example](./examples/http_test.go) HTTP request

```go
fn := func(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("status code != %d, %w", http.StatusOK, re.ErrRetryable)
	}
	return fmt.Sprintf("statusCode: %d, body: %s", res.StatusCode, body), err
}
// Retry function execution in case of an error after 10ms, 20ms, 30ms.
fn = re.Tryable(fn, slices.Values([]time.Duration{
	time.Millisecond * 10, time.Millisecond * 20, time.Millisecond * 30,
}))

x := 0
h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	statusCode := http.StatusOK
	if x > 2 || x%2 != 0 { // { 0 => 200, 1 => 418, 2 => 200, ... => 418 }
		statusCode = http.StatusTeapot
	}
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, "{ x: %d }", x)
	x++
})
srv := httptest.NewServer(h)
defer srv.Close()

for range 3 {
	v, err := fn(context.Background(), srv.URL)
	fmt.Printf("%v, err: %v\n", v, err)
}

ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
defer cancel()

v, err := fn(ctx, srv.URL)
fmt.Printf("%v, err: %v", v, err)
// Output:
// statusCode: 200, body: { x: 0 }, err: <nil>
// statusCode: 200, body: { x: 2 }, err: <nil>
// statusCode: 418, body: { x: 6 }, err: status code != 200, retryable
// statusCode: 418, body: { x: 9 }, err: context deadline exceeded
```
