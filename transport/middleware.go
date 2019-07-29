package transport

import (
	"net/http"
)

// Middleware middleware for http.RoundTripper
type Middleware func(r *http.Request, next RoundTrip) (*http.Response, error)

// RoundTrip base Round Trip func
type RoundTrip func(r *http.Request) (*http.Response, error)

// NewMiddleware create middleware for the transport
func NewMiddleware(init http.RoundTripper, mw Middleware) http.RoundTripper {
	return &middleware{
		init: init,
		mw:   mw,
	}
}

// Middleware middleware by init transport
type middleware struct {
	init http.RoundTripper
	mw   Middleware
}

func (tm *middleware) RoundTrip(r *http.Request) (*http.Response, error) {
	return tm.mw(r, tm.init.RoundTrip)
}

// Chain transport middleware
func Chain(handleFunc ...Middleware) Middleware {
	n := len(handleFunc)
	if n > 1 {
		lastI := n - 1
		return func(r *http.Request, next RoundTrip) (*http.Response, error) {
			var (
				chainHandler RoundTrip
				curI         int
			)
			chainHandler = func(currentRequest *http.Request) (*http.Response, error) {
				if curI == lastI {
					return next(currentRequest)
				}
				curI++
				res, err := handleFunc[curI](currentRequest, chainHandler)
				curI--
				return res, err

			}
			return handleFunc[0](r, chainHandler)
		}
	}

	if n == 1 {
		return handleFunc[0]
	}

	return func(r *http.Request, next RoundTrip) (*http.Response, error) {
		return next(r)
	}
}
