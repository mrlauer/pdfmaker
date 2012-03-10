package web

import (
	"testing"
)

func TestTrivialMux(t *testing.T) {
	m := CreateMux("/foo/bar/baz")
	if m.Regexp().String() != `^/foo/bar/baz$` {
		t.Errorf("Mux has regexp %q", m.Regexp().String())
	}
}

func TestMux(t *testing.T) {
	m := CreateMux("/base/:Foo/:Bar/")
	if m.Regexp().String() != `^/base/([^/]*)/([^/]*)/$` {
		t.Errorf("Mux has regexp %q", m.Regexp().String())
	}
	params := m.Matches(`/base/wox/baz/`)
	if params == nil {
		t.Errorf("No params")
	} else {
		foo, fooOk := params["Foo"]
		bar, barOk := params["Bar"]
		if !fooOk || foo != "wox" || !barOk || bar != "baz" {
			t.Errorf("bad parameters: foo = %q, bar = %q", foo, bar)
		}
	}
}
