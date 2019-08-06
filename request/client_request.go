package request

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Encoder for the body
type Encoder func(v interface{}) (io.Reader, error)

// ClientRequest make request by params method query
type ClientRequest struct {
	Encoder  Encoder
	Method   string
	Path     string
	PathArgs []interface{}
	query    url.Values
	err      error
	body     io.Reader
	ctx      context.Context
	mw       func(ctx context.Context, cr *ClientRequest,
		n func(ctx context.Context) (*http.Request, error)) (*http.Request, error)
}

// Option configure client request
type Option func(*ClientRequest)

// WithEncoder set encoder request
func WithEncoder(encoder Encoder) Option {
	return func(request *ClientRequest) {
		request.Encoder = encoder
	}
}

// WithMiddleware set middleware request
func WithMiddleware(mw ...func(ctx context.Context, cr *ClientRequest,
	n func(ctx context.Context) (*http.Request, error)) (*http.Request, error)) Option {
	return func(request *ClientRequest) {
		if request.mw != nil {
			mw = append(mw, request.mw)
		}
		request.mw = chain(mw...)
	}
}

// WithMethod set method by default GET
func WithMethod(method string) Option {
	return func(request *ClientRequest) {
		request.Method = method
	}
}

// NewPost create new post request
func NewPost(ctx context.Context, opts ...Option) ClientRequest {
	return NewRequest(ctx, append(opts, WithMethod(http.MethodPost))...)
}

// NewGet create new get request
func NewGet(ctx context.Context, opts ...Option) ClientRequest {
	return NewRequest(ctx, opts...)
}

// NewRequest create new request
func NewRequest(ctx context.Context, opts ...Option) ClientRequest {
	cl := ClientRequest{
		ctx:    ctx,
		Method: http.MethodGet,
	}
	for _, o := range opts {
		o(&cl)
	}

	return cl
}

// URI set url and args it
func (r ClientRequest) URI(path string, a ...interface{}) ClientRequest {
	r.Path = path
	r.PathArgs = a
	return r
}

// Query add values for the query
func (r ClientRequest) Query(value ...RValue) ClientRequest {
	if r.query == nil {
		r.query = make(url.Values, len(value))
	}
	for _, v := range value {
		v(r.query)
	}
	return r
}

// Header add values for the header
func (r ClientRequest) Header(value ...RValue) ClientRequest {
	return r.handle(func(ctx context.Context, _ *ClientRequest,
		n func(ctx context.Context) (*http.Request, error)) (*http.Request, error) {
		r, err := n(ctx)
		if err == nil {
			for _, v := range value {
				v(r.Header)
			}
		}
		return r, err
	})
}

// SetBasicAuth set username and password basic auth
func (r ClientRequest) SetBasicAuth(username, password string) ClientRequest {
	return r.handle(func(ctx context.Context, _ *ClientRequest,
		n func(context.Context) (*http.Request, error)) (request *http.Request, e error) {
		request, e = n(ctx)
		if e == nil {
			request.SetBasicAuth(username, password)
		}
		return
	})
}

// SetBody encode body and add to request
func (r ClientRequest) SetBody(data interface{}) ClientRequest {
	if r.err != nil {
		return r
	}
	if r.Encoder == nil {
		r.body, r.err = encoder(data)
		return r
	}
	r.body, r.err = r.Encoder(data)

	return r
}

func encoder(v interface{}) (io.Reader, error) {
	var b bytes.Buffer
	switch data := v.(type) {
	case string:
		b.WriteString(data)
	case []byte:
		b.Write(data)
	case io.Reader:
		return data, nil
	default:
		return nil, errors.New("must init encoder for the body")
	}

	return &b, nil
}

// HTTP create http Request
func (r ClientRequest) HTTP() (httpRequest *http.Request, err error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.ctx == nil {
		r.ctx = context.Background()
	}
	if r.mw != nil {
		httpRequest, err = r.mw(r.ctx, &r, r.init)
	} else {
		httpRequest, err = r.init(r.ctx)
	}

	if err != nil {
		return nil, err
	}

	return httpRequest, nil
}

func (r ClientRequest) init(ctx context.Context) (request *http.Request, e error) {
	request, e = http.NewRequest(r.Method, r.path(), r.body)
	if e == nil {
		request = request.WithContext(ctx)
	}

	return request, e
}

func (r ClientRequest) path() string {
	u := r.Path
	if len(r.PathArgs) > 0 {
		u = fmt.Sprintf(r.Path, r.PathArgs...)
	}

	if len(r.query) > 0 {
		return u + "?" + r.query.Encode()
	}

	return u
}

func (r ClientRequest) handle(
	h func(context.Context, *ClientRequest, func(context.Context) (*http.Request, error),
) (*http.Request, error)) ClientRequest {
	if r.mw == nil {
		r.mw = h
	} else {
		r.mw = chain(r.mw, h)
	}
	return r
}

// chain middleware
func chain(handleFunc ...func(ctx context.Context, cr *ClientRequest,
	n func(ctx context.Context) (*http.Request, error)) (*http.Request, error),
) func(ctx context.Context, cr *ClientRequest,
	n func(ctx context.Context) (*http.Request, error)) (*http.Request, error) {
	n := len(handleFunc)
	if n > 1 {
		lastI := n - 1
		return func(ctx context.Context, cr *ClientRequest,
			n func(context.Context) (*http.Request, error)) (*http.Request, error) {
			var (
				chainHandler func(context.Context) (*http.Request, error)
				curI         int
			)
			chainHandler = func(currentCtx context.Context) (*http.Request, error) {
				if curI == lastI {
					return n(currentCtx)
				}
				curI++
				res, err := handleFunc[curI](currentCtx, cr, chainHandler)
				curI--
				return res, err

			}
			return handleFunc[0](ctx, cr, chainHandler)
		}
	}

	if n == 1 {
		return handleFunc[0]
	}

	return func(ctx context.Context, cr *ClientRequest,
		n func(context.Context) (*http.Request, error)) (*http.Request, error) {
		return n(ctx)
	}
}
