package json

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	"github.com/go-4devs/httpclient/request"
)

// defaultEncoder marshal data and create new bytes buffer
var defaultEncoder request.Encoder = func(v interface{}) (io.Reader, error) {
	buff, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(buff), nil
}

// NewPost create new post request with json encoder for the body
func NewPost(ctx context.Context) request.ClientRequest {
	return request.NewPost(ctx, request.WithEncoder(defaultEncoder)).
		Header(
			request.StringValue("Accept", "application/json"),
			request.StringValue("Content-Type", "application/json"),
		)
}

// NewGet create new post request with json encoder for the body
func NewGet(ctx context.Context) request.ClientRequest {
	return request.NewGet(ctx, request.WithEncoder(defaultEncoder)).
		Header(
			request.StringValue("Accept", "application/json"),
			request.StringValue("Content-Type", "application/json"),
		)
}
