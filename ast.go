package dibbi

import "fmt"

type ast struct {
	Statements []*statement
}

func (a *ast) String() string {
	return fmt.Sprintf("%v", a.Statements)
}

type statementType uint

const (
	SelectType statementType = iota
	CreateTableType
	InsertType
)

func (at statementType) String() string {
	return [...]string{"Select", "Create Table", "Insert"}[at]
}

type statement struct {
	selectStatement      *selectStatement
	createTableStatement *createTableStatement
	insertStatement      *insertStatement
	statementType        statementType
}

func (a *statement) String() string {
	return fmt.Sprintf("Type: %v\nSelect: '%v'\nInsert: %v\n", a.
		statementType, a.selectStatement, a.insertStatement)
}

// Insert

type insertStatement struct {
	Table  token
	Values *[]*expression
}

func (s *insertStatement) String() string {
	return fmt.Sprintf("Table: %v, Values: %v", s.Table.value, s.Values)
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
	items *[]*selectItem
	where *expression
}

type selectItem struct {
	exp      *expression
	asterisk bool
	as       *token
}

func (s *selectStatement) String() string {
	return fmt.Sprintf("From: %v, Items: %v", s.from.value, s.items)
}

// Expressions
type expressionType uint

const (
	literalType expressionType = iota
	binaryType
)

func (lt expressionType) String() string {
	return [...]string{"literal"}[lt]
}

type binaryExpression struct {
	left     expression
	right    expression
	operator token
}

type expression struct {
	literal        *token
	binary         *binaryExpression
	expressionType expressionType
}

func (e *expression) String() string {
	return fmt.Sprintf("literal: %v, Type: %v", e.literal.value, e.expressionType)
}
