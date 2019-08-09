package decoder

import (
	"errors"
	"io"
	"net/http"
	"sync"
)

var (
	decodersMu sync.RWMutex
	decoders   = make(map[string]Decoder)
)

// decoder decode or error
type Decoder func(r io.Reader, v interface{}) error

// HTTPDecode decode by Content-Type
func HTTPDecode(r *http.Response, body io.Reader, v interface{}) error {
	ct := r.Header.Get("Content-Type")
	decodersMu.RLock()
	d, ok := decoders[ct]
	decodersMu.RUnlock()
	if ok {
		return d(body, v)
	}
	return errors.New("http client: decoder by content type'" + ct + "' not found")
}

// MustRegister register decode or panic if duplicate
func MustRegister(decoder Decoder, contentTypes ...string) {
	decodersMu.Lock()
	defer decodersMu.Unlock()
	if decoder == nil {
		panic("http client: decider is nil")
	}
	for _, ct := range contentTypes {
		if _, dup := decoders[ct]; dup {
			panic("http client:  Register called twice for decoder by content type" + ct)
		}
		decoders[ct] = decoder
	}
}
