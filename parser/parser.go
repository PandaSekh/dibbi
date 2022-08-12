package parser

import (
	"errors"
	"fmt"
	"go_dibbi/ast"
	"go_dibbi/lexer"
)

func TokenFromKeyword(k lexer.Keyword) lexer.Token {
	return lexer.Token{
		Type:  lexer.KeywordType,
		Value: string(k),
	}
}

func TokenFromSymbol(k lexer.Symbol) lexer.Token {
	return lexer.Token{
		Type:  lexer.SymbolType,
		Value: string(k),
	}
}

func assertIsToken(tokens []*lexer.Token, cursor uint, t lexer.Token) bool {
	if cursor >= uint(len(tokens)) {
		return false
	}

	return t.Equals(tokens[cursor])
}

func printHelpMessage(tokens []*lexer.Token, cursor uint, msg string) {
	var c *lexer.Token
	if cursor < uint(len(tokens)) {
		c = tokens[cursor]
	} else {
		c = tokens[cursor-1]
	}

	fmt.Printf("[%d,%d]: %s, got: %s\n", c.Location.Line, c.Location.Column, msg, c.Value)
}

func Parse(source string) (*ast.Ast, error) {
	tokens, err := lexer.Lex(source)
	if err != nil {
		return nil, err
	}

	ast := ast.Ast{}
	cursor := uint(0)
	for cursor < uint(len(tokens)) {
		stmt, newCursor, ok := parseStatement(tokens, cursor, TokenFromSymbol(lexer.SemicolonSymbol))
		if !ok {
			printHelpMessage(tokens, cursor, "Expected statement")
			return nil, errors.New("failed to parse, expected statement")
		}
		cursor = newCursor

		ast.Statements = append(ast.Statements, stmt)

		// Finds a semicolon
		semicolonIsPresent := false
		for assertIsToken(tokens, cursor, TokenFromSymbol(lexer.SemicolonSymbol)) {
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

func parseStatement(tokens []*lexer.Token, initialCursor uint, delimiter lexer.Token) (*ast.Statement, uint, bool) {
	cursor := initialCursor

	// Look for a SELECT statement
	selectStmt, newCursor, ok := parseSelectStatement(tokens, cursor, delimiter)
	if ok {
		return &ast.Statement{
			Type:            ast.SelectType,
			SelectStatement: selectStmt,
		}, newCursor, true
	}

	// Look for a INSERT statement
	insertStmt, newCursor, ok := parseInsertStatement(tokens, cursor, delimiter)
	if ok {
		return &ast.Statement{
			Type:            ast.InsertType,
			InsertStatement: insertStmt,
		}, newCursor, true
	}

	// Look for a CREATE statement
	createTableStmt, newCursor, ok := parseCreateTableStatement(tokens, cursor, delimiter)
	if ok {
		return &ast.Statement{
			Type:                 ast.CreateTableType,
			CreateTableStatement: createTableStmt,
		}, newCursor, true
	}

	return nil, initialCursor, false
}

func parseSelectStatement(tokens []*lexer.Token, initialCursor uint, delimiter lexer.Token) (*ast.SelectStatement, uint, bool) {
	return nil, 0, false
}

func parseInsertStatement(tokens []*lexer.Token, initialCursor uint, delimiter lexer.Token) (*ast.InsertStatement, uint, bool) {
	return nil, 0, false
}

func parseCreateTableStatement(tokens []*lexer.Token, initialCursor uint, delimiter lexer.Token) (*ast.CreateTableStatement, uint, bool) {
	return nil, 0, false
}
