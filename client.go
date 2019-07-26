package httpclient

import (
	"io"
	"net/http"
	"net/url"
)

// Client interface for the get response and marshaling it
type Client interface {
	Do(r *http.Request, v interface{}) error
}

// Decoder for decode body
type Decoder func(r io.Reader, v interface{}) error

// BaseClient interface for the get response and marshaling it
type BaseClient struct {
	baseURL *url.URL
	client  http.Client
	decoder Decoder
}

type options struct {
	transport  http.RoundTripper
	middleware TransportMiddleware
}

// Option for the configure client
type Option func(*options)

// WithTransportMiddleware add middleware for transport
func WithTransportMiddleware(mw ...TransportMiddleware) Option {
	return func(i *options) {
		if i.middleware != nil {
			mw = append([]TransportMiddleware{i.middleware}, mw...)
		}
		if len(mw) > 0 {
			i.middleware = chain(mw...)
		}
	}
}

// New create new client
func New(baseURL string, decoder Decoder, opts ...Option) (client BaseClient, err error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return client, err
	}

	cl := BaseClient{
		baseURL: u,
		decoder: decoder,
	}
	op := &options{}
	for _, opt := range opts {
		opt(op)
	}
	transport := http.DefaultTransport
	if op.transport == nil {
		transport = op.transport
	}

	if op.middleware != nil {
		transport = NewTransportMiddleware(transport, op.middleware)
	}
	cl.client = http.Client{
		Transport: transport,
	}
	return cl, nil

}

// Do request and decode response body
func (cl BaseClient) Do(r *http.Request, v interface{}) (err error) {
	r.URL, err = cl.baseURL.Parse(r.URL.String())
	if err != nil {
		return err
	}
	res, err := cl.client.Do(r)
	if err != nil {
		return err
	}
	defer func() {
		_ = res.Body.Close()
	}()
	return cl.decoder(res.Body, v)
}
