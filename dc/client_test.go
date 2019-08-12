package dc

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-4devs/httpclient/decoder"
	"github.com/stretchr/testify/require"
)

func ExampleNew() {
	c, err := New("https://go-search.org")
	if err != nil {
		log.Fatal(err)
	}
	req, _ := http.NewRequest(http.MethodGet, "/api?action=search&q=httpclient", nil)
	var res struct {
		Query string
		Hits  []struct {
			Name    string
			Package string
			Author  string
		}
	}
	err = c.Do(req, &res)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(res)
}

func ExampleMust() {
	req, _ := http.NewRequest(http.MethodGet, "/api?action=search&q=httpclient", nil)
	var res struct {
		Query string
		Hits  []struct {
			Name    string
			Package string
			Author  string
		}
	}
	err := Must("https://go-search.org").Do(req, &res)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(res)
}

var testDecoder = func(r io.Reader, v interface{}) error {
	return nil
}

var testMiddleware = func(r *http.Request, next func(r *http.Request) (*http.Response, error)) (*http.Response, error) {
	return next(r)
}

type testTransport struct{}

func (t testTransport) RoundTrip(*http.Request) (*http.Response, error) {
	panic("implement me")
}

func isHandle(r *http.Request, uri, method string) bool {
	return r.URL.String() == uri && r.Method == method
}

func TestNew(t *testing.T) {
	c, err := New("https://go-search.org\n")
	require.EqualError(t, err, "parse https://go-search.org\n: net/url: invalid control character in URL")
	require.Nil(t, c)

	c, err = New("https://go-search.org")
	require.Nil(t, err)
	require.Equal(t, http.DefaultClient, c.httpClient)
	u, _ := url.Parse("https://go-search.org")
	require.Equal(t, *u, c.baseURL)
	require.NotNil(t, c.with)

	c, err = New("https://go-search.org", WithDecoder(testDecoder))
	require.Nil(t, err)
	require.NotNil(t, c.decoder)

	hc := &http.Client{}
	c, err = New("https://go-search.org", WithHTTPClient(hc))
	require.Nil(t, err)
	require.Equal(t, hc, c.httpClient)
	require.Nil(t, c.decoder)

	ht := testTransport{}
	c, err = New("https://go-search.org",
		WithTransport(ht),
		WithDecoder(testDecoder),
		WithMiddleware(testMiddleware),
		WithMiddleware(testMiddleware),
	)
	require.Nil(t, err)
	require.NotEqual(t, http.DefaultClient, c.httpClient)
	require.Equal(t, ht, c.httpClient.Transport)
	require.NotNil(t, c.decoder)
	require.NotNil(t, c.middleware)
}

func TestMust(t *testing.T) {
	require.NotNil(t, Must("https://go-search.org"))

	defer func() {
		require.NotNil(t, recover())
	}()
	_ = Must("ya.ru\t\n")
}

func testServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		switch {
		case isHandle(r, "/api/ok.json", http.MethodGet):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err = w.Write([]byte(`{"ok":true}`))
			require.Nil(t, err)
		case isHandle(r, "/index.html", http.MethodGet):
			w.WriteHeader(http.StatusOK)
			_, err = w.Write([]byte(`<title>decoder not found</title>`))
		case isHandle(r, "/api/empty.json", http.MethodGet):
			w.WriteHeader(http.StatusOK)
		case isHandle(r, "/api/not-found.json", http.MethodGet):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_, err = w.Write([]byte(`{"message":"not found"}`))
		case isHandle(r, "/api/invalid.json", http.MethodGet):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(`{"err":invalid}`))
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
		require.Nil(t, err)
	}))
}

func TestClient_Do(t *testing.T) {

	s := testServer(t)
	defer s.Close()

	c := Client{}
	cl := Must(s.URL)

	decoder.MustRegister(func(r io.Reader, v interface{}) error {
		return json.NewDecoder(r).Decode(v)
	}, "application/json")
	var jsonOk struct {
		Ok bool
	}
	t.Run("json ok", func(t *testing.T) {
		r, err := http.NewRequest(http.MethodGet, s.URL+"/api/ok.json", nil)
		require.Nil(t, err)
		require.Nil(t, c.Do(r, &jsonOk))
		require.True(t, jsonOk.Ok)
	})
	t.Run("decoder not found", func(t *testing.T) {
		r, err := http.NewRequest(http.MethodGet, s.URL+"/index.html", nil)
		require.Nil(t, err)
		require.EqualError(t, c.Do(r, &jsonOk), "http client: decoder by content type'text/html; charset=utf-8' not found")
	})
	t.Run("empty body", func(t *testing.T) {
		r, err := http.NewRequest(http.MethodGet, s.URL+"/api/empty.json", nil)
		require.Nil(t, err)
		require.Equal(t, ErrEmptyBody, c.Do(r, &jsonOk))
		require.Equal(t, ErrEmptyBody, cl.Do(r, &jsonOk))
	})
	t.Run("not found", func(t *testing.T) {
		r, err := http.NewRequest(http.MethodGet, s.URL+"/api/not-found.json", nil)
		require.Nil(t, err)
		require.Nil(t, c.Do(r, &jsonOk))
		require.Equal(t, struct{ Ok bool }{Ok: false}, jsonOk)
		require.EqualError(t, cl.Do(r, &jsonOk), "not found")
	})
	t.Run("invalid json", func(t *testing.T) {
		r, err := http.NewRequest(http.MethodGet, s.URL+"/api/invalid.json", nil)
		require.Nil(t, err)
		require.EqualError(t, c.Do(r, &jsonOk), "invalid character 'i' looking for beginning of value")
	})
}
