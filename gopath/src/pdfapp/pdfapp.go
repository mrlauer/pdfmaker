/*
The application.

The routes should be, but are not yet
 - GET /static/{filename}	Get a static file.
 - GET /edit/				Edit a new document.
 - GET /					Same as /edit/.
 - GET /edit/{id}/			Edit an existing document.
 - GET /document/			Get a json list of all documents(?).
 - POST /document/			Create a new document with json provided in body.
 - GET /document/{id}/		Get an existing document in json form.
 - PUT /document/{id}/		Update an existing document with json in body.
 - DELETE /document/{id}/	Delete an existing document.
 - GET /pdf/{id}/			Get the pdf for an existing document.
Perhaps these should also switch on Accept headers.
*/
package main

import (
	"code.google.com/p/gorilla/mux"
	"db"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"local/document"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"textproc"
	"web"
)

// TemplateDir is the runtime directory for templates
var TemplateDir string

// StaticDir is the runtime directory for static files
var StaticDir string

// DB is the database
var DB document.DB

// assignId returns the (possibly null, of course) docId for the request.
func assignId(r *http.Request) db.Id {
	var id string
	web.AssignTo(&id, mux.Vars(r)["Id"])
	return db.MakeId(id)
}

// handler is the handler for the basic page.
// It simply redirects to an empty edit page
func handler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, `/edit/`, 301)
}

// editHandler is the handler for showing document edit pages
func editHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s %s\n", r.Method, r.URL.Path)
	fmt.Printf("%s\n", r.Header.Get("Accept"))
	templatePath := path.Join(TemplateDir, "main.html")
	contentPath := path.Join(TemplateDir, "content.html")
	templ, err := template.ParseFiles(templatePath, contentPath)
	header := w.Header()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fonts := textproc.ListFontFamilies()
	sort.Strings(fonts)
	defaultDoc := document.DefaultDocument()
	defaultDocJSON, err := json.Marshal(defaultDoc)
	if err != nil {
		panic(err)
	}
	fontsJSON, err := json.Marshal(fonts)
	if err != nil {
		panic(err)
	}
	id := assignId(r)
	doc, err := DB.Fetch(id)
	if !id.IsNull() && err != nil {
		web.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	docJSON, err := json.Marshal(doc)
	if err != nil {
		panic(err)
	}
	lengthREString := document.LengthREString()
	data := map[string]interface{}{"fonts": template.JS(fontsJSON),
		"doc":        template.JS(docJSON),
		"defaultDoc": template.JS(defaultDocJSON),
		"lengthRE":   lengthREString}
	header.Set("Content-Type", "text/html")
	err = templ.Execute(w, data)
	if err != nil {
		web.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// pdfhandler makes a pdf file out of the information it is passed.
func pdfhandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s %s\n", r.Method, r.URL.Path)
	fmt.Printf("%s\n", r.Header.Get("Accept"))
	header := w.Header()

	id := assignId(r)
	doc, err := DB.Fetch(id)
	if err != nil {
		web.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	header.Set("Content-Type", "application/pdf")
	//header.Set("Content-Disposition", "attachment;filename=foo.pdf")
	pdf := textproc.MakePDFStreamTextObject(w, 8.5*72, 11*72)
	defer pdf.Close()
	props := textproc.TypesettingProps{}
	props.Fontname = doc.Font
	props.Fontsize = doc.FontSize.Points()
	props.Baselineskip = doc.BaselineSkip.Points()
	props.PageWidth = doc.PageWidth.Points()
	props.PageHeight = doc.PageHeight.Points()
	props.LeftMargin = doc.LeftMargin.Points()
	props.RightMargin = doc.RightMargin.Points()
	pdf.WriteAt(doc.Text, props, props.LeftMargin, doc.TopMargin.Points()+props.Fontsize)
}

func writeDoc(w http.ResponseWriter, doc *document.Document) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(doc)
	if err != nil {
		web.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

func postDocHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s %s\n", r.Method, r.URL.Path)
	doc := document.Document{}
	json.NewDecoder(r.Body).Decode(&doc)
	err := DB.Add(&doc)
	if err != nil {
		// Try to figure out what the error was
		web.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeDoc(w, &doc)
}

func putDocHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s %s\n", r.Method, r.URL.Path)
	id := assignId(r)

	doc := document.Document{}
	json.NewDecoder(r.Body).Decode(&doc)

	if !id.IsNull() {
		doc.Id = id
	}
	DB.Update(&doc)
	writeDoc(w, &doc)
}

func getDocHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s %s\n", r.Method, r.URL.Path)
	id := assignId(r)
	doc, err := DB.Fetch(id)
	if err != nil {
		web.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeDoc(w, &doc)
}

func deleteDocHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s %s\n", r.Method, r.URL.Path)
	id := assignId(r)
	err := DB.Delete(id)
	if err != nil {
		web.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	filename := mux.Vars(r)["Filename"]
	http.ServeFile(w, r, path.Join(StaticDir, filename))
}

func panicHandler(w http.ResponseWriter, r *http.Request) {
	panic(errors.New("Oh my stars and whiskers!"))
}

func GetAppDir() string {
	apppath, err := exec.LookPath(os.Args[0])
	if err != nil {
		panic(err)
	}
	apppath, err = filepath.Abs(apppath)
	if err != nil {
		panic(err)
	}
	dir, _ := path.Split(apppath)
	return dir
}

func SetupDB(dbname string) document.DB {
	var err error
	DB, err = document.CreateMongoDB("localhost", dbname)
	if err != nil {
		panic(err)
	}
	return DB
}

func MakeRouter() http.Handler {
	r := web.MakeRouter(TemplateDir)
	r.HandleFunc(`/pdf/{Id}/`, pdfhandler).Methods("GET")
	r.HandleFunc("/static/{Filename:.*}", staticHandler).Methods("GET")
	r.HandleFunc("/", editHandler).Methods("GET")

	docPrefix := r.PathPrefix(`/document/`).Subrouter()
	docPrefix.HandleFunc(`/`, postDocHandler).Methods("POST")
	idr := docPrefix.PathPrefix(`/{Id}/`).Subrouter()
	idr.HandleFunc(`/`, putDocHandler).Methods("PUT")
	idr.HandleFunc(`/`, getDocHandler).Methods("GET")
	idr.HandleFunc(`/`, deleteDocHandler).Methods("DELETE")

	r.HandleFunc(`/edit/{Id}/`, editHandler).Methods("GET")
	r.HandleFunc(`/panic/`, panicHandler)
	return r
}

func SetPaths(topdir string) {
	TemplateDir = path.Join(topdir, "templates")
	StaticDir = path.Join(topdir, "static")
	web.SetTemplateDir(TemplateDir)
}

func main() {
	SetupDB("pdfdb")

	appdir := GetAppDir()
	SetPaths(path.Join(appdir, ".."))

	r := MakeRouter()
	http.Handle("/", r)
	fmt.Printf("listening on localhost:8080\n")
	http.ListenAndServe(":8080", nil)
}
