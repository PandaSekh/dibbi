package internal

import (
	"fmt"
	"strings"
)

type Location struct {
	Line   uint
	Column uint
}
type Keyword string
type Symbol string
type TokenType uint
type Token struct {
	Value    string
	Type     TokenType
	Location Location
}
type Cursor struct {
	Pointer  uint
	Location Location
}
type lexer func(string, Cursor) (*Token, Cursor, bool)

const (
	SelectKeyword Keyword = "select"
	FromKeyword   Keyword = "from"
	AsKeyword     Keyword = "as"
	TableKeyword  Keyword = "table"
	CreateKeyword Keyword = "create"
	InsertKeyword Keyword = "insert"
	IntoKeyword   Keyword = "into"
	ValuesKeyword Keyword = "values"
	IntKeyword    Keyword = "int"
	TextKeyword   Keyword = "text"
	BoolKeyword   Keyword = "bool"
	WhereKeyword  Keyword = "where"

	SemicolonSymbol        Symbol = ";"
	AsteriskSymbol         Symbol = "*"
	CommaSymbol            Symbol = ","
	LeftParenthesesSymbol  Symbol = "("
	RightParenthesesSymbol Symbol = ")"
	EqualsSymbol           Symbol = "="

	KeywordType TokenType = iota
	SymbolType
	IdentifierType
	StringType
	NumericType
	BooleanType
	NullType
)

func (t *Token) Equals(other *Token) bool {
	return t.Value == other.Value && t.Type == other.Type
}

func Lex(source string) ([]*Token, error) {
	var tokens []*Token
	cur := Cursor{}
	lexers := []lexer{lexKeyword, lexSymbol, lexString, lexNumeric, lexIdentifier, lexBool}

lex:
	for cur.Pointer < uint(len(source)) {
		// Try to parse Token with lexers defined above
		for _, lex := range lexers {
			if Token, newCursor, success := lex(source, cur); success {
				// One lexer was successful, update Cursor and append Token
				cur = newCursor
				if Token != nil {
					tokens = append(tokens, Token)
				}

				continue lex
			}
		}

		// No lexer was able to perform lexing. Return error
		hint := ""
		if len(tokens) > 0 {
			hint = " after " + tokens[len(tokens)-1].Value
		}
		return nil, fmt.Errorf("unable to lex Token%s, at %d:%d", hint, cur.Location.Line, cur.Location.Column)
	}

	return tokens, nil
}

////////////////////////////////
// Lex Functions
////////////////////////////////

// Finds numeric tokens
func lexNumeric(source string, initialCursor Cursor) (*Token, Cursor, bool) {
	finalCursor := initialCursor

	periodAlreadyFound := false
	exponentialAlreadyFound := false

	for ; finalCursor.Pointer < uint(len(source)); finalCursor.Pointer++ {
		char := source[finalCursor.Pointer]
		finalCursor.Location.Column++

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
				finalCursor.Location.Column++
			}
		} else if !isDigit {
			break
		}
	}

	// No characters accumulated
	if finalCursor.Pointer == initialCursor.Pointer {
		return nil, initialCursor, false
	}

	return &Token{
		Value:    source[initialCursor.Pointer:finalCursor.Pointer],
		Location: initialCursor.Location,
		Type:     NumericType,
	}, finalCursor, true
}

// Lex a string delimited by '
func lexString(source string, initialCursor Cursor) (*Token, Cursor, bool) {
	return lexCharacterDelimited(source, initialCursor, '\'')
}

// Lex a boolean
func lexBool(source string, initialCursor Cursor) (*Token, Cursor, bool) {
	finalCursor := initialCursor

	for ; finalCursor.Pointer < uint(len(source)); finalCursor.Pointer++ {
		char := source[finalCursor.Pointer]
		finalCursor.Location.Column++

		// todo refactor
		isLetterOfInterest := char == 't' || char == 'r' || char == 'u' || char == 'e' || char == 'f' || char == 'a' || char == 'l' || char == 's'

		if !isLetterOfInterest {
			break
		}
	}

	// No characters accumulated
	if finalCursor.Pointer == initialCursor.Pointer {
		return nil, initialCursor, false
	}

	// Word found is not a boolean
	if source[initialCursor.Pointer:finalCursor.Pointer] != "true" && source[initialCursor.Pointer:finalCursor.Pointer] != "false" {
		return nil, initialCursor, false
	}

	return &Token{
		Value:    source[initialCursor.Pointer:finalCursor.Pointer],
		Location: initialCursor.Location,
		Type:     BooleanType,
	}, finalCursor, true
}

func lexSymbol(source string, initialCursor Cursor) (*Token, Cursor, bool) {
	char := source[initialCursor.Pointer]
	finalCursor := initialCursor

	finalCursor.Pointer++
	finalCursor.Location.Column++

	// symbols to be ignored
	switch char {
	case '\n':
		finalCursor.Location.Line++
		finalCursor.Location.Column = 0
		fallthrough
	case '\t':
		fallthrough
	case ' ':
		return nil, finalCursor, true
	}

	symbols := []Symbol{CommaSymbol, LeftParenthesesSymbol, RightParenthesesSymbol, SemicolonSymbol, AsteriskSymbol, EqualsSymbol}
	var options []string

	for _, Symbol := range symbols {
		options = append(options, string(Symbol))
	}

	match := findLongestStringMatch(source, initialCursor, options)
	if match == "" {
		return nil, initialCursor, false
	}

	finalCursor.Pointer = initialCursor.Pointer + uint(len(match))
	finalCursor.Location.Column = initialCursor.Location.Column + uint(len(match))

	return &Token{
		Value:    match,
		Location: initialCursor.Location,
		Type:     SymbolType,
	}, finalCursor, true
}

func lexKeyword(source string, initialCursor Cursor) (*Token, Cursor, bool) {
	finalCursor := initialCursor
	keywords := []Keyword{
		SelectKeyword,
		InsertKeyword,
		ValuesKeyword,
		TableKeyword,
		CreateKeyword,
		WhereKeyword,
		FromKeyword,
		IntoKeyword,
		TextKeyword,
		IntKeyword,
		BoolKeyword,
		AsKeyword,
	}

	var options []string
	for _, Keyword := range keywords {
		options = append(options, string(Keyword))
	}

	match := findLongestStringMatch(source, initialCursor, options)
	if match == "" {
		return nil, initialCursor, false
	}

	finalCursor.Pointer = initialCursor.Pointer + uint(len(match))
	finalCursor.Location.Column = initialCursor.Location.Column + uint(len(match))

	return &Token{
		Value:    match,
		Location: initialCursor.Location,
		Type:     KeywordType,
	}, finalCursor, true
}

func lexIdentifier(source string, initialCursor Cursor) (*Token, Cursor, bool) {
	// Try to lex with helper function if it's a delimited identifier
	if Token, newCursor, ok := lexCharacterDelimited(source, initialCursor, '"'); ok {
		return Token, newCursor, true
	}

	finalCursor := initialCursor

	char := source[finalCursor.Pointer]
	// ASCII only
	isAlphabetical := (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z')
	if !isAlphabetical {
		return nil, initialCursor, false
	}

	finalCursor.Pointer++
	finalCursor.Location.Column++

	Value := []byte{char}
	for ; finalCursor.Pointer < uint(len(source)); finalCursor.Pointer++ {
		char = source[finalCursor.Pointer]
		isAlphabetical := (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z')
		isNumeric := char >= '0' && char <= '9'
		if isAlphabetical || isNumeric || char == '$' || char == '_' {
			Value = append(Value, char)
			finalCursor.Location.Column++
			continue
		}

		break
	}

	if len(Value) == 0 {
		return nil, initialCursor, false
	}

	return &Token{
		Value:    strings.ToLower(string(Value)),
		Location: initialCursor.Location,
		Type:     IdentifierType,
	}, finalCursor, true
}
