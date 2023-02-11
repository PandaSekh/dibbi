package dibbi

import (
	"testing"
)

func TestAssertEquals(t *testing.T) {
	tests := []struct {
		areEqual bool
		expected interface{}
		actual   interface{}
	}{
		{
			areEqual: true,
			expected: "105",
			actual:   "105",
		},
		{
			areEqual: true,
			expected: 105,
			actual:   105,
		},
		{
			areEqual: true,
			expected: []byte{1},
			actual:   []byte{1},
		},
		{
			areEqual: true,
			expected: "123.145",
			actual:   "123.145",
		},
		{
			areEqual: true,
			expected: true,
			actual:   true,
		},
		{
			areEqual: true,
			expected: []int{1, 2, 3},
			actual:   []int{1, 2, 3},
		},
		// false tests
		{
			areEqual: false,
			expected: "e4",
			actual:   "e5",
		},
		{
			areEqual: false,
			expected: true,
			actual:   false,
		},
		{
			areEqual: false,
			expected: "1ee4",
			actual:   false,
		},
		{
			areEqual: false,
			expected: 1,
			actual:   float32(1),
		},
	}

	for _, test := range tests {
		areEqual := objectsAreEqual(test.expected, test.actual)
		assertEqual(t, test.areEqual, areEqual)
	}
}
