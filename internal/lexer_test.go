package internal

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToken_lexFloat(t *testing.T) {
	tests := []struct {
		number bool
		value  string
	}{
		{
			number: true,
			value:  "105",
		},
		{
			number: true,
			value:  "105 ",
		},
		{
			number: true,
			value:  "123.",
		},
		{
			number: true,
			value:  "123.145",
		},
		{
			number: true,
			value:  "1e5",
		},
		{
			number: true,
			value:  "1.e21",
		},
		{
			number: true,
			value:  "1.1e2",
		},
		{
			number: true,
			value:  "1.1e-2",
		},
		{
			number: true,
			value:  "1.1e+2",
		},
		{
			number: true,
			value:  "1e-1",
		},
		{
			number: true,
			value:  ".1",
		},
		{
			number: true,
			value:  "4.",
		},
		// false tests
		{
			number: false,
			value:  "e4",
		},
		{
			number: false,
			value:  "1..",
		},
		{
			number: false,
			value:  "1ee4",
		},
		{
			number: false,
			value:  " 1",
		},
	}

	for _, test := range tests {
		tok, _, ok := lexNumeric(test.value, Cursor{})
		assert.Equal(t, test.number, ok, test.value)
		if ok {
			assert.Equal(t, strings.TrimSpace(test.value), tok.Value, test.value)
		}
	}
}

func TestToken_lexString(t *testing.T) {
	tests := []struct {
		string   bool
		value    string
		expected string
	}{
		{
			string:   false,
			value:    "a",
			expected: "a",
		},
		{
			string:   true,
			value:    "'abc'",
			expected: "'abc'",
		},
		{
			string:   true,
			value:    "'a b'",
			expected: "'a b'",
		},
		{
			string:   true,
			value:    "'a' ",
			expected: "'a' ",
		},
		{
			string:   true,
			value:    "'a '' b'",
			expected: "'a ' b'",
		},
		// false tests
		{
			string:   false,
			value:    "'",
			expected: "'",
		},
		{
			string:   false,
			value:    "",
			expected: "",
		},
		{
			string:   false,
			value:    " 'foo'",
			expected: " 'foo'",
		},
	}

	for _, test := range tests {
		tok, _, ok := lexString(test.value, Cursor{})
		assert.Equal(t, test.string, ok, test.value)
		if ok {
			test.expected = strings.TrimSpace(test.expected)
			assert.Equal(t, test.expected[1:len(test.expected)-1], tok.Value, test.expected)
		}
	}
}

func TestToken_lexSymbol(t *testing.T) {
	tests := []struct {
		symbol bool
		value  string
	}{
		{
			symbol: true,
			value:  "= ",
		},
	}

	for _, test := range tests {
		tok, _, ok := lexSymbol(test.value, Cursor{})
		assert.Equal(t, test.symbol, ok, test.value)
		if ok {
			test.value = strings.TrimSpace(test.value)
			assert.Equal(t, test.value, tok.Value, test.value)
		}
	}
}

func TestToken_lexIdentifier(t *testing.T) {
	tests := []struct {
		Identifier bool
		input      string
		value      string
	}{
		{
			Identifier: true,
			input:      "a",
			value:      "a",
		},
		{
			Identifier: true,
			input:      "abc",
			value:      "abc",
		},
		{
			Identifier: true,
			input:      "abc ",
			value:      "abc",
		},
		{
			Identifier: true,
			input:      `" abc "`,
			value:      ` abc `,
		},
		{
			Identifier: true,
			input:      "a9$",
			value:      "a9$",
		},
		{
			Identifier: true,
			input:      "userName",
			value:      "username",
		},
		{
			Identifier: true,
			input:      `"userName"`,
			value:      "userName",
		},
		// false tests
		{
			Identifier: false,
			input:      `"`,
		},
		{
			Identifier: false,
			input:      "_sadsfa",
		},
		{
			Identifier: false,
			input:      "9sadsfa",
		},
		{
			Identifier: false,
			input:      " abc",
		},
	}

	for _, test := range tests {
		tok, _, ok := lexIdentifier(test.input, Cursor{})
		assert.Equal(t, test.Identifier, ok, test.input)
		if ok {
			assert.Equal(t, test.value, tok.Value, test.input)
		}
	}
}

func TestToken_lexKeyword(t *testing.T) {
	tests := []struct {
		keyword bool
		value   string
	}{
		{
			keyword: true,
			value:   "select ",
		},
		{
			keyword: true,
			value:   "from",
		},
		{
			keyword: true,
			value:   "as",
		},
		{
			keyword: true,
			value:   "SELECT",
		},
		{
			keyword: true,
			value:   "into",
		},
		// false tests
		{
			keyword: false,
			value:   " into",
		},
		{
			keyword: false,
			value:   "flubbrety",
		},
	}

	for _, test := range tests {
		tok, _, ok := lexKeyword(test.value, Cursor{})
		assert.Equal(t, test.keyword, ok, test.value)
		if ok {
			test.value = strings.TrimSpace(test.value)
			assert.Equal(t, strings.ToLower(test.value), tok.Value, test.value)
		}
	}
}

func TestLex(t *testing.T) {
	tests := []struct {
		input  string
		Tokens []Token
		err    error
	}{
		{
			input: "select a",
			Tokens: []Token{
				{
					Location: Location{Column: 0, Line: 0},
					Value:    string(SelectKeyword),
					Type:     KeywordType,
				},
				{
					Location: Location{Column: 7, Line: 0},
					Value:    "a",
					Type:     IdentifierType,
				},
			},
		},
		{
			input: "select 1",
			Tokens: []Token{
				{
					Location: Location{Column: 0, Line: 0},
					Value:    string(SelectKeyword),
					Type:     KeywordType,
				},
				{
					Location: Location{Column: 7, Line: 0},
					Value:    "1",
					Type:     NumericType,
				},
			},
			err: nil,
		},
		{
			input: "CREATE TABLE u (id INT, name TEXT)",
			Tokens: []Token{
				{
					Location: Location{Column: 0, Line: 0},
					Value:    string(CreateKeyword),
					Type:     KeywordType,
				},
				{
					Location: Location{Column: 7, Line: 0},
					Value:    string(TableKeyword),
					Type:     KeywordType,
				},
				{
					Location: Location{Column: 13, Line: 0},
					Value:    "u",
					Type:     IdentifierType,
				},
				{
					Location: Location{Column: 15, Line: 0},
					Value:    "(",
					Type:     SymbolType,
				},
				{
					Location: Location{Column: 16, Line: 0},
					Value:    "id",
					Type:     IdentifierType,
				},
				{
					Location: Location{Column: 19, Line: 0},
					Value:    "int",
					Type:     KeywordType,
				},
				{
					Location: Location{Column: 22, Line: 0},
					Value:    ",",
					Type:     SymbolType,
				},
				{
					Location: Location{Column: 24, Line: 0},
					Value:    "name",
					Type:     IdentifierType,
				},
				{
					Location: Location{Column: 29, Line: 0},
					Value:    "text",
					Type:     KeywordType,
				},
				{
					Location: Location{Column: 33, Line: 0},
					Value:    ")",
					Type:     SymbolType,
				},
			},
		},
		{
			input: "insert into users Values (105, 233)",
			Tokens: []Token{
				{
					Location: Location{Column: 0, Line: 0},
					Value:    string(InsertKeyword),
					Type:     KeywordType,
				},
				{
					Location: Location{Column: 7, Line: 0},
					Value:    string(IntoKeyword),
					Type:     KeywordType,
				},
				{
					Location: Location{Column: 12, Line: 0},
					Value:    "users",
					Type:     IdentifierType,
				},
				{
					Location: Location{Column: 18, Line: 0},
					Value:    string(ValuesKeyword),
					Type:     KeywordType,
				},
				{
					Location: Location{Column: 25, Line: 0},
					Value:    "(",
					Type:     SymbolType,
				},
				{
					Location: Location{Column: 26, Line: 0},
					Value:    "105",
					Type:     NumericType,
				},
				{
					Location: Location{Column: 30, Line: 0},
					Value:    ",",
					Type:     SymbolType,
				},
				{
					Location: Location{Column: 32, Line: 0},
					Value:    "233",
					Type:     NumericType,
				},
				{
					Location: Location{Column: 36, Line: 0},
					Value:    ")",
					Type:     SymbolType,
				},
			},
			err: nil,
		},
		{
			input: "SELECT id FROM users;",
			Tokens: []Token{
				{
					Location: Location{Column: 0, Line: 0},
					Value:    string(SelectKeyword),
					Type:     KeywordType,
				},
				{
					Location: Location{Column: 7, Line: 0},
					Value:    "id",
					Type:     IdentifierType,
				},
				{
					Location: Location{Column: 10, Line: 0},
					Value:    string(FromKeyword),
					Type:     KeywordType,
				},
				{
					Location: Location{Column: 15, Line: 0},
					Value:    "users",
					Type:     IdentifierType,
				},
				{
					Location: Location{Column: 20, Line: 0},
					Value:    ";",
					Type:     SymbolType,
				},
			},
			err: nil,
		},
		{
			input: "SELECT * FROM my_table WHERE name = 'hello_world'",
			Tokens: []Token{
				{
					Location: Location{Column: 0, Line: 0},
					Value:    string(SelectKeyword),
					Type:     KeywordType,
				},
				{
					Location: Location{Column: 7, Line: 0},
					Value:    "*",
					Type:     SymbolType,
				},
				{
					Location: Location{Column: 9, Line: 0},
					Value:    string(FromKeyword),
					Type:     KeywordType,
				},
				{
					Location: Location{Column: 14, Line: 0},
					Value:    "my_table",
					Type:     IdentifierType,
				},
				{
					Location: Location{Column: 23, Line: 0},
					Value:    "where",
					Type:     KeywordType,
				},
				{
					Location: Location{Column: 29, Line: 0},
					Value:    "name",
					Type:     IdentifierType,
				},
				{
					Location: Location{Column: 34, Line: 0},
					Value:    "=",
					Type:     SymbolType,
				},
				{
					Location: Location{Column: 36, Line: 0},
					Value:    "hello_world",
					Type:     StringType,
				},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		tokens, err := Lex(test.input)
		assert.Equal(t, test.err, err, test.input)
		assert.Equal(t, len(test.Tokens), len(tokens), test.input)

		for i, tok := range tokens {
			assert.Equal(t, &test.Tokens[i], tok, test.input)
		}
	}
}

// /////////////////////////
// / Benchmarks
// /////////////////////////
var inputTable = []struct {
	input string
}{
	{input: "SELECT * FROM my_table WHERE name = 'hello_world'"},
	{input: "SELECT (id, name, surname, city, address, nation) FROM my_long_table_name as mltb WHERE name = 'hello_world' AND surname = 'longer_word'"},
}

func BenchmarkLex(b *testing.B) {
	for _, v := range inputTable {
		b.Run(fmt.Sprintf("input_length_%d", len(v.input)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = Lex(v.input)
			}
		})
	}
}
