package decoder

import (
	"errors"
	"io"
	"mime"
	"net/http"
	"sync"
)

var (
	decodersMu sync.RWMutex
	decoders   = make(map[string]Decoder)
)

// Decoder decode by reader
type Decoder func(r io.Reader, v interface{}) error

// HTTPDecode decode by MediaType
func HTTPDecode(r *http.Response, body io.Reader, v interface{}) error {
	mt, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return err
	}
	return Decode(mt, body, v)
}

// Decode by media type
func Decode(mediaType string, body io.Reader, v interface{}) error {
	decodersMu.RLock()
	d, ok := decoders[mediaType]
	decodersMu.RUnlock()
	if ok {
		return d(body, v)
	}
	return errors.New("http client: decoder by media type '" + mediaType + "' not found")
}

// Register decoder by media type
func Register(decoder Decoder, mediaTypes ...string) error {
	if decoder == nil || len(mediaTypes) == 0 {
		return errors.New("http client: decider and media types is required")
	}
	decodersMu.Lock()
	defer decodersMu.Unlock()
	for _, mt := range mediaTypes {
		if _, dup := decoders[mt]; dup {
			return errors.New("http client: register called twice for decoder by media type " + mt)
		}
		decoders[mt] = decoder
	}

	return nil
}

// MustRegister register decode or panic if duplicate
func MustRegister(decoder Decoder, mediaTypes ...string) {
	if err := Register(decoder, mediaTypes...); err != nil {
		panic(err)
	}
}
