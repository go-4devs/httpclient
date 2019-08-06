package testhandler

import "net/http"

//URI base url Request
type URI string

//CanHandle check url by Request
func (u URI) CanHandle(r *http.Request) bool {
	return r.URL.String() == string(u)
}
