package main

import (
	"code.google.com/p/gorilla/mux"
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
	"strconv"
	"textproc"
	"web"
)

// TemplateDir is the runtime directory for templates
var TemplateDir string

// StaticDir is the runtime directory for static files
var StaticDir string

// DB is the database
var DB document.DB

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
	var id int
	web.AssignTo(&id, mux.Vars(r)["Id"])
	doc, err := DB.Fetch(id)
	if id != 0 && err != nil {
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

	var id int
	web.AssignTo(&id, mux.Vars(r)["Id"])
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

func docHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s %s\n", r.Method, r.URL.Path)
	fmt.Printf("%s\n", r.Header.Get("Accept"))
	id := mux.Vars(r)["Id"]

	doc := document.Document{}
	json.NewDecoder(r.Body).Decode(&doc)

	var id64 int64
	if !(id == "0" || id == "") {
		var err error
		if id64, err = strconv.ParseInt(id, 10, 32); err != nil {
			web.Error(w, err.Error(), http.StatusNotFound)
			return
		} else {
			doc.Id = int(id64)
		}
	}

	switch r.Method {
	case "POST":
		DB.Add(&doc)
		writeDoc(w, &doc)
	case "GET":
		doc2, err := DB.Fetch(int(id64))
		if err != nil || id == "0" || id == "" {
			// Getting default values
			doc2 = *document.DefaultDocument()
		}
		writeDoc(w, &doc2)
	case "PUT":
		DB.Update(&doc)
		writeDoc(w, &doc)
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

func main() {
	var err error
	DB, err = document.CreateMongoDB("localhost", "pdfdb")
	if err != nil {
		panic(err)
	}
	appdir := GetAppDir()
	TemplateDir = path.Join(appdir, "../templates")
	StaticDir = path.Join(appdir, "../static")
	web.SetTemplateDir(TemplateDir)
	r := web.MakeRouter(TemplateDir)
	r.HandleFunc(`/pdf/{Id:\d*}`, pdfhandler)
	r.HandleFunc("/static/{Filename:.*}", staticHandler)
	r.HandleFunc("/", editHandler)
	r.HandleFunc(`/document/{Id:\d*}`, docHandler)
	r.HandleFunc(`/edit/{Id:\d*}`, editHandler)
	r.HandleFunc(`/panic/`, panicHandler)
	http.Handle("/", r)
	fmt.Printf("listening on localhost:8080\n")
	http.ListenAndServe(":8080", nil)
}
