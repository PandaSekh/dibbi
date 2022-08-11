package ast

import "go_dibbi/lexer"

type Ast struct {
	Statements []*Statement
}

type AstType uint

const (
	SelectType AstType = iota
	CreateTableType
	InsertType
)

type Statement struct {
	SelectStatement      *SelectStatement
	CreateTableStatement *CreateTableStatement
	InsertStatement      *InsertStatement
	Type                 AstType
}

// Insert

type InsertStatement struct {
	Table  lexer.Token
	Values *[]*Expression
}

type ExpressionType uint

const (
	LiteralType ExpressionType = iota
)

type Expression struct {
	Literal        *lexer.Token
	ExpressionType ExpressionType
}

// Create

type ColumnDefinition struct {
	Name     *lexer.Token
	Datatype *lexer.Token
}

type CreateTableStatement struct {
	Name    *lexer.Token
	Columns *[]*ColumnDefinition
}

// Select

type SelectStatement struct {
	From  *lexer.Token
	Items []*Expression
}
