package internal

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"
	"strconv"
)

type MemoryCell []byte

func (mc MemoryCell) AsInt() int32 {
	var i int32
	err := binary.Read(bytes.NewBuffer(mc), binary.BigEndian, &i)
	if err != nil {
		fmt.Printf("Corrupted data [%s]: %s\n", mc, err)
		return 0
	}

	return i
}

func (mc MemoryCell) AsText() string {
	return string(mc)
}

func (mc MemoryCell) AsBool() bool {
	b, err := strconv.ParseBool(string(mc))
	if err != nil {
		return false
	}
	return b
}

func (mc MemoryCell) equals(b MemoryCell) bool {
	if mc == nil || b == nil {
		return mc == nil && b == nil
	}

	return bytes.Equal(mc, b)
}

// literalToMemoryCell maps a Go value into a memory cell
func literalToMemoryCell(t *Token) MemoryCell {
	if t.Type == NumericType {
		buf := new(bytes.Buffer)
		i, err := strconv.Atoi(t.Value)
		if err != nil {
			fmt.Printf("Corrupted data [%s]: %s\n", t.Value, err)
			return nil
		}

		err = binary.Write(buf, binary.BigEndian, new(big.Int).SetInt64(int64(i)).Bytes())
		if err != nil {
			fmt.Printf("Corrupted data [%s]: %s\n", buf.String(), err)
			return nil
		}
		return buf.Bytes()
	}

	if t.Type == StringType {
		return MemoryCell(t.Value)
	}

	if t.Type == BooleanType {
		if t.Value == "true" {
			return []byte{1}
		}

		return []byte{0}
	}

	return nil
}

type dbTable struct {
	columns     []string
	columnTypes []ColumnType
	rows        [][]MemoryCell
}

type MemoryBackend struct {
	tables map[string]*dbTable
}

func NewMemoryBackend() *MemoryBackend {
	return &MemoryBackend{
		tables: map[string]*dbTable{},
	}
}

func (mb *MemoryBackend) CreateTable(crt *CreateTableStatement) error {
	t := dbTable{}
	mb.tables[crt.Name.Value] = &t
	if crt.Columns == nil {

		return nil
	}

	for _, col := range *crt.Columns {
		t.columns = append(t.columns, col.Name.Value)

		var dt ColumnType
		switch col.Datatype.Value {
		case "int":
			dt = IntType
		case "text":
			dt = TextType
		case "bool":
			dt = BoolType
		default:
			return ErrInvalidDatatype
		}

		t.columnTypes = append(t.columnTypes, dt)
	}

	return nil
}

func (mb *MemoryBackend) Insert(inst *InsertStatement) error {
	table, ok := mb.tables[inst.Table.Value]
	if !ok {
		return ErrTableDoesNotExist
	}

	if inst.Values == nil {
		return nil
	}

	var row []MemoryCell

	if len(*inst.Values) != len(table.columns) {
		return ErrMissingValues
	}

	for _, value := range *inst.Values {
		if value.ExpressionType != LiteralType {
			fmt.Println("Skipping non-literal.")
			continue
		}

		row = append(row, mb.tokenToCell(value.Literal))
	}

	table.rows = append(table.rows, row)
	return nil
}

func (mb *MemoryBackend) tokenToCell(t *Token) MemoryCell {
	if t.Type == NumericType {
		buf := new(bytes.Buffer)
		i, err := strconv.Atoi(t.Value)
		if err != nil {
			panic(err)
		}

		err = binary.Write(buf, binary.BigEndian, int32(i))
		if err != nil {
			panic(err)
		}
		return buf.Bytes()
	}

	if t.Type == StringType {
		return MemoryCell(t.Value)
	}

	if t.Type == BooleanType {
		return literalToMemoryCell(t)
	}

	return nil
}

func (mb *MemoryBackend) Select(selectStatement *SelectStatement) (*Results, error) {
	table, ok := mb.tables[selectStatement.From.Value]
	if !ok {
		return nil, ErrTableDoesNotExist
	}

	var results [][]Cell
	var columns []struct {
		Type ColumnType
		Name string
	}

	for i, row := range table.rows {
		var result []Cell
		isFirstRow := i == 0

		for _, exp := range selectStatement.Items {
			if exp.ExpressionType != LiteralType {
				// Unsupported, doesn't currently exist, ignore.
				fmt.Println("Skipping non-literal expression.")
				continue
			}

			lit := exp.Literal

			if isSelectFromAllExpression(exp) {
				return selectStar(table)
			}

			if lit.Type == IdentifierType {
				found := false
				for i, tableCol := range table.columns {
					if tableCol == lit.Value {
						if isFirstRow {
							columns = append(columns, struct {
								Type ColumnType
								Name string
							}{
								Type: table.columnTypes[i],
								Name: table.columns[i],
							})
						}

						result = append(result, row[i])
						found = true
						break
					}
				}

				if !found {
					return nil, ErrColumnDoesNotExist
				}

				continue
			}

			return nil, ErrColumnDoesNotExist
		}

		results = append(results, result)
	}

	return &Results{
		Columns: columns,
		Rows:    results,
	}, nil
}

func isSelectFromAllExpression(exp *Expression) bool {
	return exp.Literal.Type == SymbolType && exp.Literal.Value == TokenFromSymbol(AsteriskSymbol).Value
}

func selectStar(table *dbTable) (*Results, error) {
	var results [][]Cell
	var columns []struct {
		Type ColumnType
		Name string
	}

	for i, row := range table.rows {
		var result []Cell
		isFirstRow := i == 0

		for i, tableCol := range table.columns {
			if isFirstRow {
				columns = append(columns, struct {
					Type ColumnType
					Name string
				}{
					Type: table.columnTypes[i],
					Name: tableCol,
				})
			}

			result = append(result, row[i])
		}

		results = append(results, result)
	}

	return &Results{
		Columns: columns,
		Rows:    results,
	}, nil
}
