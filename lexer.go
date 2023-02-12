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
type tokenType uint
type token struct {
	value     string
	tokenType tokenType
	location  location
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
	AndKeyword        keyword = "and"
	OrKeyword         keyword = "or"
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
	LeftParenthesesSymbol  symbol = "("
	RightParenthesesSymbol symbol = ")"
	EqualsSymbol           symbol = "="
	NeqSymbol              symbol = "<>"
	NeqSymbol2             symbol = "!="
	ConcatSymbol           symbol = "||"
	PlusSymbol             symbol = "+"
	LtSymbol               symbol = "<"
	LteSymbol              symbol = "<="
	GtSymbol               symbol = ">"
	GteSymbol              symbol = ">="

	KeywordType tokenType = iota
	SymbolType
	IdentifierType
	StringType
	NumericType
	BooleanType
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
		LeftParenthesesSymbol,
		RightParenthesesSymbol,
		SemicolonSymbol,
		AsteriskSymbol,
		EqualsSymbol,
		NeqSymbol,
		NeqSymbol2,
		LtSymbol,
		LteSymbol,
		GtSymbol,
		GteSymbol,
		ConcatSymbol,
		PlusSymbol,
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
		match = string(NeqSymbol)
	}

	return &token{
		value:     match,
		location:  initialCursor.location,
		tokenType: SymbolType,
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
		AndKeyword,
		OrKeyword,
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

	tokenType := KeywordType
	//if match == string(TrueKeyword) || match == string(FalseKeyword) {
	//	tokenType = BooleanType
	//}

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
		tokenType: IdentifierType,
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
		tokenType: BooleanType,
	}, finalCursor, true
}
