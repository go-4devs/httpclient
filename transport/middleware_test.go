package transport

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func testMW(msg string) func(r *http.Request,
	next func(r *http.Request) (*http.Response, error)) (*http.Response, error) {
	return func(r *http.Request,
		next func(r *http.Request) (*http.Response, error)) (response *http.Response, e error) {
		res, err := next(r)
		return res, errors.New(err.Error() + msg)
	}
}
func testErr(t *testing.T, msg string) func(res *http.Response, err error) {
	return func(res *http.Response, err error) {
		require.Nil(t, res)
		require.EqualError(t, err, "err handle"+msg)
	}
}

// nolint: bodyclose
func TestChain(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	require.Nil(t, err)
	handleErr := func(r *http.Request) (*http.Response, error) {
		return nil, errors.New("err handle")
	}
	testErr(t, "")(Chain()(r, handleErr))
	testErr(t, "one")(Chain(testMW("one"))(r, handleErr))
	testErr(t, " three two one")(Chain(
		testMW(" one"),
		testMW(" two"),
		testMW(" three"),
	)(r, handleErr))
}

type testRoundTrip struct{}

func (t testRoundTrip) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, errors.New("err handle")
}

// nolint: bodyclose
func TestNewMiddleware(t *testing.T) {
	mw := NewMiddleware(testRoundTrip{}, testMW(""))

	testErr(t, "")(mw.RoundTrip(nil))
}
