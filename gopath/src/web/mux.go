package web

import (
	"regexp"
)

type Mux struct {
	text      string
	muxRegexp *regexp.Regexp
	variables []string
}

func (m *Mux)Text() string {
	return m.text
}

func (m *Mux)Regexp() *regexp.Regexp {
	return m.muxRegexp
}

func CreateMux(text string) *Mux {
	// parse the text into a regexp
	re := regexp.MustCompile(`:[A-Z]\w*`)
	matches := re.FindAllStringIndex(text, -1)

	restring := "^"
	pos := 0
	variables := []string{}
	for _, match := range matches {
		start := match[0]
		end := match[1]
		// Everything since the last match
		// Should check for evil characters
		restring += regexp.QuoteMeta(text[pos:start])
		pos = end
		// Turn this into an RE
		// Should the RE be more restrictive?
		restring += `([^/]*)`
		variables = append(variables, text[start+1:end])
	}
	restring += regexp.QuoteMeta(text[pos:])
	restring += "$"
	muxre := regexp.MustCompile(restring)

	return &Mux{text: text, muxRegexp: muxre, variables: variables}
}

func (m *Mux) Matches(text string) map[string]string {
	matches := m.Regexp().FindStringSubmatch(text)
	if matches != nil {
		variables := make(map[string]string)
		for i, v := range m.variables {
			variables[v] = matches[i+1]
		}
		return variables
	}
	return nil
}
