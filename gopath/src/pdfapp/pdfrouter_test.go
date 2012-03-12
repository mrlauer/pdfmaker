package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"local/document"
	"net/http"
	"net/http/httptest"
	"path"
	"runtime"
	"testing"
)

func do_request(t *testing.T, req *http.Request, expStatus int) []byte {
	response, err := http.DefaultClient.Do(req)
	url := req.URL.Path
	if err != nil {
		t.Errorf("Could not get %s", url)
		return nil
	}
	status := response.StatusCode
	if status != expStatus {
		t.Errorf("%s returned status %d", url, status)
	}
	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	return body
}

func test_get(t *testing.T, url string, expStatus int) []byte {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Errorf("Could not create request %s", url)
		return nil
	}
	return do_request(t, req, expStatus)
}

func get_top_dir() string {
	// Yuck
	_, file, _, _ := runtime.Caller(0)
	dir := path.Dir(file)
	top := path.Join(dir, "..", "..", "..")
	return top
}

func TestRouter(t *testing.T) {
	SetPaths(get_top_dir())

	db := SetupDB("apptestdb")
	defer db.DeleteAll()
	r := MakeRouter()
	server := httptest.NewServer(r)
	defer server.Close()
	base := server.URL
	test_get(t, base+"/document/3/", http.StatusNotFound)
	test_get(t, base, http.StatusOK)

	// Test putting and getting
	doc := document.DefaultDocument()
	text := "This is some text БДЖ"
	doc.Text = text
	fontsz, _ := document.LengthFromString("13pt")
	doc.FontSize = fontsz
	jsonRep, err := json.Marshal(doc)
	if err != nil {
		t.Errorf("Could not marshal json")
	}
	req, err := http.NewRequest("POST", base+"/document/", bytes.NewReader(jsonRep))
	if err != nil {
		t.Errorf("Could not create request")
	}
	body := do_request(t, req, http.StatusOK)
	var doc2 document.Document
	err = json.Unmarshal(body, &doc2)
	if err != nil {
		t.Errorf("Could not unmarshall body")
	}
	if doc2.Text != doc.Text {
		t.Errorf("Returned document had wrong text %q", doc2.Text)
	}
	id := doc2.Id

	// Check that it's in the database
	{
		url := fmt.Sprintf("%s/document/%s/", base, id)
		body = test_get(t, url, http.StatusOK)
		var doc2 document.Document
		err = json.Unmarshal(body, &doc2)
		if err != nil {
			t.Errorf("Could not unmarshall body")
		}
		if doc2.Text != doc.Text {
			t.Errorf("Returned document had wrong text %q", doc2.Text)
		}
		if doc2.FontSize != doc.FontSize {
			t.Errorf("Returned document had wrong fontsize %s", doc2.FontSize)
		}
	}
}
