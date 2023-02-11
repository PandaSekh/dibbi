package dibbi

import (
	"bytes"
	"reflect"
	"testing"
)

func assertEqual(t *testing.T, expected interface{}, actual interface{}) {
	if !objectsAreEqual(expected, actual) {
		t.Fatalf("Objects are not equal. Expected: %s. Got: %s", expected, actual)
	}
}

func objectsAreEqual(expected, actual interface{}) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}

	exp, ok := expected.([]byte)
	if !ok {
		return reflect.DeepEqual(expected, actual)
	}

	act, ok := actual.([]byte)
	if !ok {
		return false
	}
	if exp == nil || act == nil {
		return exp == nil && act == nil
	}
	return bytes.Equal(exp, act)
}
