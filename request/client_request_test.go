package request

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func clientEncoder(v interface{}) (reader io.Reader, e error) {
	return nil, nil
}

func ExampleClientRequest_Query() {
	ctx := context.TODO()
	req, err := NewRequest(ctx, WithEncoder(clientEncoder)).
		Query(
			StringValue("q", "search"),
			Int64Value("id", 2),
			TimeValue("ts", time.Now(), time.RFC3339Nano),
		).
		HTTP()
	if err != nil {
		log.Fatal(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	fmt.Println(res.Body)

}

func ExampleClientRequest_Path() {
	ctx := context.TODO()
	req, err := NewRequest(ctx, WithEncoder(clientEncoder)).URI("/users/%d", 1).HTTP()
	if err != nil {
		log.Fatal(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	fmt.Println(res.Body)
}

func TestClientRequest_Query(t *testing.T) {
	r := ClientRequest{}
	h, err := r.Query(
		StringValue("key", "value"),
		StringValue("key2", "value2"),
		StringValue("key3", " строка"),
		TimeValue("time", requireTime(t, "2019-07-30T08:08:48+03:00"), time.RFC3339),
		Int64Value("key2", 456),
	).HTTP()

	require.Nil(t, err)
	require.NotNil(t, h)
	require.Equal(t, "?key=value"+
		"&key2=value2&key2=456"+
		"&key3=+%D1%81%D1%82%D1%80%D0%BE%D0%BA%D0%B0&time=2019-07-30T08%3A08%3A48%2B03%3A00", h.URL.String())
}

func TestClientRequest_SetBasicAuth(t *testing.T) {
	r := ClientRequest{}
	h, e := r.SetBasicAuth("username", "password").HTTP()
	require.Nil(t, e)
	require.Equal(t, "Basic dXNlcm5hbWU6cGFzc3dvcmQ=", h.Header.Get("Authorization"))

	h, e = r.SetBasicAuth("username2", "").HTTP()
	require.Nil(t, e)
	require.Equal(t, "Basic dXNlcm5hbWUyOg==", h.Header.Get("Authorization"))
}

func TestClientRequest_Header(t *testing.T) {
	r := ClientRequest{}
	h, e := r.Header(
		StringValue("key", "value"),
		Int64Value("key2int", 123),
		TimeValue("key3time", requireTime(t, "2019-07-30T08:08:48+03:00"), time.RFC3339),
	).HTTP()
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
	h, e := r.SetBody("some data in body").HTTP()
	require.Nil(t, e)
	b, e := ioutil.ReadAll(h.Body)
	require.Nil(t, e)
	require.Equal(t, "some data in body", string(b))
}

func TestNewRequest(t *testing.T) {
	ctx := context.Background()
	r, e := NewRequest(ctx, WithEncoder(func(v interface{}) (io.Reader, error) {
		return nil, nil
	}), WithMethod(http.MethodDelete)).
		URI("/path/to/the/page").
		Header(StringValue("x-header", "data")).
		Query(Int64Value("id", 42)).
		SetBody([]byte(`some data`)).
		HTTP()

	require.Nil(t, e)
	ex, e := http.NewRequest(http.MethodDelete, "/path/to/the/page?id=42", nil)
	require.Nil(t, e)
	ex.Header.Add("x-header", "data")

	require.Equal(t, ex.WithContext(ctx), r)
}

func BenchmarkNewRequest(b *testing.B) {
	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		_, e := NewRequest(ctx, WithEncoder(func(v interface{}) (reader io.Reader, e error) {
			return nil, nil
		}), WithMethod(http.MethodDelete)).
			URI("/path/to/the/page").
			Header(StringValue("x-header", "data")).
			Query(
				Int64Value("id", 42),
				StringValue("id", "42"),
			).
			SetBody([]byte(`some data`)).
			HTTP()

		if e != nil {
			b.Fatalf("Unexpected error: %s", e)
		}
	}
}

func requireTime(t *testing.T, value string) time.Time {
	ti, e := time.Parse(time.RFC3339, value)
	require.Nil(t, e)

	return ti
}
