package re_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/da440dil/go-re"
)

func Example_http() {
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
		return fmt.Sprintf("{ statusCode: %d, body: %s }", res.StatusCode, body), err
	}
	// Use constant delay between retries 10 ms with maximum number of retries 1.
	fn = re.Tryable(fn, re.Constant(time.Millisecond*10), re.MaxRetries(1))

	x := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		statusCode := http.StatusOK
		if x > 2 || x%2 != 0 { // { 0 => 200, 1 => 418, 2 => 200, ... => 418 }
			statusCode = http.StatusTeapot
		}
		w.WriteHeader(statusCode)
		fmt.Fprintf(w, "{ x: %d }", x)
		x++
	}))
	defer srv.Close()

	for i := 0; i < 3; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
		defer cancel()

		v, err := fn(ctx, srv.URL)
		fmt.Printf("{ v: %v, err: %v }\n", v, err)
	}
	// Output:
	// { v: { statusCode: 200, body: { x: 0 } }, err: <nil> }
	// { v: { statusCode: 200, body: { x: 2 } }, err: <nil> }
	// { v: { statusCode: 418, body: { x: 4 } }, err: status code != 200, retryable }
}
