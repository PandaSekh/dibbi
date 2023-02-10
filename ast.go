package dibbi

import "fmt"

type ast struct {
	Statements []*statement
}

func (a *ast) String() string {
	return fmt.Sprintf("%v", a.Statements)
}

type astType uint

const (
	SelectType astType = iota
	CreateTableType
	InsertType
)

func (at astType) String() string {
	return [...]string{"Select", "Create Table", "Insert"}[at]
}

type statement struct {
	selectStatement      *selectStatement
	createTableStatement *createTableStatement
	InsertStatement      *InsertStatement
	Type                 astType
}

func (a *statement) String() string {
	return fmt.Sprintf("Type: %v\nSelect: '%v'\nInsert: %v\n", a.
		Type, a.selectStatement, a.InsertStatement)
}

// Insert

type InsertStatement struct {
	Table  token
	Values *[]*expression
}

func (s *InsertStatement) String() string {
	return fmt.Sprintf("Table: %v, Values: %v", s.Table.value, s.Values)
}

type ExpressionType uint

const (
	LiteralType ExpressionType = iota
)

func (lt ExpressionType) String() string {
	return [...]string{"Literal"}[lt]
}

type expression struct {
	Literal        *token
	ExpressionType ExpressionType
}

func (e *expression) String() string {
	return fmt.Sprintf("Literal: %v, Type: %v", e.Literal.value, e.ExpressionType)
}

// Create

type columnDefinition struct {
	Name     *token
	Datatype *token
}

type createTableStatement struct {
	Name    *token
	Columns *[]*columnDefinition
}

// Select

type selectStatement struct {
	from  *token
	items []*expression
}

func (s *selectStatement) String() string {
	return fmt.Sprintf("From: %v, Items: %v", s.from.value, s.items)
}
