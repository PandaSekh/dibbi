package dibbi

import (
	"strings"
)

// Finds string Tokens delimited by delimiter
func lexCharacterDelimited(source string, initialCursor Cursor, delimiter byte) (*token, Cursor, bool) {
	finalCursor := initialCursor

	if len(source[finalCursor.Pointer:]) == 0 {
		return nil, initialCursor, false
	}

	// Source must start with delimiter
	if source[finalCursor.Pointer] != delimiter {
		return nil, initialCursor, false
	}

	finalCursor.location.column++
	finalCursor.Pointer++

	var value []byte
	for ; finalCursor.Pointer < uint(len(source)); finalCursor.Pointer++ {
		char := source[finalCursor.Pointer]

		if char == delimiter {
			if finalCursor.Pointer+1 >= uint(len(source)) || // if Cursor is at end of source
				source[finalCursor.Pointer+1] != delimiter { // or next char is not a delimiter
				// the full value was read
				finalCursor.Pointer++
				finalCursor.location.column++

				return &token{
					value:     string(value),
					location:  initialCursor.location,
					tokenType: stringType,
				}, finalCursor, true
			} else {
				// next char is a delimiter. In SQL a double delimiter is an escape
				value = append(value, delimiter)
				finalCursor.Pointer++
				finalCursor.location.column++
			}
		} else {
			value = append(value, char)
			finalCursor.location.column++
		}
	}

	return nil, initialCursor, false
}

// given a source, Cursor and an array of possible matches, returns the longest match found
func findLongestStringMatch(source string, initialCursor Cursor, options []string) string {
	var value []byte
	var skipList []int
	var match string

	finalCursor := initialCursor

	for finalCursor.Pointer < uint(len(source)) {
		value = append(value, strings.ToLower(string(source[finalCursor.Pointer]))...)
		finalCursor.Pointer++

	match:
		for i, option := range options {
			for _, skip := range skipList {
				if i == skip {
					continue match
				}
			}

			if option == string(value) {
				skipList = append(skipList, i)
				if len(option) > len(match) {
					match = option
				}

				continue
			}

			sharesPrefix := string(value) == option[:finalCursor.Pointer-initialCursor.Pointer]
			tooLong := len(value) > len(option)
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
