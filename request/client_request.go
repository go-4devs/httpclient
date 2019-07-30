package request

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
)

// Encoder for the body
type Encoder func(v interface{}) (io.Reader, error)

// ClientRequest make request by params method query
type ClientRequest struct {
	Encoder Encoder
	Method  string
	URL     string
	user    struct {
		username, password string
	}
	query  url.Values
	header url.Values
	err    error
	body   io.Reader
	ctx    context.Context
}

// NewPost create new post request
func NewPost(ctx context.Context, path string, encoder Encoder) *ClientRequest {
	return NewRequest(ctx, path, http.MethodPost, encoder)
}

// NewGet create new get request
func NewGet(ctx context.Context, path string, encoder Encoder) *ClientRequest {
	return NewRequest(ctx, path, http.MethodGet, encoder)
}

// NewRequest create new request
func NewRequest(ctx context.Context, path, method string, encoder Encoder) *ClientRequest {
	return &ClientRequest{
		ctx:     ctx,
		URL:     path,
		Encoder: encoder,
		Method:  method,
	}
}

// Query add values for the query
func (r *ClientRequest) Query(value ...RValue) *ClientRequest {
	if r.query == nil {
		r.query = url.Values{}
	}
	for _, v := range value {
		v(r.query)
	}
	return r
}

// Header add values for the header
func (r *ClientRequest) Header(value ...RValue) *ClientRequest {
	if r.header == nil {
		r.header = url.Values{}
	}
	for _, v := range value {
		v(r.header)
	}
	return r
}

// SetBasicAuth set username and password basic auth
func (r *ClientRequest) SetBasicAuth(username, password string) *ClientRequest {
	r.user.username = username
	r.user.password = password
	return r
}

// SetBody encode body and add to request
func (r *ClientRequest) SetBody(data interface{}) *ClientRequest {
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
func (r *ClientRequest) HTTP() (*http.Request, error) {
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

func (r *ClientRequest) path() string {

	values := r.query.Encode()
	if values == "" {
		return r.URL
	}

	return r.URL + "?" + values
}
