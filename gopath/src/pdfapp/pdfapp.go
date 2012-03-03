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
	"code.google.com/p/gorilla/mux"
)

// TemplateDir is the runtime directory for templates
var TemplateDir string

// StaticDir is the runtime directory for static files
var StaticDir string

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
	fmt.Printf("%s %s\n", r.Method, r.URL.Path)
	fmt.Printf("%s\n", r.Header.Get("Accept"))
	header := w.Header()

	re := regexp.MustCompile(`/\w*(/(\w+))?/?`)
	idstr := re.FindStringSubmatch(r.URL.Path)[2]
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		fmt.Printf("could not parse, %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	doc, err := document.FetchDocument(int(id))
	if err != nil {
		fmt.Printf("could not find doc %d, %s\n", id, err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
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
		http.Error(w, err.Error(), http.StatusNotFound)
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
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else {
			doc.Id = int(id64)
		}
	}

	switch r.Method {
	case "POST":
		document.AddDocument(&doc)
		writeDoc(w, &doc)
	case "GET":
		doc2, err := document.FetchDocument(int(id64))
		if err != nil || id == "0" || id == "" {
			// Getting default values
			doc2 = *document.DefaultDocument()
		}
		writeDoc(w, &doc2)
	case "PUT":
		document.UpdateDocument(&doc)
		writeDoc(w, &doc)
	}

}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	filename := mux.Vars(r)["Filename"]
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
	r := mux.NewRouter()
	r.HandleFunc(`/pdf/{Id:\d*}`, pdfhandler)
	r.HandleFunc("/static/{Filename:.*}", staticHandler)
	r.HandleFunc("/", editHandler)
	r.HandleFunc(`/document/{Id:\d*}`, docHandler)
	r.HandleFunc(`/edit/{Id:\d*}`, editHandler)
	http.Handle("/", r)
	fmt.Printf("listening on localhost:8080\n")
	http.ListenAndServe(":8080", nil)
}
