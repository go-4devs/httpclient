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
	URL      string
	pathArgs []interface{}
	user     struct {
		username, password string
	}
	query  url.Values
	header url.Values
	err    error
	body   io.Reader
	ctx    context.Context
}

// NewPost create new post request
func NewPost(ctx context.Context, encoder Encoder) ClientRequest {
	return NewRequest(ctx, http.MethodPost, encoder)
}

// NewGet create new get request
func NewGet(ctx context.Context, encoder Encoder) ClientRequest {
	return NewRequest(ctx, http.MethodGet, encoder)
}

// NewRequest create new request
func NewRequest(ctx context.Context, method string, encoder Encoder) ClientRequest {
	return ClientRequest{
		ctx:     ctx,
		Encoder: encoder,
		Method:  method,
	}
}

// Path set url and args it
func (r ClientRequest) Path(path string, a ...interface{}) ClientRequest {
	r.URL = path
	r.pathArgs = a
	return r
}

// Query add values for the query
func (r ClientRequest) Query(value ...RValue) ClientRequest {
	if r.query == nil {
		r.query = url.Values{}
	}
	for _, v := range value {
		v(r.query)
	}
	return r
}

// Header add values for the header
func (r ClientRequest) Header(value ...RValue) ClientRequest {
	if r.header == nil {
		r.header = url.Values{}
	}
	for _, v := range value {
		v(r.header)
	}
	return r
}

// SetBasicAuth set username and password basic auth
func (r ClientRequest) SetBasicAuth(username, password string) ClientRequest {
	r.user.username = username
	r.user.password = password
	return r
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
func (r ClientRequest) HTTP() (*http.Request, error) {
	if r.err != nil {
		return nil, r.err
	}
	httpRequest, err := http.NewRequest(r.Method, r.path(), r.body)
	if err != nil {
		return nil, err
	}

	if r.user.username != "" {
		httpRequest.SetBasicAuth(r.user.username, r.user.password)
	}

	for n := range r.header {
		httpRequest.Header.Add(n, r.header.Get(n))
	}

	return httpRequest, nil
}

func (r ClientRequest) path() string {
	u := r.URL
	if len(r.pathArgs) > 0 {
		u = fmt.Sprintf(r.URL, r.pathArgs...)
	}

	values := r.query.Encode()
	if values == "" {
		return u
	}

	return u + "?" + values
}
