package internal

import "fmt"

type Ast struct {
	Statements []*Statement
}

func (a *Ast) String() string {
	return fmt.Sprintf("%v", a.Statements)
}

type AstType uint

const (
	SelectType AstType = iota
	CreateTableType
	InsertType
)

func (at AstType) String() string {
	return [...]string{"Select", "Create Table", "Insert"}[at]
}

type Statement struct {
	SelectStatement      *SelectStatement
	CreateTableStatement *CreateTableStatement
	InsertStatement      *InsertStatement
	Type                 AstType
}

func (a *Statement) String() string {
	return fmt.Sprintf("Type: %v\nSelect: '%v'\nInsert: %v\n", a.Type, a.SelectStatement, a.InsertStatement)
}

// Insert

type InsertStatement struct {
	Table  Token
	Values *[]*Expression
}

func (s *InsertStatement) String() string {
	return fmt.Sprintf("Table: %v, Values: %v", s.Table.Value, s.Values)
}

type ExpressionType uint

const (
	LiteralType ExpressionType = iota
)

func (lt ExpressionType) String() string {
	return [...]string{"Literal"}[lt]
}

type Expression struct {
	Literal        *Token
	ExpressionType ExpressionType
}

func (e *Expression) String() string {
	return fmt.Sprintf("Literal: %v, Type: %v", e.Literal.Value, e.ExpressionType)
}

// Create

type ColumnDefinition struct {
	Name     *Token
	Datatype *Token
}

type CreateTableStatement struct {
	Name    *Token
	Columns *[]*ColumnDefinition
}

// Select

type SelectStatement struct {
	From  *Token
	Items []*Expression
}

func (s *SelectStatement) String() string {
	return fmt.Sprintf("From: %v, Items: %v", s.From.Value, s.Items)
}
