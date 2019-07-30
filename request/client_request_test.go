package request

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClientRequest_Query(t *testing.T) {
	r := ClientRequest{}
	r.Query(
		StringValue("key", "value"),
		StringValue("key2", "value2"),
		StringValue("key3", " строка"),
		TimeValue("time", requireTime(t, "2019-07-30T08:08:48+03:00"), time.RFC3339),
		Int64Value("key2", 456),
	)

	h, err := r.HTTP()
	require.Nil(t, err)
	require.Equal(t, "?key=value"+
		"&key2=value2&key2=456"+
		"&key3=+%D1%81%D1%82%D1%80%D0%BE%D0%BA%D0%B0&time=2019-07-30T08%3A08%3A48%2B03%3A00", h.URL.String())
}

func TestClientRequest_SetBasicAuth(t *testing.T) {
	r := ClientRequest{}
	r.SetBasicAuth("username", "password")
	h, e := r.HTTP()
	require.Nil(t, e)
	require.Equal(t, "Basic dXNlcm5hbWU6cGFzc3dvcmQ=", h.Header.Get("Authorization"))

	r.SetBasicAuth("username2", "")
	h, e = r.HTTP()
	require.Nil(t, e)
	require.Equal(t, "Basic dXNlcm5hbWUyOg==", h.Header.Get("Authorization"))
}

func TestClientRequest_Header(t *testing.T) {
	r := ClientRequest{}
	r.Header(
		StringValue("key", "value"),
		Int64Value("key2int", 123),
		TimeValue("key3time", requireTime(t, "2019-07-30T08:08:48+03:00"), time.RFC3339),
	)
	h, e := r.HTTP()
	require.Nil(t, e)
	require.Equal(t, "value", h.Header.Get("key"))
	require.Equal(t, "123", h.Header.Get("key2int"))
	require.Equal(t, "2019-07-30T08:08:48+03:00", h.Header.Get("key3time"))
}

func TestClientRequest_HTTP_InvalidMethod(t *testing.T) {
	r := ClientRequest{Method: "my awesome method"}
	h, err := r.HTTP()
	require.Nil(t, h)
	require.Error(t, err)
}

func TestClientRequest_SetBody(t *testing.T) {
	r := ClientRequest{}
	r.SetBody("some data in body")

	h, e := r.HTTP()
	require.Nil(t, e)
	b, e := ioutil.ReadAll(h.Body)
	require.Nil(t, e)
	require.Equal(t, "some data in body", string(b))
}

func TestNewRequest(t *testing.T) {
	ctx := context.Background()
	r, e := NewRequest(ctx, "/path/to/the/page", http.MethodDelete, func(v interface{}) (reader io.Reader, e error) {
		return nil, nil
	}).
		Header(StringValue("x-header", "data")).
		Query(Int64Value("id", 42)).
		SetBody([]byte(`some data`)).
		HTTP()

	require.Nil(t, e)
	ex, e := http.NewRequest(http.MethodDelete, "/path/to/the/page?id=42", nil)
	require.Nil(t, e)
	ex.Header.Add("x-header", "data")

	require.Equal(t, ex, r)
}

func requireTime(t *testing.T, value string) time.Time {
	ti, e := time.Parse(time.RFC3339, value)
	require.Nil(t, e)

	return ti
}