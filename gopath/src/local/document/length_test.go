package document

import (
	"testing"
)

func TestFrac(t *testing.T) {
	testStr := "6 7/8"
	result := parseFrac(testStr)
	if result != 6.875 {
		t.Errorf("parseFloat(%q) returned %g", testStr, result)
	}
}
