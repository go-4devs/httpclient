package httpclient

import (
	"net/http"
)

// TransportMiddleware middleware for http.RoundTripper
type TransportMiddleware func(r *http.Request, next RoundTrip) (*http.Response, error)

// RoundTrip base Round Trip func
type RoundTrip func(r *http.Request) (*http.Response, error)

// NewTransportMiddleware create middleware for the transport
func NewTransportMiddleware(
	transport http.RoundTripper,
	middleware TransportMiddleware,
) http.RoundTripper {
	return &transportMiddleware{
		init: transport,
		mw:   middleware,
	}
}

// TransportMiddleware middleware by init transport
type transportMiddleware struct {
	init http.RoundTripper
	mw   TransportMiddleware
}

func (tm *transportMiddleware) RoundTrip(r *http.Request) (*http.Response, error) {
	return tm.mw(r, tm.init.RoundTrip)
}

func chain(handleFunc ...TransportMiddleware) TransportMiddleware {
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
