package dc

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/go-4devs/httpclient"
	"github.com/go-4devs/httpclient/transport"
)

// Decoder for decode body
type Decoder func(r io.Reader, v interface{}) error

// HTTPClient get response and marshaling it by decoder
type Client struct {
	HTTPClient http.Client
	Decoder    Decoder
	baseURL    url.URL
	fetch      *fetch
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

// Must create clint or panic
func Must(baseURL string, decoder Decoder, opts ...Option) Client {
	cl, err := New(baseURL, decoder, opts...)
	if err != nil {
		panic(err)
	}

	return cl
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
func (c Client) Do(r *http.Request, v interface{}) (err error) {
	return c.Fetch(r).Decode(v)
}

type fetch struct {
	body     io.Reader
	response *http.Response
	err      error
}

func (c *Client) Error() error {
	if c.fetch.err != nil {
		return c.fetch.err
	}

	return nil
}

func (c *Client) Fetch(r *http.Request) httpclient.Fetch {
	c.fetch = &fetch{}
	r.URL, c.fetch.err = c.baseURL.Parse(r.URL.String())
	if c.fetch.err != nil {
		return c
	}
	res, err := c.HTTPClient.Do(r)
	if err != nil {
		c.fetch.err = err
		return c
	}
	var b bytes.Buffer
	if _, err := io.Copy(&b, c.fetch.response.Body); err != nil {
		c.fetch.err = err
	}
	if b.Len() != 0 {
		c.fetch.body = &b
	}
	_ = res.Body.Close()
	c.fetch.response = res

	return c
}

func (c *Client) IsStatus(httpStatus int) bool {
	if c.fetch != nil {
		return c.fetch.response.StatusCode == httpStatus
	}

	return false
}

func (c *Client) With(h func(r *http.Response, b io.Reader) error) httpclient.Fetch {
	if c.fetch.err != nil {
		if err := h(c.fetch.response, c.fetch.body); err != nil {
			c.fetch.err = err
		}
	}

	return c
}

func (c *Client) Decode(v interface{}) error {
	if c.fetch.err != nil {
		return c.fetch.err
	}

	if c.Decoder == nil {
		return errors.New("must init decoder")
	}

	return c.Decoder(c.fetch.body, c)
}

func (c *Client) Body() io.Reader {
	return c.fetch.body
}
