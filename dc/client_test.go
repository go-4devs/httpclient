package dc

import (
	"log"
	"net/http"
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
