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
	transport http.RoundTripper
}

// Option for the configure client
type Option func(*options)

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
	op := &options{
		transport: http.DefaultTransport,
	}
	for _, opt := range opts {
		opt(op)
	}
	cl.client = http.Client{
		Transport: op.transport,
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
