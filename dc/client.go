package dc

import (
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/go-4devs/httpclient/transport"
)

// Decoder for decode body
type Decoder func(r io.Reader, v interface{}) error

// HTTPClient get response and marshaling it by decoder
type Client struct {
	HTTPClient http.Client
	Decoder    Decoder
	baseURL    url.URL
}

type options struct {
	transport  http.RoundTripper
	middleware transport.Middleware
}

// Option for the configure HTTPClient
type Option func(*options)

// WithTransportMiddleware add middleware for transport
func WithTransportMiddleware(mw ...transport.Middleware) Option {
	return func(i *options) {
		if i.middleware != nil {
			mw = append([]transport.Middleware{i.middleware}, mw...)
		}
		if len(mw) > 0 {
			i.middleware = transport.Chain(mw...)
		}
	}
}

// New create new HTTPClient
func New(baseURL string, decoder Decoder, opts ...Option) (client Client, err error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return client, err
	}

	cl := Client{
		baseURL: *u,
		Decoder: decoder,
	}
	op := &options{}
	for _, opt := range opts {
		opt(op)
	}
	tr := http.DefaultTransport
	if op.transport == nil {
		tr = op.transport
	}

	if op.middleware != nil {
		tr = transport.NewMiddleware(tr, op.middleware)
	}
	cl.HTTPClient = http.Client{
		Transport: tr,
	}
	return cl, nil

}

// Do request and decode response body
func (cl Client) Do(r *http.Request, v interface{}) (err error) {
	r.URL, err = cl.baseURL.Parse(r.URL.String())
	if err != nil {
		return err
	}
	res, err := cl.HTTPClient.Do(r)
	if err != nil {
		return err
	}
	defer func() {
		_ = res.Body.Close()
	}()
	return cl.decode(res.Body, v)
}

func (cl *Client) decode(body io.Reader, v interface{}) (err error) {
	if cl.Decoder == nil {
		return errors.New("must init decoder")
	}

	return cl.Decoder(body, v)
}
