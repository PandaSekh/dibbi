package dibbi

import "errors"

// ColumnType defines the available column types
type ColumnType uint

const (
	TextType ColumnType = iota
	IntType
	BoolType
)

func (c ColumnType) String() string {
	switch c {
	case TextType:
		return "TextType"
	case IntType:
		return "IntType"
	case BoolType:
		return "BoolType"
	default:
		return "Error"
	}
}

type Cell interface {
	AsText() *string
	AsInt() *int32
	AsBool() *bool
}

type Results struct {
	Columns []ResultColumn
	Rows    [][]Cell
}

type ResultColumn struct {
	Type    ColumnType
	Name    string
	NotNull bool
}

type Index struct {
	Name       string
	Exp        string
	Type       string
	Unique     bool
	PrimaryKey bool
}

var (
	ErrTableDoesNotExist  = errors.New("table does not exist")
	ErrColumnDoesNotExist = errors.New("column does not exist")
	ErrInvalidSelectItem  = errors.New("select item is not valid")
	ErrInvalidDatatype    = errors.New("invalid datatype")
	ErrMissingValues      = errors.New("missing values")
)

type TableMetadata struct {
	Name    string
	Columns []ResultColumn
	Indexes []Index
}

type Database interface {
	CreateTable(*createTableStatement) error
	//DropTable(*dropTableStatement) error
	//CreateIndex(*createIndexStatement) error
	Insert(*InsertStatement) error
	Select(*selectStatement) (*Results, error)
	//GetTables() []TableMetadata
}
