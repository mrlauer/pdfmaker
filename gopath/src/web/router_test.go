package web

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestRouter(t *testing.T) {
	responseString := "Ohai!"
	router := MakeRouter("")
	router.HandleFunc("/foo/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(responseString))
	})
	router.HandleFunc("/panic/", func(w http.ResponseWriter, r *http.Request) {
		panic("O Noes!")
	})
	server := httptest.NewServer(router)
	defer server.Close()

	type TestData struct {
		Path   string
		Status int
		Body   string
	}

	data := []TestData{
		TestData{"/foo/", http.StatusOK, responseString},
		TestData{"/foo", http.StatusOK, responseString},
		TestData{"/bar", http.StatusNotFound, "404 Not Found\n"},
		TestData{"/panic", http.StatusInternalServerError, "O Noes!\n"}}

	for _, d := range data {
		url := server.URL + d.Path
		res, err := http.Get(url)
		if err != nil {
			t.Errorf("Could not get %s", url)
		}
		if res.StatusCode != d.Status {
			t.Errorf("Status code was %d", res.StatusCode)
		}
		body, _ := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if string(body) != d.Body {
			t.Errorf("Body was %q", body)
		}
	}
}

func ExampleRouter() {
	responseString := "Ohai!"
	router := MakeRouter("")
	router.HandleFunc("/foo/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(responseString))
	})
	router.HandleFunc("/panic/", func(w http.ResponseWriter, r *http.Request) {
		panic("O Noes!")
	})
	server := httptest.NewServer(router)
	defer server.Close()

	r, _ := http.Get(server.URL + "/foo")
	fmt.Printf("Status: %d\n", r.StatusCode)
	fmt.Printf("Body: ")
	io.Copy(os.Stdout, r.Body)
	r.Body.Close()
	fmt.Printf("\n")

	// Output:
	// Status: 200
	// Body: Ohai!
}
