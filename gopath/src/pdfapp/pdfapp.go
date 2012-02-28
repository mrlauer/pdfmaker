package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
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

var TemplateDir string
var StaticDir string

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
	defaultDoc := DefaultDocument()
	defaultDocJSON, err := json.Marshal(defaultDoc)
	if err != nil {
		panic(err)
	}
	fontsJSON, err := json.Marshal(fonts)
	if err != nil {
		panic(err)
	}
	data := map[string]interface{}{"fonts": template.JS(fontsJSON), "defaultDoc" : template.JS(defaultDocJSON)}
	header.Set("Content-Type", "text/html")
	err = templ.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

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

type Length struct {
	definition string
	points	   float64
}

var lengthRE *regexp.Regexp

func init() {
	lengthRE = regexp.MustCompile(`^\s*(\d+(\.\d*)?|\.\d+)\s*("|in|pt)\s*$`)
}

// units
const (
	_ = iota
	Points
	Inches
	Mils
	Centimeters
	Millimeters
)

type LengthUnit int

func getUnit(unitStr string) (LengthUnit, error) {
	switch unitStr {
	case `"`:
		return Inches, nil
	case `in`:
		return Inches, nil
	case `mil`:
		return Mils, nil
	case `pt`:
		return Points, nil
	case `cm`:
		return Centimeters, nil
	case `mm`:
		return Millimeters, nil
	}
	fmt.Printf("Invalid unit string %s\n", unitStr)
	return Points, errors.New("Invalid unit string")
}

func normalizedUnitString(unitStr string) string {
	switch unitStr {
	case `in`:
		return `"`
	}
	return unitStr
}

func getUnitToPoints(unit LengthUnit) float64 {
	switch unit {
	case Inches:
		return 72.0
	case Mils:
		return 0.072
	case Centimeters:
		return 72.0 / 2.54
	case Millimeters:
		return 72.0 / 25.4
	}
	return 1.0
}

func translateLength(def string) (string, float64, LengthUnit, error) {
	if match := lengthRE.FindStringSubmatch(def); match != nil {
		l, err := strconv.ParseFloat(match[1], 64)
		if err == nil {
			unitStr := match[3]
			unit, err := getUnit(unitStr)
			if err == nil {
				scale := getUnitToPoints(unit)
				normalized := match[1] + normalizedUnitString(unitStr)
				return normalized, l * scale, unit, err
			}
		}
	}
	return def, 0.0, Points, errors.New("Could not parse length")
}

func LengthFromString(definition string) (Length, error) {
	normalized, points, _, err := translateLength(definition)
	if err != nil {
		return Length{}, err
	}
	return Length{definition: normalized, points: points}, nil
}

func LengthFromPoints(points float64) Length {
	str := strconv.FormatFloat(points, 'g', -1, 64) + "pt"
	return Length{definition: str, points: points}
}

func (l Length) String() string {
	if l.definition == "" {
		return "0pt"
	}
	return l.definition
}

func (l Length) Points() float64 {
	return l.points
}

// implements json marshal/unmarshall
func (l Length) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.definition)
}

func (l *Length) UnmarshalJSON(data []byte) error {
	var def string
	err := json.Unmarshal(data, &def)
	if err == nil {
		*l, err = LengthFromString(def)
	}
	return err
}

// TODO: use a database, you moron!
// This is the structure that is translated to/from JS
type Document struct {
	Font		 string
	Text		 string
	FontSize	 Length
	BaselineSkip Length
	LeftMargin	 Length
	RightMargin  Length
	TopMargin	 Length
	BottomMargin Length
	PageHeight	 Length
	PageWidth	 Length
	Id			 int `json:"id,omitempty"`
}

func DefaultDocument() *Document {
	doc := Document{}
	doc.Font = "Adobe Garamond Pro"
	doc.Text = "Lorem Ipsum"
	doc.FontSize = LengthFromPoints(12)
	doc.BaselineSkip = LengthFromPoints(15)
	doc.LeftMargin = LengthFromPoints(72)
	doc.RightMargin = LengthFromPoints(72)
	doc.TopMargin = LengthFromPoints(72)
	doc.BottomMargin = LengthFromPoints(72)
	doc.PageHeight, _ = LengthFromString(`11in`)
	doc.PageWidth, _ = LengthFromString(`8.5"`)
	return &doc
}

func writeDoc(w http.ResponseWriter, doc *Document) {
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

	doc := Document{}
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
			pdoc = DefaultDocument()
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
