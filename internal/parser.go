package internal

import (
	"errors"
	"fmt"
)

func TokenFromKeyword(k Keyword) Token {
	return Token{
		Type:  KeywordType,
		Value: string(k),
	}
}

func TokenFromSymbol(k Symbol) Token {
	return Token{
		Type:  SymbolType,
		Value: string(k),
	}
}

func assertIsToken(tokens []*Token, cursor uint, t Token) bool {
	if cursor >= uint(len(tokens)) {
		return false
	}

	return t.Equals(tokens[cursor])
}

func printHelpMessage(tokens []*Token, cursor uint, msg string) {
	var c *Token
	if cursor < uint(len(tokens)) {
		c = tokens[cursor]
	} else {
		c = tokens[cursor-1]
	}

	fmt.Printf("[%d,%d]: %s, got: %s\n", c.Location.Line, c.Location.Column, msg, c.Value)
}

func Parse(source string) (*Ast, error) {
	tokens, err := Lex(source)
	if err != nil {
		return nil, err
	}

	ast := Ast{}
	cursor := uint(0)
	for cursor < uint(len(tokens)) {
		stmt, newCursor, ok := parseStatement(tokens, cursor, TokenFromSymbol(SemicolonSymbol))
		if !ok {
			printHelpMessage(tokens, cursor, "Expected statement")
			return nil, errors.New("failed to parse, expected statement")
		}
		cursor = newCursor

		ast.Statements = append(ast.Statements, stmt)

		// Finds a semicolon
		semicolonIsPresent := false
		for assertIsToken(tokens, cursor, TokenFromSymbol(SemicolonSymbol)) {
			cursor++
			semicolonIsPresent = true
		}

		if !semicolonIsPresent {
			printHelpMessage(tokens, cursor, "Expected semi-colon delimiter between statements")
			return nil, errors.New("missing semi-colon between statements")
		}
	}

	return &ast, nil
}

func parseStatement(tokens []*Token, initialCursor uint, delimiter Token) (*Statement, uint, bool) {
	cursor := initialCursor

	// Look for a SELECT statement
	selectStmt, newCursor, ok := parseSelectStatement(tokens, cursor, delimiter)
	if ok {
		return &Statement{
			Type:            SelectType,
			SelectStatement: selectStmt,
		}, newCursor, true
	}

	// Look for a INSERT statement
	insertStmt, newCursor, ok := parseInsertStatement(tokens, cursor)
	if ok {
		return &Statement{
			Type:            InsertType,
			InsertStatement: insertStmt,
		}, newCursor, true
	}

	// Look for a CREATE statement
	createTableStmt, newCursor, ok := parseCreateTableStatement(tokens, cursor)
	if ok {
		return &Statement{
			Type:                 CreateTableType,
			CreateTableStatement: createTableStmt,
		}, newCursor, true
	}

	return nil, initialCursor, false
}

func parseSelectStatement(tokens []*Token, initialCursor uint, delimiter Token) (*SelectStatement, uint, bool) {
	cursor := initialCursor
	if !assertIsToken(tokens, cursor, TokenFromKeyword(SelectKeyword)) {
		return nil, initialCursor, false
	}

	cursor++

	selectStmt := SelectStatement{}

	expressions, newCursor, ok := parseExpressions(tokens, cursor, []Token{TokenFromKeyword(FromKeyword), delimiter})
	if !ok {
		return nil, initialCursor, false
	}

	selectStmt.Items = *expressions
	cursor = newCursor

	if assertIsToken(tokens, cursor, TokenFromKeyword(FromKeyword)) {
		cursor++

		from, newCursor, ok := parseToken(tokens, cursor, IdentifierType)
		if !ok {
			printHelpMessage(tokens, cursor, "Expected FROM token")
			return nil, initialCursor, false
		}

		selectStmt.From = from
		cursor = newCursor
	}

	return &selectStmt, cursor, true

}

func parseInsertStatement(tokens []*Token, initialCursor uint) (*InsertStatement, uint, bool) {
	cursor := initialCursor

	// Look for INSERT
	if !assertIsToken(tokens, cursor, TokenFromKeyword(InsertKeyword)) {
		return nil, initialCursor, false
	}
	cursor++

	// Look for INTO
	if !assertIsToken(tokens, cursor, TokenFromKeyword(IntoKeyword)) {
		printHelpMessage(tokens, cursor, "Expected into")
		return nil, initialCursor, false
	}
	cursor++

	// Look for table name
	table, newCursor, ok := parseToken(tokens, cursor, IdentifierType)
	if !ok {
		printHelpMessage(tokens, cursor, "Expected table name")
		return nil, initialCursor, false
	}
	cursor = newCursor

	// Look for VALUES
	if !assertIsToken(tokens, cursor, TokenFromKeyword(ValuesKeyword)) {
		printHelpMessage(tokens, cursor, "Expected VALUES")
		return nil, initialCursor, false
	}
	cursor++

	// Look for left paren
	if !assertIsToken(tokens, cursor, TokenFromSymbol(LeftParenthesesSymbol)) {
		printHelpMessage(tokens, cursor, "Expected left paren")
		return nil, initialCursor, false
	}
	cursor++

	// Look for expression list
	values, newCursor, ok := parseExpressions(tokens, cursor, []Token{TokenFromSymbol(RightParenthesesSymbol)})
	if !ok {
		return nil, initialCursor, false
	}
	cursor = newCursor

	// Look for right paren
	if !assertIsToken(tokens, cursor, TokenFromSymbol(RightParenthesesSymbol)) {
		printHelpMessage(tokens, cursor, "Expected right paren")
		return nil, initialCursor, false
	}
	cursor++

	return &InsertStatement{
		Table:  *table,
		Values: values,
	}, cursor, true
}

func parseCreateTableStatement(tokens []*Token, initialCursor uint) (*CreateTableStatement, uint, bool) {
	cursor := initialCursor

	if !assertIsToken(tokens, cursor, TokenFromKeyword(CreateKeyword)) {
		return nil, initialCursor, false
	}
	cursor++

	if !assertIsToken(tokens, cursor, TokenFromKeyword(TableKeyword)) {
		return nil, initialCursor, false
	}
	cursor++

	name, newCursor, ok := parseToken(tokens, cursor, IdentifierType)
	if !ok {
		printHelpMessage(tokens, cursor, "Expected table name")
		return nil, initialCursor, false
	}
	cursor = newCursor

	if !assertIsToken(tokens, cursor, TokenFromSymbol(LeftParenthesesSymbol)) {
		printHelpMessage(tokens, cursor, "Expected left parenthesis")
		return nil, initialCursor, false
	}
	cursor++

	cols, newCursor, ok := parseColumnDefinitions(tokens, cursor, TokenFromSymbol(RightParenthesesSymbol))
	if !ok {
		return nil, initialCursor, false
	}
	cursor = newCursor

	if !assertIsToken(tokens, cursor, TokenFromSymbol(RightParenthesesSymbol)) {
		printHelpMessage(tokens, cursor, "Expected right parenthesis")
		return nil, initialCursor, false
	}
	cursor++

	return &CreateTableStatement{
		Name:    name,
		Columns: cols,
	}, cursor, true
}

func parseToken(tokens []*Token, initialCursor uint, tokenType TokenType) (*Token, uint, bool) {
	cursor := initialCursor

	if cursor >= uint(len(tokens)) {
		return nil, initialCursor, false
	}

	current := tokens[cursor]
	if current.Type == tokenType {
		return current, cursor + 1, true
	}

	return nil, initialCursor, false
}

func parseExpressions(tokens []*Token, initialCursor uint, delimiters []Token) (*[]*Expression, uint, bool) {
	cursor := initialCursor

	var exps []*Expression
outer:
	for {
		if cursor >= uint(len(tokens)) {
			return nil, initialCursor, false
		}

		// Look for delimiter
		current := tokens[cursor]
		for _, delimiter := range delimiters {
			if delimiter.Equals(current) {
				break outer
			}
		}

		// Look for comma
		if len(exps) > 0 {
			if !assertIsToken(tokens, cursor, TokenFromSymbol(CommaSymbol)) {
				printHelpMessage(tokens, cursor, "Expected comma")
				return nil, initialCursor, false
			}

			cursor++
		}

		// Check if it's star
		t, newCursor, ok := parseToken(tokens, cursor, SymbolType)
		if ok {
			exps = append(exps, &Expression{
				Literal:        t,
				ExpressionType: LiteralType,
			})
			return &exps, newCursor, true
		}

		// Look for expression
		exp, newCursor, ok := parseExpression(tokens, cursor, TokenFromSymbol(CommaSymbol))
		if !ok {
			printHelpMessage(tokens, cursor, "Expected expression")
			return nil, initialCursor, false
		}
		cursor = newCursor

		exps = append(exps, exp)
	}

	return &exps, cursor, true
}

func parseExpression(tokens []*Token, initialCursor uint, _ Token) (*Expression, uint, bool) {
	cursor := initialCursor

	kinds := []TokenType{IdentifierType, NumericType, StringType}
	for _, kind := range kinds {
		t, newCursor, ok := parseToken(tokens, cursor, kind)
		if ok {
			return &Expression{
				Literal:        t,
				ExpressionType: LiteralType,
			}, newCursor, true
		}
	}

	return nil, initialCursor, false
}

func parseColumnDefinitions(tokens []*Token, initialCursor uint, delimiter Token) (*[]*ColumnDefinition, uint, bool) {
	cursor := initialCursor

	var cds []*ColumnDefinition
	for {
		if cursor >= uint(len(tokens)) {
			return nil, initialCursor, false
		}

		// Look for a delimiter
		current := tokens[cursor]
		if delimiter.Equals(current) {
			break
		}

		// Look for a comma
		if len(cds) > 0 {
			if !assertIsToken(tokens, cursor, TokenFromSymbol(CommaSymbol)) {
				printHelpMessage(tokens, cursor, "Expected comma")
				return nil, initialCursor, false
			}

			cursor++
		}

		// Look for a column name
		id, newCursor, ok := parseToken(tokens, cursor, IdentifierType)
		if !ok {
			printHelpMessage(tokens, cursor, "Expected column name")
			return nil, initialCursor, false
		}
		cursor = newCursor

		// Look for a column type
		ty, newCursor, ok := parseToken(tokens, cursor, KeywordType)
		if !ok {
			printHelpMessage(tokens, cursor, "Expected column type")
			return nil, initialCursor, false
		}
		cursor = newCursor

		cds = append(cds, &ColumnDefinition{
			Name:     id,
			Datatype: ty,
		})
	}

	return &cds, cursor, true
}
