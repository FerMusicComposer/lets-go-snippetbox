package assert

import (
	"strings"
	"testing"
)

// Equal is a function for comparing actual and expected values in Go testing.
//
// The function takes the testing.T object, actual and expected values as parameters.
// Parameters are of generic type, so actual and expected can be of any comparable type
// It does not return anything.
func Equal[T comparable](t *testing.T, actual, expected T) {
	// The t.Helper() function indicates to the Go
	// test runner that our Equal() function is a test helper. This means that when t.Errorf()
	// is called from our Equal() function, the Go test runner will report the filename and line
	// number of the code which called our Equal() function in the output.
	t.Helper()

	if actual != expected {
		t.Errorf("got: %v; want:%v;", actual, expected)
	}
}

// StringContains is a function to test if a string contains a specific substring.
//
// It takes parameters t *testing.T, actual string, and expectedSubstring string.
// It does not return anything.
func StringContains(t *testing.T, actual, expectedSubstring string) {
	t.Helper()
	if !strings.Contains(actual, expectedSubstring) {
		t.Errorf("got: %q; expected to contain: %q", actual, expectedSubstring)
	}
}

func NilError(t *testing.T, actual error) {
	t.Helper()
	if actual != nil {
		t.Errorf("got: %v; expected: nil", actual)
	}
}
