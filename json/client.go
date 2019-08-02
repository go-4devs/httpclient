package json

import (
	"encoding/json"
	"io"

	"github.com/go-4devs/httpclient/dc"
	"github.com/go-4devs/httpclient/decoder"
)

// DefaultDecoder default json decoder
var DefaultDecoder dc.Decoder = func(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

// RegisterDecoder by application/json with aliases content type
func RegisterDecoder(aliases ...string) {
	decoder.MustRegister(func(r io.Reader, v interface{}) error {
		return json.NewDecoder(r).Decode(v)
	}, append(aliases, "application/json")...)
}

// NewClient create client with json decoder
func NewClient(baseURL string, opts ...dc.Option) (dc.Client, error) {
	opts = append(opts, ds.WithDecoder(DefaultDecoder))
	return dc.New(baseURL, opts...)
}
