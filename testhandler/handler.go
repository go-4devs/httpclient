package testhandler

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
)

//Handle test handler.
type Handle struct {
	Request Request
	code    int
	data    string
	cases   []func(*testing.T, *http.Request)
}

// CanHandle check request.
func (h Handle) CanHandle(r *http.Request) bool {
	return h.Request != nil && h.Request.CanHandle(r)
}

// Write response.
func (h Handle) Write(w http.ResponseWriter) Handle {
	if h.code != 0 {
		w.WriteHeader(h.code)
	}
	if h.data != "" {
		if _, err := w.Write([]byte(h.data)); err != nil {
			panic(err)
		}
	}
	return h
}

// Cases run cases by request.
func (h Handle) Cases(t *testing.T, r *http.Request) {
	for _, rt := range h.cases {
		rt(t, r)
	}
}

// Option for handle.
type Option func(*Handle)

// WithCode set custom code.
func WithCode(code int) Option {
	return func(handle *Handle) {
		handle.code = code
	}
}

// WithCodeNotFound set 404 code.
func WithCodeNotFound() Option {
	return func(handle *Handle) {
		handle.code = http.StatusNotFound
	}
}

// WithCodeBadRequest set 400 code.
func WithCodeBadRequest() Option {
	return func(handle *Handle) {
		handle.code = http.StatusBadRequest
	}
}

// WithCodeUnauthorized set 401 code.
func WithCodeUnauthorized() Option {
	return func(handle *Handle) {
		handle.code = http.StatusUnauthorized
	}
}

// WithTestRequest set cases Request.
func WithTestRequest(t ...func(*testing.T, *http.Request)) Option {
	return func(handle *Handle) {
		handle.cases = append(handle.cases, t...)
	}
}

// Request set Request to handler.
type Request interface {
	CanHandle(r *http.Request) bool
}

// NewHandle create handler.
func NewHandle(req Request, data string, opts ...Option) Handle {
	jh := &Handle{
		Request: req,
		data:    data,
		code:    http.StatusOK,
	}
	for _, o := range opts {
		o(jh)
	}
	return *jh
}

// NewHTTPHandler http handler.
func NewHTTPHandler(t *testing.T, h ...Handle) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.GetBody == nil {
			snapshot, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Error(err)
			}
			r.GetBody = func() (io.ReadCloser, error) {
				r := bytes.NewBuffer(snapshot)
				return ioutil.NopCloser(r), nil
			}
		}
		for i := range h {
			if h[i].CanHandle(r) {
				h[i].Write(w).Cases(t, r)
				return
			}
		}
		w.WriteHeader(http.StatusNotImplemented)
	})
}
