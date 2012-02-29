package document

import (
	"testing"
)

func TestFrac(t *testing.T) {
	testStr := "6  7/8"
	expectedNorm := "6 7/8"
	normalized, result, err := parseFrac(testStr)
	if err != nil {
		t.Errorf("could not parse %q", testStr)
	}
	if normalized != expectedNorm {
		t.Errorf("normalized %q = ", testStr)
	}
	if result != 6.875 {
		t.Errorf("parseFloat(%q) returned %g", testStr, result)
	}
}

func TestLengths(t *testing.T) {
	type data struct {
		Str		string
		Normstr	string
		Points	float64
		Ok		bool
	}
	var testData []data = []data { data{"1 cm", "1cm", 72.0 / 2.54, true},
		data{"1in", `1"`, 72.0, true},
		data{"6-7/8 pt", `6 7/8pt`, 6.875, true},
		data{"12", "", 0.0, false} }

	for _, d := range testData {
		length, err := LengthFromString(d.Str)
		if err != nil {
			if d.Ok {
				t.Errorf("length failed for %q\n", d.Str)
			}
		} else if !d.Ok {
			t.Errorf("length didn't fail for %q\n", d.Str)
		} else {
			if length.String() != d.Normstr {
				t.Errorf("length for %q has string %q\n", d.Str, length.String())
			}
			if length.Points() != d.Points {
				t.Errorf("length for %q has %g points\n", d.Str, length.Points())
			}
		}
	}
}

