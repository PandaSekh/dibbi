package internal

import (
	"strings"
)

// Finds string Tokens delimited by delimiter
func lexCharacterDelimited(source string, initialCursor Cursor, delimiter byte) (*Token, Cursor, bool) {
	finalCursor := initialCursor

	if len(source[finalCursor.Pointer:]) == 0 {
		return nil, initialCursor, false
	}

	// Source must start with delimiter
	if source[finalCursor.Pointer] != delimiter {
		return nil, initialCursor, false
	}

	finalCursor.Location.Column++
	finalCursor.Pointer++

	var Value []byte
	for ; finalCursor.Pointer < uint(len(source)); finalCursor.Pointer++ {
		char := source[finalCursor.Pointer]

		if char == delimiter {
			if finalCursor.Pointer+1 >= uint(len(source)) || // if Cursor is at end of source
				source[finalCursor.Pointer+1] != delimiter { // or next char is not a delimiter
				// the full Value was read
				finalCursor.Pointer++
				finalCursor.Location.Column++

				return &Token{
					Value:    string(Value),
					Location: initialCursor.Location,
					Type:     StringType,
				}, finalCursor, true
			} else {
				// next char is a delimiter. In SQL a double delimiter is an escape
				Value = append(Value, delimiter)
				finalCursor.Pointer++
				finalCursor.Location.Column++
			}
		} else {
			Value = append(Value, char)
			finalCursor.Location.Column++
		}
	}

	return nil, initialCursor, false
}

// given a source, Cursor and an array of possible matches, returns the longest match found
func findLongestStringMatch(source string, initialCursor Cursor, options []string) string {
	var Value []byte
	var skipList []int
	var match string

	finalCursor := initialCursor

	for finalCursor.Pointer < uint(len(source)) {
		Value = append(Value, strings.ToLower(string(source[finalCursor.Pointer]))...)
		finalCursor.Pointer++

	match:
		for i, option := range options {
			for _, skip := range skipList {
				if i == skip {
					continue match
				}
			}

			if option == string(Value) {
				skipList = append(skipList, i)
				if len(option) > len(match) {
					match = option
				}

				continue
			}

			sharesPrefix := string(Value) == option[:finalCursor.Pointer-initialCursor.Pointer]
			tooLong := len(Value) > len(option)
			if tooLong || !sharesPrefix {
				skipList = append(skipList, i)
			}
		}
		if len(skipList) == len(options) {
			break
		}
	}

	return match
}
