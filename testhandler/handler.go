package testhandler

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

//Handle ...
type Handle struct {
	Request Request
	code    int
	data    io.WriterTo
	cases   []func(*testing.T, *http.Request)
}

// CanHandle check request
func (h Handle) CanHandle(r *http.Request) bool {
	return h.Request != nil && h.Request.CanHandle(r)
}

// Write response
func (h Handle) Write(w http.ResponseWriter) Handle {
	if h.code != 0 {
		w.WriteHeader(h.code)
	}
	if h.data != nil {
		if _, err := h.data.WriteTo(w); err != nil {
			panic(err)
		}
	}
	return h
}

// Cases run cases by request
func (h Handle) Cases(t *testing.T, r *http.Request) {
	for _, rt := range h.cases {
		rt(t, r)
	}
}

// Option for handle
type Option func(*Handle)

// WithCode set custom code
func WithCode(code int) Option {
	return func(handle *Handle) {
		handle.code = code
	}
}

// WithCodeNotFound set 404 code
func WithCodeNotFound() Option {
	return func(handle *Handle) {
		handle.code = http.StatusNotFound
	}
}

// WithCodeBadRequest set 400 code
func WithCodeBadRequest() Option {
	return func(handle *Handle) {
		handle.code = http.StatusBadRequest
	}
}

// WithCodeUnauthorized set 401 code
func WithCodeUnauthorized() Option {
	return func(handle *Handle) {
		handle.code = http.StatusUnauthorized
	}
}

// WithTestRequest set cases Request
func WithTestRequest(t ...func(*testing.T, *http.Request)) Option {
	return func(handle *Handle) {
		handle.cases = append(handle.cases, t...)
	}
}

// Request set Request to handler
type Request interface {
	CanHandle(r *http.Request) bool
}

// NewHandle create handler
func NewHandle(req Request, data string, opts ...Option) Handle {
	jh := &Handle{
		Request: req,
		data:    bytes.NewBufferString(data),
		code:    http.StatusOK,
	}
	for _, o := range opts {
		o(jh)
	}
	return *jh
}

// NewHTTPHandler http handler
func NewHTTPHandler(t *testing.T, h ...Handle) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for i := range h {
			if h[i].CanHandle(r) {
				h[i].Write(w).Cases(t, r)
				return
			}
		}
		w.WriteHeader(http.StatusNotImplemented)
	})
}
