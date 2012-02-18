package main

import(
	"net/http"
	"fmt"
	"textproc"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"html/template"
)

var TemplateDir string

func handler(w http.ResponseWriter, r *http.Request) {
	templatePath := path.Join(TemplateDir, "main.html")
	contentPath := path.Join(TemplateDir, "content.html")
	templ, err := template.ParseFiles(templatePath, contentPath)
	header := w.Header()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fonts := textproc.ListFontFamilies()
	data := map[string]interface{} { "text" : "Ohai there!", "fonts" : fonts }
	header.Set("Content-Type", "text/html")
	err = templ.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func pdfhandler(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	header.Set("Content-Type", "application/pdf")
	pdf := textproc.MakePDFStreamTextObject(w, 8.5 * 72, 11 * 72)
	props := textproc.TypesettingProps{Fontname:"Adobe Garamond Pro", Fontsize:12.0, Baselineskip:15.0}
	pdf.WriteAt("Ohai there", props, 10.0, 15.0)
	pdf.Close()
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
	http.HandleFunc("/pdf", pdfhandler)
	http.HandleFunc("/", handler)
	fmt.Printf("listening on localhost:8080\n")
	http.ListenAndServe(":8080", nil)
}
