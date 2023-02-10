package dibbi

import (
	"errors"
	"fmt"
)

func tokenFromKeyword(k keyword) token {
	return token{
		tokenType: KeywordType,
		value:     string(k),
	}
}

func tokenFromSymbol(k symbol) token {
	return token{
		tokenType: SymbolType,
		value:     string(k),
	}
}

func assertIsToken(tokens []*token, cursor uint, t token) bool {
	if cursor >= uint(len(tokens)) {
		return false
	}

	return t.Equals(tokens[cursor])
}

func printHelpMessage(tokens []*token, cursor uint, msg string) {
	var c *token
	if cursor < uint(len(tokens)) {
		c = tokens[cursor]
	} else {
		c = tokens[cursor-1]
	}

	fmt.Printf("[%d,%d]: %s, got: %s\n", c.location.line, c.location.column, msg, c.value)
}

func parse(source string) (*ast, error) {
	tokens, err := Lex(source)
	if err != nil {
		return nil, err
	}

	ast := ast{}
	cursor := uint(0)
	for cursor < uint(len(tokens)) {
		stmt, newCursor, ok := parseStatement(tokens, cursor, tokenFromSymbol(SemicolonSymbol))
		if !ok {
			printHelpMessage(tokens, cursor, "Expected statement")
			return nil, errors.New("failed to parse, expected statement")
		}
		cursor = newCursor

		ast.Statements = append(ast.Statements, stmt)

		// Finds a semicolon
		semicolonIsPresent := false
		for assertIsToken(tokens, cursor, tokenFromSymbol(SemicolonSymbol)) {
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

func parseStatement(tokens []*token, initialCursor uint, delimiter token) (*statement, uint, bool) {
	cursor := initialCursor

	// Look for a SELECT statement
	selectStmt, newCursor, ok := parseSelectStatement(tokens, cursor, delimiter)
	if ok {
		return &statement{
			Type:            SelectType,
			selectStatement: selectStmt,
		}, newCursor, true
	}

	// Look for a INSERT statement
	insertStmt, newCursor, ok := parseInsertStatement(tokens, cursor)
	if ok {
		return &statement{
			Type:            InsertType,
			InsertStatement: insertStmt,
		}, newCursor, true
	}

	// Look for a CREATE statement
	createTableStmt, newCursor, ok := parseCreateTableStatement(tokens, cursor)
	if ok {
		return &statement{
			Type:                 CreateTableType,
			createTableStatement: createTableStmt,
		}, newCursor, true
	}

	return nil, initialCursor, false
}

func parseSelectStatement(tokens []*token, initialCursor uint, delimiter token) (*selectStatement, uint, bool) {
	cursor := initialCursor
	if !assertIsToken(tokens, cursor, tokenFromKeyword(SelectKeyword)) {
		return nil, initialCursor, false
	}

	cursor++

	selectStmt := selectStatement{}

	expressions, newCursor, ok := parseExpressions(tokens, cursor, []token{tokenFromKeyword(FromKeyword), delimiter})
	if !ok {
		return nil, initialCursor, false
	}

	selectStmt.items = *expressions
	cursor = newCursor

	if assertIsToken(tokens, cursor, tokenFromKeyword(FromKeyword)) {
		cursor++

		from, newCursor, ok := parseToken(tokens, cursor, IdentifierType)
		if !ok {
			printHelpMessage(tokens, cursor, "Expected FROM token")
			return nil, initialCursor, false
		}

		selectStmt.from = from
		cursor = newCursor
	}

	return &selectStmt, cursor, true

}

func parseInsertStatement(tokens []*token, initialCursor uint) (*InsertStatement, uint, bool) {
	cursor := initialCursor

	// Look for INSERT
	if !assertIsToken(tokens, cursor, tokenFromKeyword(InsertKeyword)) {
		return nil, initialCursor, false
	}
	cursor++

	// Look for INTO
	if !assertIsToken(tokens, cursor, tokenFromKeyword(IntoKeyword)) {
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
	if !assertIsToken(tokens, cursor, tokenFromKeyword(ValuesKeyword)) {
		printHelpMessage(tokens, cursor, "Expected VALUES")
		return nil, initialCursor, false
	}
	cursor++

	// Look for left paren
	if !assertIsToken(tokens, cursor, tokenFromSymbol(LeftParenthesesSymbol)) {
		printHelpMessage(tokens, cursor, "Expected left paren")
		return nil, initialCursor, false
	}
	cursor++

	// Look for expression list
	values, newCursor, ok := parseExpressions(tokens, cursor, []token{tokenFromSymbol(RightParenthesesSymbol)})
	if !ok {
		return nil, initialCursor, false
	}
	cursor = newCursor

	// Look for right paren
	if !assertIsToken(tokens, cursor, tokenFromSymbol(RightParenthesesSymbol)) {
		printHelpMessage(tokens, cursor, "Expected right paren")
		return nil, initialCursor, false
	}
	cursor++

	return &InsertStatement{
		Table:  *table,
		Values: values,
	}, cursor, true
}

func parseCreateTableStatement(tokens []*token, initialCursor uint) (*createTableStatement, uint, bool) {
	cursor := initialCursor

	if !assertIsToken(tokens, cursor, tokenFromKeyword(CreateKeyword)) {
		return nil, initialCursor, false
	}
	cursor++

	if !assertIsToken(tokens, cursor, tokenFromKeyword(TableKeyword)) {
		return nil, initialCursor, false
	}
	cursor++

	name, newCursor, ok := parseToken(tokens, cursor, IdentifierType)
	if !ok {
		printHelpMessage(tokens, cursor, "Expected table name")
		return nil, initialCursor, false
	}
	cursor = newCursor

	if !assertIsToken(tokens, cursor, tokenFromSymbol(LeftParenthesesSymbol)) {
		printHelpMessage(tokens, cursor, "Expected left parenthesis")
		return nil, initialCursor, false
	}
	cursor++

	cols, newCursor, ok := parsecolumnDefinitions(tokens, cursor, tokenFromSymbol(RightParenthesesSymbol))
	if !ok {
		return nil, initialCursor, false
	}
	cursor = newCursor

	if !assertIsToken(tokens, cursor, tokenFromSymbol(RightParenthesesSymbol)) {
		printHelpMessage(tokens, cursor, "Expected right parenthesis")
		return nil, initialCursor, false
	}
	cursor++

	return &createTableStatement{
		Name:    name,
		Columns: cols,
	}, cursor, true
}

func parseToken(tokens []*token, initialCursor uint, tokenType tokenType) (*token, uint, bool) {
	cursor := initialCursor

	if cursor >= uint(len(tokens)) {
		return nil, initialCursor, false
	}

	current := tokens[cursor]
	if current.tokenType == tokenType {
		return current, cursor + 1, true
	}

	return nil, initialCursor, false
}

func parseExpressions(tokens []*token, initialCursor uint, delimiters []token) (*[]*expression, uint, bool) {
	cursor := initialCursor

	var exps []*expression
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
			if !assertIsToken(tokens, cursor, tokenFromSymbol(CommaSymbol)) {
				printHelpMessage(tokens, cursor, "Expected comma")
				return nil, initialCursor, false
			}

			cursor++
		}

		// Check if it's star
		t, newCursor, ok := parseToken(tokens, cursor, SymbolType)
		if ok {
			exps = append(exps, &expression{
				Literal:        t,
				ExpressionType: LiteralType,
			})
			return &exps, newCursor, true
		}

		// Look for expression
		exp, newCursor, ok := parseExpression(tokens, cursor, tokenFromSymbol(CommaSymbol))
		if !ok {
			printHelpMessage(tokens, cursor, "Expected expression")
			return nil, initialCursor, false
		}
		cursor = newCursor

		exps = append(exps, exp)
	}

	return &exps, cursor, true
}

func parseExpression(tokens []*token, initialCursor uint, _ token) (*expression, uint, bool) {
	cursor := initialCursor

	kinds := []tokenType{IdentifierType, NumericType, StringType}
	for _, kind := range kinds {
		t, newCursor, ok := parseToken(tokens, cursor, kind)
		if ok {
			return &expression{
				Literal:        t,
				ExpressionType: LiteralType,
			}, newCursor, true
		}
	}

	return nil, initialCursor, false
}

func parsecolumnDefinitions(tokens []*token, initialCursor uint, delimiter token) (*[]*columnDefinition, uint, bool) {
	cursor := initialCursor

	var cds []*columnDefinition
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
			if !assertIsToken(tokens, cursor, tokenFromSymbol(CommaSymbol)) {
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

		// Look for a column Type
		ty, newCursor, ok := parseToken(tokens, cursor, KeywordType)
		if !ok {
			printHelpMessage(tokens, cursor, "Expected column Type")
			return nil, initialCursor, false
		}
		cursor = newCursor

		cds = append(cds, &columnDefinition{
			Name:     id,
			Datatype: ty,
		})
	}

	return &cds, cursor, true
}
