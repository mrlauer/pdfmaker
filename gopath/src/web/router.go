package web

import (
	"code.google.com/p/gorilla/mux"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strings"
)

// Router wraps the gorilla router to do a few standard things.
//
// - Handle end slashes: the desired behavior is that for GET requests
// (and only get requests) if the url does not end with a slash, but
// matches something that does, to redirect to slashed one.
//
// - Convert panics in handlers to Server Error responses.
//
// - Put up standard not found/server error messages
type Router struct {
	grouter     *mux.Router
	templateDir string
}

func (r *Router) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route {
	return r.grouter.HandleFunc(path, f)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			msg := fmt.Sprintf("%v", err)
			Error(w, msg, http.StatusInternalServerError)
		}
	}()
	r.grouter.ServeHTTP(w, req)
}

func MakeRouter(templateDir string) *Router {
	grouter := mux.NewRouter()
	notFoundHandler := func(w http.ResponseWriter, req *http.Request) {
		path := req.URL.Path
		if !strings.HasSuffix(path, "/") && req.Method == "GET" {
			url := *req.URL
			url.Path += "/"
			newReq, err := http.NewRequest("GET", url.String(), nil)
			var rm mux.RouteMatch
			if err == nil && grouter.Match(newReq, &rm) {
				http.Redirect(w, req, url.String(), http.StatusMovedPermanently)
				return
			}
		}
		// Fall through and call the regular one
		NotFoundHandler(w, req)
	}
	grouter.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	return &Router{grouter: grouter, templateDir: templateDir}
}

func Error(w http.ResponseWriter, errmsg string, code int) {
	templatePath := path.Join(TemplateDir, "error.html")
	templ, err := template.ParseFiles(templatePath)
	if err != nil {
		http.Error(w, errmsg, code)
		return
	}
	header := w.Header()
	data := make(map[string]interface{})
	data["code"] = code
	data["errmsg"] = errmsg
	header.Set("Content-Type", "text/html")
	w.WriteHeader(code)
	err = templ.Execute(w, data)
	if err != nil {
		http.Error(w, errmsg, code)
		return
	}
}

func NotFoundHandler(w http.ResponseWriter, req *http.Request) {
	Error(w, "404 Not Found", http.StatusNotFound)
}

var TemplateDir string

func SetTemplateDir(dir string) {
	TemplateDir = dir
}
