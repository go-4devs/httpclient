package json

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-4devs/httpclient/dc"
	"github.com/go-4devs/httpclient/transport"
)

// NewJSONDecoder new JSON decoder
func NewJSONDecoder() dc.Decoder {
	return func(r io.Reader, v interface{}) error {
		return json.NewDecoder(r).Decode(v)
	}
}

// NewJSONClient create client with json decoder
func NewJSONClient(baseURL string, opts ...dc.Option) (dc.Client, error) {
	opts = append(opts,
		dc.WithTransportMiddleware(func(r *http.Request, next transport.RoundTrip) (*http.Response, error) {
			r.Header.Add("Content-Type", "application/json")
			return next(r)
		}),
	)

	return dc.New(baseURL, NewJSONDecoder(), opts...)
}
