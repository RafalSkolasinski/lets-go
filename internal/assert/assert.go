package assert

import "testing"

func Equal[T comparable](t *testing.T, actual, expected T) {
	// Indicate to the Go test runner that our Equal() function is a test helper.
	t.Helper()

	if actual != expected {
		t.Errorf("got: %v; want %v", actual, expected)
	}
}
