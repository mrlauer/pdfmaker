// package document encapsulates simple text documents.
package document

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"sync"
)

// Length represents a page-length value.
// It preserves its original string, or a normalized version of it, to ensure user intent is preserved.
type Length struct {
	definition string
	points     float64
}

// lengthRE is the regular expression for parsing lengths. It is set at initialization.
var lengthRE *regexp.Regexp

func init() {
	decimalString := `\d+(?:\.\d*)?|\.\d+`
	fracString := `(?:\d+(?:\s+|-))?\d+/[1-9]\d*`
	unitString := `("|in|pt|cm|mm)`
	lengthREString := `^\s*(` + decimalString + `|` + fracString + `)\s*` + unitString + `\s*$`
	lengthRE = regexp.MustCompile(lengthREString)
}

func LengthREString() string {
	return lengthRE.String()
}

// an enumeration for unit types.
const (
	_ = iota
	Points
	Inches
	Mils
	Centimeters
	Millimeters
)

// LengthUnit is the "enum" type for, you guessed it, units of length.
type LengthUnit int

// getUnit returns the LengthUnit for a unit string.
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

// normalizedUnitString returns a normalized string for units.
// It converts "in" to "\""
func normalizedUnitString(unitStr string) string {
	switch unitStr {
	case `in`:
		return `"`
	}
	return unitStr
}

// getUnitToPoints returns scale to convert a length in the given units to points.
// If given an invalid unit, it does not return an error or panic; it just returns 1.
// This might be a bad idea.
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

// parseFrac parses a "fraction"
func parseFrac(fracstr string) (string, float64, error) {
	re := regexp.MustCompile(`(?:(\d+)(?:\s+|-))?(\d+)/(\d+)`)
	m := re.FindStringSubmatch(fracstr)
	if m != nil {
		f := 0.0
		normalized := ""
		if m[1] != "" {
			f, _ = strconv.ParseFloat(m[1], 64)
			normalized = m[1] + " "
		}
		num, _ := strconv.ParseFloat(m[2], 64)
		denom, _ := strconv.ParseFloat(m[3], 64)
		f += num / denom
		normalized += m[2] + "/" + m[3]
		return normalized, f, nil
	}
	return fracstr, 0.0, errors.New("Could not parse fraction")

}

// parseLength attempts to parse decimals and fractions
func parseLength(str string) (string, float64, error) {
	l, err := strconv.ParseFloat(str, 64)
	if err == nil {
		return str, l, nil
	}
	return parseFrac(str)

}

// translateLength takes a length string and returns
// - a normalized string
// - the length in points
// - the units
// - an error, in the string is not valid.
func translateLength(def string) (string, float64, LengthUnit, error) {
	if match := lengthRE.FindStringSubmatch(def); match != nil {
		normalized, l, err := parseLength(match[1])
		if err == nil {
			unitStr := match[2]
			unit, err := getUnit(unitStr)
			if err == nil {
				scale := getUnitToPoints(unit)
				normalized += normalizedUnitString(unitStr)
				return normalized, l * scale, unit, err
			}
		}
	}
	return def, 0.0, Points, errors.New("Could not parse length")
}

// LengthFromString returns a Length for a length string. It can fail if the string is invalid.
func LengthFromString(definition string) (Length, error) {
	normalized, points, _, err := translateLength(definition)
	if err != nil {
		return Length{}, err
	}
	return Length{definition: normalized, points: points}, nil
}

// LengthFromPoints returns a Length for a given point value. It always succeeds.
func LengthFromPoints(points float64) Length {
	str := strconv.FormatFloat(points, 'g', -1, 64) + "pt"
	return Length{definition: str, points: points}
}

// String returns the defining string, or "0pt" if there is none.
func (l Length) String() string {
	if l.definition == "" {
		return "0pt"
	}
	return l.definition
}

// Points returns the point value.
func (l Length) Points() float64 {
	return l.points
}

// MarshalJSON uses the defining string.
func (l Length) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.definition)
}

// UnmarshalJSON uses the defining string.
func (l *Length) UnmarshalJSON(data []byte) error {
	var def string
	err := json.Unmarshal(data, &def)
	if err == nil {
		*l, err = LengthFromString(def)
	}
	return err
}

// TODO: use a database, you moron!

// Document encapsulates the defining properties of a document.
type Document struct {
	Font         string
	Text         string
	FontSize     Length
	BaselineSkip Length
	LeftMargin   Length
	RightMargin  Length
	TopMargin    Length
	BottomMargin Length
	PageHeight   Length
	PageWidth    Length
	// Id is the document identifier. It is serialized to JSON is "id", 
	// and omitted if empty.
	Id int `json:"id,omitempty"`
}

// DefaultDocument returns a Document with reasonable default values.
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

type DB interface {
	Add(doc *Document)
	Update(doc *Document) error
	Fetch(id int) (Document, error)
	Delete(id int) error
	Close()
}

// documents is a map that serves as a fake database
type FakeDB struct {
	documents map[int]Document
	docIdx int
	lock sync.RWMutex
}

func CreateFakeDB() *FakeDB {
	db := new(FakeDB)
	db.documents = make(map[int]Document)
	return db
}

// AddDocument adds a new document and sets the id of its argument
func (d *FakeDB)Add(doc *Document) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.docIdx += 1
	doc.Id = d.docIdx
	d.documents[doc.Id] = *doc
}

func (d *FakeDB)Update(doc *Document) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	_, ok := d.documents[doc.Id]
	if ok {
		d.documents[doc.Id] = *doc
		return nil
	}
	return errors.New("document does not exist")
}

func (d *FakeDB)Fetch(id int) (Document, error) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	doc, ok := d.documents[id]
	if ok {
		return doc, nil
	}
	return doc, errors.New("document does not exist")
}

func (d *FakeDB)Delete(id int) error {
	delete(d.documents, id)
	return nil
}

func (d *FakeDB)Close() {
}
