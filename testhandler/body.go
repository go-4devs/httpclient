package testhandler

import (
	"io/ioutil"
	"net/http"
)

// Body check json body request
type Body string

// CanHandle check request by body
func (j Body) CanHandle(r *http.Request) bool {
	if body, e := r.GetBody(); e == nil && body != nil {
		b, err := ioutil.ReadAll(body)
		return err == nil && string(b) == string(j)
	}
	return false
}
