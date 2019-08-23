package json

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	"github.com/go-4devs/httpclient/request"
)

// DefaultEncoder marshal data and create new bytes buffer
var DefaultEncoder request.Encoder = func(v interface{}) (io.Reader, error) {
	buff, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(buff), nil
}

// DefaultOptions default option for the json
var DefaultOptions = []request.Option{
	request.WithEncoder(DefaultEncoder),
	request.WithHeader(
		request.StringValue("Accept", "application/json"),
		request.StringValue("Content-Type", "application/json"),
	),
}

// Post create new post request with json encoder for the body
func Post(ctx context.Context, opts ...request.Option) request.ClientRequest {
	return request.NewPost(ctx, append(opts, DefaultOptions...)...)
}

// Get create new get request with json encoder for the body
func Get(ctx context.Context, opts ...request.Option) request.ClientRequest {
	return request.NewGet(ctx, append(opts, DefaultOptions...)...)
}

// Request create new post request with json encoder for the body
func Request(ctx context.Context, opts ...request.Option) request.ClientRequest {
	return request.NewRequest(ctx, append(opts, DefaultOptions...)...)
}
