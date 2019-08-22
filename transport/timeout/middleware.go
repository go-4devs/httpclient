package timeout

import (
	"context"
	"net/http"
	"time"

	"github.com/go-4devs/httpclient/transport"
)

// New create new timeout middleware
func New(timeout time.Duration) transport.Middleware {
	return func(r *http.Request, n func(r *http.Request) (*http.Response, error)) (*http.Response, error) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		return n(r.WithContext(ctx))
	}
}
