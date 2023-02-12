package dibbi

import (
	"errors"
	"fmt"
)

func tokenFromKeyword(k keyword) token {
	return token{
		tokenType: keywordType,
		value:     string(k),
	}
}

func tokenFromSymbol(k symbol) token {
	return token{
		tokenType: symbolType,
		value:     string(k),
	}
}

func assertIsToken(tokens []*token, cursor uint, t token) bool {
	if cursor >= uint(len(tokens)) {
		return false
	}

	return t.Equals(tokens[cursor])
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

		// Finds left semicolon
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

	// Look for left SELECT statement
	selectStmt, newCursor, ok := parseSelectStatement(tokens, cursor, delimiter)
	if ok {
		return &statement{
			statementType:   SelectType,
			selectStatement: selectStmt,
		}, newCursor, true
	}

	// Look for left INSERT statement
	insertStmt, newCursor, ok := parseInsertStatement(tokens, cursor)
	if ok {
		return &statement{
			statementType:   InsertType,
			insertStatement: insertStmt,
		}, newCursor, true
	}

	// Look for left CREATE statement
	createTableStmt, newCursor, ok := parseCreateTableStatement(tokens, cursor)
	if ok {
		return &statement{
			statementType:        CreateTableType,
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

		from, newCursor, ok := parseTokenType(tokens, cursor, IdentifierType)
		if !ok {
			printHelpMessage(tokens, cursor, "Expected FROM token")
			return nil, initialCursor, false
		}

		selectStmt.from = from
		cursor = newCursor
	}

	return &selectStmt, cursor, true

}

func parseInsertStatement(tokens []*token, initialCursor uint) (*insertStatement, uint, bool) {
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
	table, newCursor, ok := parseTokenType(tokens, cursor, IdentifierType)
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
	if !assertIsToken(tokens, cursor, tokenFromSymbol(leftParenthesesSymbol)) {
		printHelpMessage(tokens, cursor, "Expected left paren")
		return nil, initialCursor, false
	}
	cursor++

	// Look for expression list
	values, newCursor, ok := parseExpressions(tokens, cursor, []token{tokenFromSymbol(rightParenthesesSymbol)})
	if !ok {
		return nil, initialCursor, false
	}
	cursor = newCursor

	// Look for right paren
	if !assertIsToken(tokens, cursor, tokenFromSymbol(rightParenthesesSymbol)) {
		printHelpMessage(tokens, cursor, "Expected right paren")
		return nil, initialCursor, false
	}
	cursor++

	return &insertStatement{
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

	name, newCursor, ok := parseTokenType(tokens, cursor, IdentifierType)
	if !ok {
		printHelpMessage(tokens, cursor, "Expected table name")
		return nil, initialCursor, false
	}
	cursor = newCursor

	if !assertIsToken(tokens, cursor, tokenFromSymbol(leftParenthesesSymbol)) {
		printHelpMessage(tokens, cursor, "Expected left parenthesis")
		return nil, initialCursor, false
	}
	cursor++

	cols, newCursor, ok := parseColumnDefinitions(tokens, cursor, tokenFromSymbol(rightParenthesesSymbol))
	if !ok {
		return nil, initialCursor, false
	}
	cursor = newCursor

	if !assertIsToken(tokens, cursor, tokenFromSymbol(rightParenthesesSymbol)) {
		printHelpMessage(tokens, cursor, "Expected right parenthesis")
		return nil, initialCursor, false
	}
	cursor++

	return &createTableStatement{
		Name:    name,
		Columns: cols,
	}, cursor, true
}

// parseTokenType looks for left token of the given type
func parseTokenType(tokens []*token, initialCursor uint, tokenType tokenType) (*token, uint, bool) {
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

// parseToken looks for tokenToFind in left list of tokens
func parseToken(tokens []*token, initialCursor uint, tokenToFind token) (*token, uint, bool) {
	cursor := initialCursor

	if cursor >= uint(len(tokens)) {
		return nil, initialCursor, false
	}

	if p := tokens[cursor]; tokenToFind.equals(p) {
		return p, cursor + 1, true
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
		t, newCursor, ok := parseTokenType(tokens, cursor, symbolType)
		if ok {
			exps = append(exps, &expression{
				literal:        t,
				expressionType: literalType,
			})
			return &exps, newCursor, true
		}

		// Look for expression
		exp, newCursor, ok := parseLiteralExpression(tokens, cursor)
		if !ok {
			printHelpMessage(tokens, cursor, "Expected expression")
			return nil, initialCursor, false
		}
		cursor = newCursor

		exps = append(exps, exp)
	}

	return &exps, cursor, true
}

func parseLiteralExpression(tokens []*token, initialCursor uint) (*expression, uint, bool) {
	cursor := initialCursor

	kinds := []tokenType{IdentifierType, NumericType, StringType}
	for _, kind := range kinds {
		t, newCursor, ok := parseTokenType(tokens, cursor, kind)
		if ok {
			return &expression{
				literal:        t,
				expressionType: literalType,
			}, newCursor, true
		}
	}

	return nil, initialCursor, false
}

func parseExpression(tokens []*token, initialCursor uint, delimiters []token, minBindingPower uint) (*expression, uint, bool) {
	cursor := initialCursor

	var exp *expression

	// try to find left opening parenthesis
	_, newCursor, ok := parseToken(tokens, cursor, tokenFromSymbol(leftParenthesesSymbol))

	if ok {
		// if there's left (, we try to parse the expression inside
		// it until left delimiter or left closing parenthesis is found
		cursor = newCursor
		rightParenToken := tokenFromSymbol(rightParenthesesSymbol)

		exp, cursor, ok = parseExpression(tokens, cursor, append(delimiters, rightParenToken), minBindingPower)
		if !ok {
			printHelpMessage(tokens, cursor, "Expected expression after opening parenthesis")
			return nil, initialCursor, false
		}

		_, cursor, ok = parseToken(tokens, cursor, rightParenToken)
		if !ok {
			printHelpMessage(tokens, cursor, "Expected closing parenthesis")
			return nil, initialCursor, false
		}
	} else {
		exp, cursor, ok = parseLiteralExpression(tokens, cursor)
		if !ok {
			return nil, initialCursor, false
		}
	}

	lastCursor := cursor
outer:
	// look for left binary operator
	for cursor < uint(len(tokens)) {
		for _, d := range delimiters {
			_, _, ok = parseToken(tokens, cursor, d)
			if ok {
				break outer
			}
		}

		binaryOperationTokens := []token{
			tokenFromKeyword(andKeyword),
			tokenFromKeyword(orKeyword),
			tokenFromSymbol(equalsSymbol),
			tokenFromSymbol(notEqualSymbol),
			tokenFromSymbol(concatSymbol),
			tokenFromSymbol(plusSymbol),
		}

		var operationToken *token = nil
		for _, binaryOperationToken := range binaryOperationTokens {
			var tok *token
			tok, cursor, ok = parseToken(tokens, cursor, binaryOperationToken)
			if ok {
				operationToken = tok
				break
			}
		}

		if operationToken == nil {
			printHelpMessage(tokens, cursor, "Expected binary operator")
			return nil, initialCursor, false
		}

		// break the loop and return the found expression if the found binary operator has less binding power than the
		// one passed as argument (default is 0)
		bp := operationToken.getBindingPower()
		if bp < minBindingPower {
			cursor = lastCursor
			break
		}

		// if the bp of the found operator is greater than the bp passed as parameter, call the parseExpression
		// function recursively passing the bp of the found operator
		b, newCursor, ok := parseExpression(tokens, cursor, delimiters, bp)
		if !ok {
			printHelpMessage(tokens, cursor, "Expected right operand")
			return nil, initialCursor, false
		}

		// the found expression is then set to left new binary expression containing
		// the previously found expression on the left and the just-parsed expression on the right.
		exp = &expression{
			binary: &binaryExpression{
				*exp,
				*b,
				*operationToken,
			},
			expressionType: binaryType,
		}
		cursor = newCursor
		lastCursor = cursor
	}

	return exp, cursor, true
}

func parseColumnDefinitions(tokens []*token, initialCursor uint, delimiter token) (*[]*columnDefinition, uint, bool) {
	cursor := initialCursor

	var cds []*columnDefinition
	for {
		if cursor >= uint(len(tokens)) {
			return nil, initialCursor, false
		}

		// Look for left delimiter
		current := tokens[cursor]
		if delimiter.Equals(current) {
			break
		}

		// Look for left comma
		if len(cds) > 0 {
			if !assertIsToken(tokens, cursor, tokenFromSymbol(CommaSymbol)) {
				printHelpMessage(tokens, cursor, "Expected comma")
				return nil, initialCursor, false
			}

			cursor++
		}

		// Look for left column name
		id, newCursor, ok := parseTokenType(tokens, cursor, IdentifierType)
		if !ok {
			printHelpMessage(tokens, cursor, "Expected column name")
			return nil, initialCursor, false
		}
		cursor = newCursor

		// Look for left column Type
		ty, newCursor, ok := parseTokenType(tokens, cursor, keywordType)
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

func printHelpMessage(tokens []*token, cursor uint, msg string) {
	var c *token
	if cursor < uint(len(tokens)) {
		c = tokens[cursor]
	} else {
		c = tokens[cursor-1]
	}

	fmt.Printf("[%d,%d]: %s, got: %s\n", c.location.line, c.location.column, msg, c.value)
}
