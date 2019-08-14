package decoder

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	decoder := func(r io.Reader, v interface{}) error {
		return nil
	}
	require.EqualError(t, Register(decoder), "http client: decider and media types is required")
	require.EqualError(t, Register(nil, "text/html"), "http client: decider and media types is required")
	require.Nil(t, Register(decoder, "text/html"))
	require.EqualError(t, Register(decoder, "text/html"),
		"http client: register called twice for decoder by media type text/html")
	require.Nil(t, Decode("text/html", nil, nil))
}

func TestMustRegister(t *testing.T) {
	decoder := func(r io.Reader, v interface{}) error {
		return nil
	}
	MustRegister(decoder, "multipart/form-data")
	require.Nil(t, Decode("multipart/form-data", nil, nil))
	defer func() {
		require.NotNil(t, recover())
	}()
	MustRegister(decoder, "multipart/form-data")
}

func TestDecode(t *testing.T) {
	MustRegister(func(r io.Reader, v interface{}) error {
		return errors.New("error decode")
	}, "application/msword")
	require.EqualError(t, Decode("application/pdf", nil, nil),
		"http client: decoder by media type 'application/pdf' not found")
	require.EqualError(t, Decode("application/msword", nil, nil), "error decode")

	MustRegister(func(r io.Reader, v interface{}) error {
		return xml.NewDecoder(r).Decode(v)
	}, "application/xml")
	var data struct {
		Title string `xml:"title"`
	}
	b := bytes.NewBufferString(`<xml><title>some text</title></xml>`)
	require.Nil(t, Decode("application/xml", b, &data))
	require.Equal(t, "some text", data.Title)
}

func TestHTTPDecode(t *testing.T) {
	r := &http.Response{}
	require.EqualError(t, HTTPDecode(r, nil, nil), "mime: no media type")
	MustRegister(func(r io.Reader, v interface{}) error {
		return nil
	}, "audio/aac")
	r.Header = http.Header{}
	r.Header.Add("Content-Type", "audio/aac")
	require.Nil(t, HTTPDecode(r, nil, nil))
}
