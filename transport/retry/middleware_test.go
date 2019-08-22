package retry

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/go-4devs/httpclient/transport"
	"github.com/stretchr/testify/require"
)

func ExampleNew() {
	mw := New(100,
		WithBackOffLinear(time.Millisecond*100),
		WithStatusCode5XX(),
		WithStatusCode(http.StatusConflict, http.StatusPermanentRedirect),
		WithRetriable(func(response *http.Response) bool {
			return response.ContentLength > 0
		}),
	)

	cl := http.Client{
		Transport: transport.NewMiddleware(http.DefaultTransport, mw),
	}
	r, err := cl.Get("http://google.com")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	log.Print(r)
}

var errResp = errors.New("failed get response")

func testRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "http://google.com", nil)
	return req
}
func requireErrResp(t *testing.T, err error) {
	require.EqualError(t, err, "failed get response")
}

// nolint: bodyclose
func TestNew(t *testing.T) {
	r := New(2, WithBackOffLinear(time.Millisecond*20))
	var cnt int
	d, e := r(testRequest(), func(*http.Request) (*http.Response, error) {
		cnt++
		return nil, errResp
	})

	require.Nil(t, d)
	requireErrResp(t, e)
	require.Equal(t, 3, cnt)

	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(time.Millisecond/10, cancel)
	d, e = r(testRequest().WithContext(ctx), func(*http.Request) (*http.Response, error) {
		cnt++
		return nil, errResp
	})
	require.Nil(t, d)
	require.Equal(t, 4, cnt)
	requireErrResp(t, e)
}

// nolint: bodyclose
func TestWithBackOff(t *testing.T) {
	r := New(3, WithBackOffLinear(time.Second/2))
	start := time.Now()
	var cnt int
	d, e := r(testRequest(), func(*http.Request) (*http.Response, error) {
		cnt++
		return nil, errResp
	})
	require.Nil(t, d)
	require.True(t, time.Since(start) >= time.Second/2)
	requireErrResp(t, e)
}

func TestWithStatusCode(t *testing.T) {
	r := New(100, WithStatusCode5XX(), WithStatusCode(400))
	var cnt int
	d, e := r(testRequest(), func(*http.Request) (*http.Response, error) {
		res := &http.Response{
			Body: ioutil.NopCloser(&bytes.Buffer{}),
		}
		switch cnt {
		case 0:
			res.StatusCode = http.StatusBadRequest
		case 1:
			res.StatusCode = http.StatusInternalServerError
		case 2:
			res.StatusCode = http.StatusNetworkAuthenticationRequired
		default:
			res.StatusCode = http.StatusOK
		}
		cnt++
		return res, nil
	})
	defer d.Body.Close()
	require.Nil(t, e)
	require.Equal(t, http.StatusOK, d.StatusCode)
	require.Equal(t, 4, cnt)
}

// nolint: bodyclose
func TestWithRetriable(t *testing.T) {
	r := New(100)
	var cnt int
	d, e := r(testRequest(), func(*http.Request) (*http.Response, error) {
		cnt++
		return nil, nil
	})
	require.Nil(t, e)
	require.Nil(t, d)
	require.Equal(t, 1, cnt)
}
