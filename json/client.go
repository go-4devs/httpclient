package json

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-4devs/httpclient"
)

// NewJSONDecoder new JSON decoder
func NewJSONDecoder() httpclient.Decoder {
	return func(r io.Reader, v interface{}) error {
		return json.NewDecoder(r).Decode(v)
	}
}

// NewJSONClient create client with json decoder
func NewJSONClient(baseURL string, opts ...httpclient.Option) (httpclient.BaseClient, error) {
	opts = append(opts, httpclient.WithTransportMiddleware(func(r *http.Request, next httpclient.RoundTrip) (*http.Response, error) {
		r.Header.Add("Content-Type", "application/json")
		return next(r)
	}))

	return httpclient.New(baseURL, NewJSONDecoder(), opts...)
}
