package web

import (
	"testing"
)

type testStruct struct {
	Int     int
	Int32   int32
	Int64   int64
	Float32 float32
	Float54 float64
	Bool1   bool
	Bool2   bool
	String  string
}

func TestAssignToStruct(t *testing.T) {
	testVals := map[string]string{"Int": "1",
		"Int32":   "32",
		"Int64":   "64",
		"Float32": "3.14159",
		"Float64": "2e27",
		"Bool1":   "true",
		"Bool2":   "false",
		"String":  "hello world"}

	var s testStruct
	err := AssignToStruct(&s, testVals)
	if err != nil {
		t.Errorf("assignToStruct returned an error")
	}
	if s.Int != 1 {
		t.Errorf("int was %q", s.Int)
	}
	if s.Int32 != 32 {
		t.Errorf("int was %q", s.Int32)
	}
	if s.Int64 != 64 {
		t.Errorf("int was %q", s.Int64)
	}
	if s.Float32 != 3.14159 {
		t.Errorf("float32 was %q", s.Float32)
	}
	if s.Bool1 != true {
		t.Errorf("bool1 was %q", s.Bool1)
	}
	if s.Bool2 != false {
		t.Errorf("bool2 was %q", s.Bool2)
	}
	if s.String != "hello world" {
		t.Errorf("string was %q", s.String)
	}
}

func TestAssignTo(t *testing.T) {
	var i int
	AssignTo(&i, "42")
	if i != 42 {
		t.Errorf("int was %q", i)
	}

	var i32 int32
	AssignTo(&i32, "32")
	if i32 != 32 {
		t.Errorf("int32 was %q", i)
	}

	var i64 int64
	AssignTo(&i64, "64")
	if i64 != 64 {
		t.Errorf("int32 was %q", i)
	}

}
