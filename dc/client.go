// client with decoder
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

var _ httpclient.Fetcher = &Client{}

// Client get response and marshaling it by decoder
type Client struct {
	HTTPClient *http.Client
	decoder    decoder.Decoder
	baseURL    url.URL
	with       func(*http.Response, io.Reader) error
	middleware transport.Middleware
}

// Option for the configure Client
type Option func(*Client)

// WithMiddleware add middleware do request
func WithMiddleware(mw ...transport.Middleware) Option {
	return func(i *Client) {
		if i.middleware != nil {
			mw = append([]transport.Middleware{i.middleware}, mw...)
		}
		if len(mw) > 0 {
			i.middleware = transport.Chain(mw...)
		}
	}
}

// WithTransport set transport
func WithTransport(tr http.RoundTripper) Option {
	return func(i *Client) {
		if i.HTTPClient == http.DefaultClient {
			i.HTTPClient = &http.Client{}
		}
		i.HTTPClient.Transport = tr
	}
}

// WithFetchMiddleware add middleware for transport
// nolint: bodyclose
func WithFetchMiddleware(mw ...func(*http.Response, io.Reader) error) Option {
	return func(i *Client) {
		if i.with != nil {
			mw = append([]func(*http.Response, io.Reader) error{i.with}, mw...)
		}

		i.with = func(response *http.Response, reader io.Reader) error {
			for _, h := range mw {
				if e := h(response, reader); e != nil {
					return e
				}
			}
			return nil
		}
	}
}

// WithErrorMiddleware add middleware for transport
// nolint: bodyclose
func WithErrorMiddleware(minStatusCode int,
	errFactory func() error,
	decoder func(*http.Response, io.Reader, interface{}) error) Option {
	return func(i *Client) {
		WithFetchMiddleware(func(r *http.Response, b io.Reader) (err error) {
			if r.StatusCode >= minStatusCode {
				err = errFactory()
				if derr := decoder(r, b, err); derr != nil {
					return derr
				}
			}

			return
		})(i)
	}
}

// WithDecoder set decoder body
func WithDecoder(decoder decoder.Decoder) Option {
	return func(i *Client) {
		i.decoder = decoder
	}
}

// Must create client or panic
func Must(baseURL string, opts ...Option) *Client {
	cl, err := New(baseURL, opts...)
	if err != nil {
		panic(err)
	}

	return cl
}

// New create new Client with default http client
func New(baseURL string, opts ...Option) (*Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	cl := &Client{
		baseURL:    *u,
		HTTPClient: http.DefaultClient,
	}
	for _, opt := range opts {
		opt(cl)
	}

	if cl.with == nil {
		errDecoder := decoder.HTTPDecode
		if cl.decoder != nil {
			errDecoder = func(r *http.Response, body io.Reader, v interface{}) error {
				return cl.decoder(body, v)
			}
		}
		WithErrorMiddleware(http.StatusBadRequest, apierrors.MessageFactory, errDecoder)(cl)
	}

	return cl, nil
}

// Do request and decode response body
func (c *Client) Do(r *http.Request, v interface{}) error {
	return c.Fetch(r).With(c.with).Decode(v)
}

func (c *Client) Fetch(r *http.Request) httpclient.Fetch {
	f := fetch{
		decode: c.decode,
	}
	r.URL, f.err = c.baseURL.Parse(r.URL.String())
	if f.err != nil {
		return f
	}
	res, err := func(req *http.Request) (*http.Response, error) {
		if c.middleware != nil {
			return c.middleware(r, c.HTTPClient.Do)
		}
		return c.HTTPClient.Do(req)
	}(r)
	if err != nil {
		f.err = err
		return f
	}
	if res.Body != nil {
		var b bytes.Buffer
		if _, err := io.Copy(&b, res.Body); err != nil {
			f.err = err
		}
		if b.Len() != 0 {
			f.body = &b
		}
		_ = res.Body.Close()
	}

	return f
}

type fetch struct {
	body     io.Reader
	response *http.Response
	err      error
	decode   func(r *http.Response, body io.Reader, v interface{}) error
}

func (f fetch) Error() error {
	return f.err
}

func (f fetch) IsStatus(httpStatus int) bool {
	return f.response.StatusCode == httpStatus
}

func (f fetch) With(h func(r *http.Response, b io.Reader) error) httpclient.Fetch {
	if f.err != nil {
		return f
	}
	if err := h(f.response, f.body); err != nil {
		f.err = err
	}
	return f
}

func (f fetch) Decode(v interface{}) error {
	if f.err != nil {
		return f.err
	}
	return f.decode(f.response, f.body, v)
}

func (f fetch) Body() io.Reader {
	return f.body
}

func (c *Client) decode(r *http.Response, body io.Reader, v interface{}) error {
	if c.decoder != nil {
		return c.decoder(body, v)
	}
	return decoder.HTTPDecode(r, body, v)
}
