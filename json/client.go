package json

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-4devs/httpclient/dc"
)

// DefaultDecoder default json decoder
var DefaultDecoder dc.Decoder = func(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

// NewClient create client with json decoder
func NewClient(baseURL string, opts ...dc.Option) (dc.Client, error) {
	opts = append(opts,
		dc.WithTransportMiddleware(func(
			r *http.Request,
			next func(r *http.Request) (*http.Response, error),
		) (*http.Response, error) {
			r.Header.Add("Content-Type", "application/json")
			return next(r)
		}),
	)

	return dc.New(baseURL, DefaultDecoder, opts...)
}
