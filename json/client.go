package json

import (
	"encoding/json"
	"io"

	"github.com/go-4devs/httpclient/dc"
	"github.com/go-4devs/httpclient/decoder"
)

// defaultDecoder default json decoder
var defaultDecoder decoder.Decoder = func(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

// RegisterDecoder by application/json with aliases content type
func RegisterDecoder(aliases ...string) {
	decoder.MustRegister(defaultDecoder, append(aliases, "application/json")...)
}

// NewClient create client with json decoder
func NewClient(baseURL string, opts ...dc.Option) (dc.Client, error) {
	opts = append(opts, dc.WithDecoder(defaultDecoder))
	return dc.New(baseURL, opts...)
}
