package dibbi

import (
	"fmt"
	"strings"
	"testing"
)

func TestToken_lexFloat(t *testing.T) {
	tests := []struct {
		isInt bool
		value string
	}{
		{
			isInt: true,
			value: "105",
		},
		{
			isInt: true,
			value: "105 ",
		},
		{
			isInt: true,
			value: "123.",
		},
		{
			isInt: true,
			value: "123.145",
		},
		{
			isInt: true,
			value: "1e5",
		},
		{
			isInt: true,
			value: "1.e21",
		},
		{
			isInt: true,
			value: "1.1e2",
		},
		{
			isInt: true,
			value: "1.1e-2",
		},
		{
			isInt: true,
			value: "1.1e+2",
		},
		{
			isInt: true,
			value: "1e-1",
		},
		{
			isInt: true,
			value: ".1",
		},
		{
			isInt: true,
			value: "4.",
		},
		// false tests
		{
			isInt: false,
			value: "e4",
		},
		{
			isInt: false,
			value: "1..",
		},
		{
			isInt: false,
			value: "1ee4",
		},
		{
			isInt: false,
			value: " 1",
		},
	}

	for _, test := range tests {
		tok, _, ok := lexNumeric(test.value, Cursor{})
		assertEqual(t, test.isInt, ok)
		if ok {
			assertEqual(t, strings.TrimSpace(test.value), tok.value)
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
			value:    "left",
			expected: "left",
		},
		{
			string:   true,
			value:    "'abc'",
			expected: "'abc'",
		},
		{
			string:   true,
			value:    "'left right'",
			expected: "'left right'",
		},
		{
			string:   true,
			value:    "'left' ",
			expected: "'left' ",
		},
		{
			string:   true,
			value:    "'left '' right'",
			expected: "'left ' right'",
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
		assertEqual(t, test.string, ok)
		if ok {
			test.expected = strings.TrimSpace(test.expected)
			assertEqual(t, test.expected[1:len(test.expected)-1], tok.value)
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
		assertEqual(t, test.symbol, ok)
		if ok {
			test.value = strings.TrimSpace(test.value)
			assertEqual(t, test.value, tok.value)
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
			input:      "left",
			value:      "left",
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
		assertEqual(t, test.Identifier, ok)
		if ok {
			assertEqual(t, test.value, tok.value)
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
		assertEqual(t, test.keyword, ok)
		if ok {
			test.value = strings.TrimSpace(test.value)
			assertEqual(t, strings.ToLower(test.value), tok.value)
		}
	}
}

func TestLex(t *testing.T) {
	tests := []struct {
		input  string
		Tokens []token
		err    error
	}{
		{
			input: "select left",
			Tokens: []token{
				{
					location:  location{column: 0, line: 0},
					value:     string(SelectKeyword),
					tokenType: keywordType,
				},
				{
					location:  location{column: 7, line: 0},
					value:     "left",
					tokenType: IdentifierType,
				},
			},
		},
		{
			input: "select 1",
			Tokens: []token{
				{
					location:  location{column: 0, line: 0},
					value:     string(SelectKeyword),
					tokenType: keywordType,
				},
				{
					location:  location{column: 7, line: 0},
					value:     "1",
					tokenType: NumericType,
				},
			},
			err: nil,
		},
		{
			input: "CREATE TABLE u (id INT, name TEXT)",
			Tokens: []token{
				{
					location:  location{column: 0, line: 0},
					value:     string(CreateKeyword),
					tokenType: keywordType,
				},
				{
					location:  location{column: 7, line: 0},
					value:     string(TableKeyword),
					tokenType: keywordType,
				},
				{
					location:  location{column: 13, line: 0},
					value:     "u",
					tokenType: IdentifierType,
				},
				{
					location:  location{column: 15, line: 0},
					value:     "(",
					tokenType: symbolType,
				},
				{
					location:  location{column: 16, line: 0},
					value:     "id",
					tokenType: IdentifierType,
				},
				{
					location:  location{column: 19, line: 0},
					value:     "int",
					tokenType: keywordType,
				},
				{
					location:  location{column: 22, line: 0},
					value:     ",",
					tokenType: symbolType,
				},
				{
					location:  location{column: 24, line: 0},
					value:     "name",
					tokenType: IdentifierType,
				},
				{
					location:  location{column: 29, line: 0},
					value:     "text",
					tokenType: keywordType,
				},
				{
					location:  location{column: 33, line: 0},
					value:     ")",
					tokenType: symbolType,
				},
			},
		},
		{
			input: "insert into users Values (105, 233)",
			Tokens: []token{
				{
					location:  location{column: 0, line: 0},
					value:     string(InsertKeyword),
					tokenType: keywordType,
				},
				{
					location:  location{column: 7, line: 0},
					value:     string(IntoKeyword),
					tokenType: keywordType,
				},
				{
					location:  location{column: 12, line: 0},
					value:     "users",
					tokenType: IdentifierType,
				},
				{
					location:  location{column: 18, line: 0},
					value:     string(ValuesKeyword),
					tokenType: keywordType,
				},
				{
					location:  location{column: 25, line: 0},
					value:     "(",
					tokenType: symbolType,
				},
				{
					location:  location{column: 26, line: 0},
					value:     "105",
					tokenType: NumericType,
				},
				{
					location:  location{column: 30, line: 0},
					value:     ",",
					tokenType: symbolType,
				},
				{
					location:  location{column: 32, line: 0},
					value:     "233",
					tokenType: NumericType,
				},
				{
					location:  location{column: 36, line: 0},
					value:     ")",
					tokenType: symbolType,
				},
			},
			err: nil,
		},
		{
			input: "SELECT id FROM users;",
			Tokens: []token{
				{
					location:  location{column: 0, line: 0},
					value:     string(SelectKeyword),
					tokenType: keywordType,
				},
				{
					location:  location{column: 7, line: 0},
					value:     "id",
					tokenType: IdentifierType,
				},
				{
					location:  location{column: 10, line: 0},
					value:     string(FromKeyword),
					tokenType: keywordType,
				},
				{
					location:  location{column: 15, line: 0},
					value:     "users",
					tokenType: IdentifierType,
				},
				{
					location:  location{column: 20, line: 0},
					value:     ";",
					tokenType: symbolType,
				},
			},
			err: nil,
		},
		{
			input: "SELECT * FROM my_table WHERE name = 'hello_world'",
			Tokens: []token{
				{
					location:  location{column: 0, line: 0},
					value:     string(SelectKeyword),
					tokenType: keywordType,
				},
				{
					location:  location{column: 7, line: 0},
					value:     "*",
					tokenType: symbolType,
				},
				{
					location:  location{column: 9, line: 0},
					value:     string(FromKeyword),
					tokenType: keywordType,
				},
				{
					location:  location{column: 14, line: 0},
					value:     "my_table",
					tokenType: IdentifierType,
				},
				{
					location:  location{column: 23, line: 0},
					value:     "where",
					tokenType: keywordType,
				},
				{
					location:  location{column: 29, line: 0},
					value:     "name",
					tokenType: IdentifierType,
				},
				{
					location:  location{column: 34, line: 0},
					value:     "=",
					tokenType: symbolType,
				},
				{
					location:  location{column: 36, line: 0},
					value:     "hello_world",
					tokenType: StringType,
				},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		tokens, err := Lex(test.input)
		assertEqual(t, test.err, err)
		assertEqual(t, len(test.Tokens), len(tokens))

		for i, tok := range tokens {
			assertEqual(t, &test.Tokens[i], tok)
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
