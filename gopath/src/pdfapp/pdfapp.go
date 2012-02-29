package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"local/document"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"textproc"
)

// TemplateDir is the runtime directory for templates
var TemplateDir string

// StaticDir is the runtime directory for static files
var StaticDir string

// handler is the handler for the basic page.
func handler(w http.ResponseWriter, r *http.Request) {
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
	lengthREString := document.LengthREString()
	data := map[string]interface{}{"fonts": template.JS(fontsJSON),
		"defaultDoc": template.JS(defaultDocJSON),
		"lengthRE":   lengthREString}
	header.Set("Content-Type", "text/html")
	err = templ.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// pdfhandler makes a pdf file out of the information it is passed.
func pdfhandler(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	r.ParseForm()
	text := r.Form.Get("text")
	font := r.Form.Get("font")
	fontsz := 12.0
	topmargin := 72.0
	leftmargin := 72.0

	header.Set("Content-Type", "application/pdf")
	pdf := textproc.MakePDFStreamTextObject(w, 8.5*72, 11*72)
	props := textproc.TypesettingProps{Fontname: font, Fontsize: 12.0, Baselineskip: 15.0}
	props.PageWidth = 72.0 * 8.5
	props.LeftMargin = leftmargin
	props.RightMargin = leftmargin
	pdf.WriteAt(text, props, leftmargin, topmargin+fontsz)
	pdf.Close()
}

func writeDoc(w http.ResponseWriter, doc *document.Document) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(doc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

func docHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s %s\n", r.Method, r.URL.Path)
	fmt.Printf("%s\n", r.Header.Get("Accept"))
	re := regexp.MustCompile(`/\w*(/(\w+))?/?`)
	id := re.FindStringSubmatch(r.URL.Path)[2]

	doc := document.Document{}
	json.NewDecoder(r.Body).Decode(&doc)

	if !(id == "0" || id == "") {
		if id64, err := strconv.ParseInt(id, 10, 32); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else {
			doc.Id = int(id64)
		}
	}

	switch r.Method {
	case "POST":
		doc.Id = 37
		writeDoc(w, &doc)
	case "GET":
		pdoc := &doc
		if id == "0" || id == "" {
			// Getting default values
			pdoc = document.DefaultDocument()
		}
		writeDoc(w, pdoc)
	case "PUT":
		writeDoc(w, &doc)
	}

}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	re := regexp.MustCompile(`/static/(.*)`)
	filename := re.FindStringSubmatch(r.URL.Path)[1]
	http.ServeFile(w, r, path.Join(StaticDir, filename))
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
	appdir := GetAppDir()
	TemplateDir = path.Join(appdir, "../templates")
	StaticDir = path.Join(appdir, "../static")
	http.HandleFunc("/pdf", pdfhandler)
	http.HandleFunc("/static/", staticHandler)
	http.HandleFunc("/", handler)
	http.HandleFunc("/document/", docHandler)
	fmt.Printf("listening on localhost:8080\n")
	http.ListenAndServe(":8080", nil)
}
