package httpclient

import (
	"io"
	"net/http"
)

// Fetch interface for the get response and processed it
type Fetch interface {
	IsStatus(httpStatus int) bool
	With(func(r *http.Response, b io.Reader) error) Fetch
	Decode(v interface{}) error
	Body() io.Reader
	Error() error
}

// Fetcher fetch response
type Fetcher interface {
	Client
	Fetch(r *http.Request) Fetch
}
