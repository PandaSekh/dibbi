package dibbi

import "errors"

// ColumnType defines the available column types
type ColumnType uint

const (
	TextType ColumnType = iota
	IntType
	BoolType
)

func (ct ColumnType) String() string {
	return [...]string{"Text", "Int", "Bool"}[ct]
}

type Cell interface {
	AsText() *string
	AsInt() *int32
	AsBool() *bool
}

type QueryResults struct {
	Columns []ResultColumn
	Rows    [][]Cell
}

type ResultColumn struct {
	Type    ColumnType
	Name    string
	NotNull bool
}

var (
	ErrTableDoesNotExist  = errors.New("table does not exist")
	ErrcolumnDoesNotExist = errors.New("column does not exist")
	ErrInvalidSelectItem  = errors.New("select item is not valid")
	ErrInvalidDatatype    = errors.New("invalid datatype")
	ErrMissingValues      = errors.New("missing values")
)

type Database interface {
	CreateTable(*createTableStatement) error
	Insert(*InsertStatement) error
	Select(*selectStatement) (*QueryResults, error)
}
