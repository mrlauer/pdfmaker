package main

import (
	"net/http"
	"net/http/httptest"
	"path"
	"runtime"
	"testing"
)

func test_route(t *testing.T, url string, expStatus int) {
	response, err := http.Get(url)
	if err != nil {
		t.Errorf("Could not get %s", url)
		return
	}
	status := response.StatusCode
	if status != expStatus {
		t.Errorf("%s returned status %d", url, status)
	}
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
	test_route(t, base+"/document/", http.StatusOK)
	test_route(t, base+"/document/3", http.StatusNotFound)
	test_route(t, base, http.StatusOK)
}
