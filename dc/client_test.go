package dc

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
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

func isHandle(r *http.Request, uri, method string) bool {
	return r.URL.String() == uri && r.Method == method
}
func TestClient_Do(t *testing.T) {

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		switch {
		case isHandle(r, "/api/ok.json", http.MethodGet):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err = w.Write([]byte(`{"ok":true}`))
			require.Nil(t, err)
		case isHandle(r, "/err/decoder/not/found.html", http.MethodGet):
			w.WriteHeader(http.StatusOK)
			_, err = w.Write([]byte(`<title>decoder not found</title>`))
		case isHandle(r, "/api/empty/body.json", http.MethodGet):
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
		require.Nil(t, err)
	}))
	defer s.Close()
	c := Client{}
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
		r, err := http.NewRequest(http.MethodGet, s.URL+"/err/decoder/not/found.html", nil)
		require.Nil(t, err)
		require.EqualError(t, c.Do(r, &jsonOk), "http client: decoder by content type'text/html; charset=utf-8' not found")
	})
	t.Run("empty body", func(t *testing.T) {
		r, err := http.NewRequest(http.MethodGet, s.URL+"/api/empty/body.json", nil)
		require.Nil(t, err)
		require.Equal(t, ErrEmptyBody, c.Do(r, &jsonOk))
	})
}
