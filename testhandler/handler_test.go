package testhandler

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func ExampleNewHTTPHandler_json() {
	var t *testing.T
	handler := NewHTTPHandler(t,
		NewHandle(URI("/user/1"), `{"id":1}`),
		NewHandle(URI("/user/3"), `{"message":"user with id 2 not found"}`, WithCodeNotFound()),
		NewHandle(URI("/user/wrong"), `{"message":"bad request"}`, WithCodeBadRequest()),
	)
	s := httptest.NewServer(handler)
	defer s.Close()

	res, err := http.Get(s.URL + "/user/1")
	require.Nil(t, err)
	body, err := ioutil.ReadAll(res.Body)
	require.Nil(t, err)
	defer func() {
		_ = res.Body.Close()
	}()
	require.Equal(t, string(body), `{"id":1}`)
}
