package dibbi

import (
	"fmt"
	"strings"
)

type location struct {
	line   uint
	column uint
}
type keyword string
type symbol string

// token represents a lexed object from an input
type token struct {
	value     string
	tokenType tokenType
	location  location
}
type tokenType uint

func (t *token) equals(other *token) bool {
	return t.value == other.value && t.tokenType == other.tokenType
}

// getBindingPower returns the bp of a given token.
// Reference: https://matklad.github.io/2020/04/13/simple-but-powerful-pratt-parsing.html
func (t *token) getBindingPower() uint {
	switch t.tokenType {
	case keywordType:
		switch keyword(t.value) {
		case andKeyword:
			fallthrough
		case orKeyword:
			return 1
		}
	case symbolType:
		switch symbol(t.value) {
		case equalsSymbol:
			fallthrough
		case notEqualSymbol:
			return 2

		case lessThanSymbol:
			fallthrough
		case greaterThanSymbol:
			return 3

		case lessThanEqualSymbol:
			fallthrough
		case greaterThanEqualSymbol:
			return 4

		case concatSymbol:
			fallthrough
		case plusSymbol:
			return 5
		}
	}

	return 0
}

type Cursor struct {
	Pointer  uint
	location location
}
type lexer func(string, Cursor) (*token, Cursor, bool)

const (
	SelectKeyword     keyword = "select"
	FromKeyword       keyword = "from"
	AsKeyword         keyword = "as"
	TableKeyword      keyword = "table"
	CreateKeyword     keyword = "create"
	InsertKeyword     keyword = "insert"
	IntoKeyword       keyword = "into"
	ValuesKeyword     keyword = "values"
	IntKeyword        keyword = "int"
	TextKeyword       keyword = "text"
	BoolKeyword       keyword = "bool"
	WhereKeyword      keyword = "where"
	andKeyword        keyword = "and"
	orKeyword         keyword = "or"
	TrueKeyword       keyword = "true"
	FalseKeyword      keyword = "false"
	UniqueKeyword     keyword = "unique"
	IndexKeyword      keyword = "index"
	OnKeyword         keyword = "on"
	PrimaryKeyKeyword keyword = "primary key"
	NullKeyword       keyword = "null"
	LimitKeyword      keyword = "limit"
	OffsetKeyword     keyword = "offset"
	DropKeyword       keyword = "drop"

	SemicolonSymbol        symbol = ";"
	AsteriskSymbol         symbol = "*"
	CommaSymbol            symbol = ","
	leftParenthesesSymbol  symbol = "("
	rightParenthesesSymbol symbol = ")"
	equalsSymbol           symbol = "="
	notEqualSymbol         symbol = "<>"
	NeqSymbol2             symbol = "!="
	concatSymbol           symbol = "||"
	plusSymbol             symbol = "+"
	lessThanSymbol         symbol = "<"
	lessThanEqualSymbol    symbol = "<="
	greaterThanSymbol      symbol = ">"
	greaterThanEqualSymbol symbol = ">="

	keywordType tokenType = iota
	symbolType
	identifierType
	stringType
	NumericType
	booleanType
	NullType
)

func (t *token) Equals(other *token) bool {
	return t.value == other.value && t.tokenType == other.tokenType
}

func Lex(source string) ([]*token, error) {
	var tokens []*token
	cur := Cursor{}
	lexers := []lexer{lexKeyword, lexSymbol, lexString, lexNumeric, lexIdentifier, lexBool}

lex:
	for cur.Pointer < uint(len(source)) {
		// Try to parse token with lexers defined above
		for _, lex := range lexers {
			if token, newCursor, success := lex(source, cur); success {
				// One lexer was successful, update Cursor and append token
				cur = newCursor
				if token != nil {
					tokens = append(tokens, token)
				}

				continue lex
			}
		}

		// No lexer was able to perform lexing. Return error
		hint := ""
		if len(tokens) > 0 {
			hint = " after " + tokens[len(tokens)-1].value
		}
		return nil, fmt.Errorf("unable to lex token%s, at %d:%d", hint, cur.location.line, cur.location.column)
	}

	return tokens, nil
}

////////////////////////////////
// Lex Functions
////////////////////////////////

// Finds numeric tokens
func lexNumeric(source string, initialCursor Cursor) (*token, Cursor, bool) {
	finalCursor := initialCursor

	periodAlreadyFound := false
	exponentialAlreadyFound := false

	for ; finalCursor.Pointer < uint(len(source)); finalCursor.Pointer++ {
		char := source[finalCursor.Pointer]
		finalCursor.location.column++

		isDigit := char >= '0' && char <= '9'
		isPeriod := char == '.'
		isExponentialMarker := char == 'e'

		// Validate start of expression (should be a digit)
		if finalCursor.Pointer == initialCursor.Pointer {
			if !isDigit && !isPeriod {
				return nil, initialCursor, false
			}

			periodAlreadyFound = isPeriod
		} else if isPeriod {
			if periodAlreadyFound {
				// period was already found
				return nil, initialCursor, false
			}

			periodAlreadyFound = true
		} else if isExponentialMarker {
			if exponentialAlreadyFound {
				return nil, initialCursor, false
			}

			// No periods can be present after exponential
			exponentialAlreadyFound = true
			periodAlreadyFound = true

			// Exponential must be followed by digits
			if finalCursor.Pointer == uint(len(source)-1) {
				return nil, initialCursor, false
			}

			nextChar := source[finalCursor.Pointer+1]
			if nextChar == '-' || nextChar == '+' {
				finalCursor.Pointer++
				finalCursor.location.column++
			}
		} else if !isDigit {
			break
		}
	}

	// No characters accumulated
	if finalCursor.Pointer == initialCursor.Pointer {
		return nil, initialCursor, false
	}

	return &token{
		value:     source[initialCursor.Pointer:finalCursor.Pointer],
		location:  initialCursor.location,
		tokenType: NumericType,
	}, finalCursor, true
}

// Lex a string delimited by '
func lexString(source string, initialCursor Cursor) (*token, Cursor, bool) {
	return lexCharacterDelimited(source, initialCursor, '\'')
}

func lexSymbol(source string, initialCursor Cursor) (*token, Cursor, bool) {
	char := source[initialCursor.Pointer]
	finalCursor := initialCursor

	finalCursor.Pointer++
	finalCursor.location.column++

	// symbols to be ignored
	switch char {
	case '\n':
		finalCursor.location.line++
		finalCursor.location.column = 0
		fallthrough
	case '\t':
		fallthrough
	case ' ':
		return nil, finalCursor, true
	}

	symbols := []symbol{
		CommaSymbol,
		leftParenthesesSymbol,
		rightParenthesesSymbol,
		SemicolonSymbol,
		AsteriskSymbol,
		equalsSymbol,
		notEqualSymbol,
		NeqSymbol2,
		lessThanSymbol,
		lessThanEqualSymbol,
		greaterThanSymbol,
		greaterThanEqualSymbol,
		concatSymbol,
		plusSymbol,
	}
	var options []string

	for _, symbol := range symbols {
		options = append(options, string(symbol))
	}

	match := findLongestStringMatch(source, initialCursor, options)
	if match == "" {
		return nil, initialCursor, false
	}

	finalCursor.Pointer = initialCursor.Pointer + uint(len(match))
	finalCursor.location.column = initialCursor.location.column + uint(len(match))

	// != is rewritten as <>
	if match == string(NeqSymbol2) {
		match = string(notEqualSymbol)
	}

	return &token{
		value:     match,
		location:  initialCursor.location,
		tokenType: symbolType,
	}, finalCursor, true
}

func lexKeyword(source string, initialCursor Cursor) (*token, Cursor, bool) {
	finalCursor := initialCursor
	keywords := []keyword{
		SelectKeyword,
		InsertKeyword,
		ValuesKeyword,
		TableKeyword,
		CreateKeyword,
		DropKeyword,
		WhereKeyword,
		FromKeyword,
		IntoKeyword,
		TextKeyword,
		BoolKeyword,
		IntKeyword,
		andKeyword,
		orKeyword,
		AsKeyword,
		//TrueKeyword,
		//FalseKeyword,
		UniqueKeyword,
		IndexKeyword,
		OnKeyword,
		PrimaryKeyKeyword,
		NullKeyword,
		LimitKeyword,
		OffsetKeyword,
	}

	var options []string
	for _, keyword := range keywords {
		options = append(options, string(keyword))
	}

	match := findLongestStringMatch(source, initialCursor, options)
	if match == "" {
		return nil, initialCursor, false
	}

	finalCursor.Pointer = initialCursor.Pointer + uint(len(match))
	finalCursor.location.column = initialCursor.location.column + uint(len(match))

	tokenType := keywordType
	if match == string(NullKeyword) {
		tokenType = NullType
	}

	return &token{
		value:     match,
		location:  initialCursor.location,
		tokenType: tokenType,
	}, finalCursor, true
}

func lexIdentifier(source string, initialCursor Cursor) (*token, Cursor, bool) {
	// Try to lex with helper function if it's a delimited identifier
	if token, newCursor, ok := lexCharacterDelimited(source, initialCursor, '"'); ok {
		return token, newCursor, true
	}

	finalCursor := initialCursor

	char := source[finalCursor.Pointer]
	// ASCII only
	isAlphabetical := (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z')
	if !isAlphabetical {
		return nil, initialCursor, false
	}

	finalCursor.Pointer++
	finalCursor.location.column++

	value := []byte{char}
	for ; finalCursor.Pointer < uint(len(source)); finalCursor.Pointer++ {
		char = source[finalCursor.Pointer]
		isAlphabetical := (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z')
		isNumeric := char >= '0' && char <= '9'
		if isAlphabetical || isNumeric || char == '$' || char == '_' {
			value = append(value, char)
			finalCursor.location.column++
			continue
		}

		break
	}

	if len(value) == 0 {
		return nil, initialCursor, false
	}

	return &token{
		value:     strings.ToLower(string(value)),
		location:  initialCursor.location,
		tokenType: identifierType,
	}, finalCursor, true
}

func lexBool(source string, initialCursor Cursor) (*token, Cursor, bool) {
	finalCursor := initialCursor

	for ; finalCursor.Pointer < uint(len(source)); finalCursor.Pointer++ {
		char := source[finalCursor.Pointer]
		finalCursor.location.column++

		// todo refactor
		isLetterOfInterest := char == 't' || char == 'r' || char == 'u' || char == 'e' ||
			char == 'f' || char == 'a' || char == 'l' || char == 's'

		if !isLetterOfInterest {
			break
		}
	}

	// No characters accumulated
	if finalCursor.Pointer == initialCursor.Pointer {
		return nil, initialCursor, false
	}

	// Word found is not a boolean
	if source[initialCursor.Pointer:finalCursor.Pointer] != string(TrueKeyword) &&
		source[initialCursor.Pointer:finalCursor.Pointer] != string(FalseKeyword) {
		return nil, initialCursor, false
	}

	return &token{
		value:     source[initialCursor.Pointer:finalCursor.Pointer],
		location:  initialCursor.location,
		tokenType: booleanType,
	}, finalCursor, true
}
