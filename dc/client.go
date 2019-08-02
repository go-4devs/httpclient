package dc

import (
	"bytes"
	"io"
	"net/http"
	"net/url"

	"github.com/go-4devs/httpclient"
	"github.com/go-4devs/httpclient/apierrors"
	"github.com/go-4devs/httpclient/decoder"
	"github.com/go-4devs/httpclient/transport"
)

var _ httpclient.Fetch = &Client{}

// HTTPClient get response and marshaling it by decoder
type Client struct {
	HTTPClient http.Client
	decoder    decoder.Decoder
	baseURL    url.URL
	fetch      *fetch
	with       []func(r *http.Response, b io.Reader) error
}

type options struct {
	transport  http.RoundTripper
	middleware transport.Middleware
	decoder    decoder.Decoder
	with       []func(r *http.Response, b io.Reader) error
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

// WithFetchMiddleware add middleware for transport
func WithFetchMiddleware(mw ...func(r *http.Response, b io.Reader) error) Option {
	return func(i *options) {
		i.with = append(i.with, mw...)
	}
}

// WithErrorMiddleware add middleware for transport
func WithErrorMiddleware(minStatusCode int, errFactory func() error) Option {
	return func(i *options) {
		i.with = append(i.with, func(r *http.Response, b io.Reader) (err error) {
			if r.StatusCode >= minStatusCode {
				err = errFactory()
				if i.decoder == nil {
					return decoder.HTTPDecode(r, b, err)
				}
				if derr := i.decoder(b, err); derr != nil {
					return derr
				}
			}

			return
		})
	}
}

// WithDecoder set decoder body
func WithDecoder(decoder decoder.Decoder) Option {
	return func(i *options) {
		i.decoder = decoder
	}
}

// Must create clint or panic
func Must(baseURL string, opts ...Option) Client {
	cl, err := New(baseURL, opts...)
	if err != nil {
		panic(err)
	}

	return cl
}

// New create new HTTPClient
func New(baseURL string, opts ...Option) (client Client, err error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return client, err
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

	if len(op.with) == 0 {
		WithErrorMiddleware(http.StatusBadRequest, apierrors.MessageFactory)(op)
	}

	cl := Client{
		baseURL: *u,
		decoder: op.decoder,
		HTTPClient: http.Client{
			Transport: tr,
		},
		with: op.with,
	}

	return cl, nil

}

// Do request and decode response body
func (c *Client) Do(r *http.Request, v interface{}) error {
	f := c.Fetch(r)
	for _, w := range c.with {
		f.With(w)

	}

	return f.Decode(v)
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
		return c
	}
	if err := h(c.fetch.response, c.fetch.body); err != nil {
		c.fetch.err = err
	}
	return c
}

func (c *Client) Decode(v interface{}) error {
	if c.fetch.err != nil {
		return c.fetch.err
	}
	return c.decode(v)
}

func (c *Client) Body() io.Reader {
	return c.fetch.body
}

func (c *Client) decode(v interface{}) error {
	if c.decoder != nil {
		return c.decoder(c.fetch.body, v)
	}
	return decoder.HTTPDecode(c.fetch.response, c.fetch.body, v)
}
