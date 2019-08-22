package timeout

import (
	"errors"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/go-4devs/httpclient/transport"
	"github.com/stretchr/testify/require"
)

func ExampleNew() {
	mw := New(time.Nanosecond * 200)

	cl := http.Client{
		Transport: transport.NewMiddleware(http.DefaultTransport, mw),
	}
	r, err := cl.Get("http://google.com")
	if err != nil {
		//&url.Error{Op:"Get", URL:"http://google.com", Err:context.deadlineExceededError{}}
		log.Fatalf("%#v", err)
	}
	defer r.Body.Close()
	log.Print(r)
}

// nolint: bodyclose
func TestNew(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://google.com", nil)

	r := New(time.Millisecond)
	d, e := r(req, func(r *http.Request) (*http.Response, error) {
		select {
		case <-r.Context().Done():
			return nil, errors.New("cancel")
		case <-time.After(time.Second):
		}

		return nil, errors.New("failed")
	})
	require.Nil(t, d)
	require.EqualError(t, e, "cancel")
}
