package httpclient

import (
	"net/http"
)

// Client interface for the get response and marshaling it
type Client interface {
	Do(r *http.Request, v interface{}) error
}
