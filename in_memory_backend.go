package dibbi

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"
	"strconv"
)

type MemoryCell []byte

func (mc MemoryCell) AsInt() *int32 {
	var i int32
	err := binary.Read(bytes.NewBuffer(mc), binary.BigEndian, &i)
	if err != nil {
		fmt.Printf("Corrupted data [%s]: %s\n", mc, err)
		return nil
	}

	return &i
}

func (mc MemoryCell) AsText() *string {
	if len(mc) == 0 {
		return nil
	}

	s := string(mc)
	return &s
}

func (mc MemoryCell) AsBool() *bool {
	if len(mc) == 0 {
		return nil
	}

	b := mc[0] == 1
	return &b
}

func (mc MemoryCell) equals(b MemoryCell) bool {
	if mc == nil || b == nil {
		return mc == nil && b == nil
	}

	return bytes.Equal(mc, b)
}

// literalToMemoryCell maps a Go value into a memory Cell
func literalToMemoryCell(t *token) MemoryCell {
	if t.tokenType == NumericType {
		buf := new(bytes.Buffer)
		i, err := strconv.Atoi(t.value)
		if err != nil {
			fmt.Printf("Corrupted data [%s]: %s\n", t.value, err)
			return nil
		}

		err = binary.Write(buf, binary.BigEndian, new(big.Int).SetInt64(int64(i)).Bytes())
		if err != nil {
			fmt.Printf("Corrupted data [%s]: %s\n", buf.String(), err)
			return nil
		}
		return buf.Bytes()
	}

	if t.
		tokenType == StringType {
		return MemoryCell(t.value)
	}

	if t.
		tokenType == BooleanType {
		if t.value == "true" {
			return []byte{1}
		}

		return []byte{0}
	}

	return nil
}

type dbTable struct {
	Columns     []string
	columnTypes []ColumnType
	rows        [][]MemoryCell
}

type memoryBackend struct {
	tables map[string]*dbTable
}

func newMemoryBackend() *memoryBackend {
	fmt.Println("Created new memory backend")
	return &memoryBackend{
		tables: map[string]*dbTable{},
	}
}

func (mb *memoryBackend) CreateTable(crt *createTableStatement) error {
	t := dbTable{}
	mb.tables[crt.Name.value] = &t
	if crt.Columns == nil {

		return nil
	}

	for _, col := range *crt.Columns {
		t.Columns = append(t.Columns, col.Name.value)

		var dt ColumnType
		switch col.Datatype.value {
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

func (mb *memoryBackend) Insert(inst *InsertStatement) error {
	table, ok := mb.tables[inst.Table.value]
	if !ok {
		return ErrTableDoesNotExist
	}

	if inst.Values == nil {
		return nil
	}

	var row []MemoryCell

	if len(*inst.Values) != len(table.Columns) {
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

func (mb *memoryBackend) tokenToCell(t *token) MemoryCell {
	if t.
		tokenType == NumericType {
		buf := new(bytes.Buffer)
		i, err := strconv.Atoi(t.value)
		if err != nil {
			panic(err)
		}

		err = binary.Write(buf, binary.BigEndian, int32(i))
		if err != nil {
			panic(err)
		}
		return buf.Bytes()
	}

	if t.
		tokenType == StringType {
		return MemoryCell(t.value)
	}

	if t.
		tokenType == BooleanType {
		return literalToMemoryCell(t)
	}

	return nil
}

func (mb *memoryBackend) Select(selectStatement *selectStatement) (*QueryResults, error) {
	table, ok := mb.tables[selectStatement.from.value]
	if !ok {
		return nil, ErrTableDoesNotExist
	}

	var results [][]Cell
	var Columns []ResultColumn

	for i, row := range table.rows {
		var result []Cell
		isFirstRow := i == 0

		for _, exp := range selectStatement.items {
			if exp.ExpressionType != LiteralType {
				// Unsupported, doesn't currently exist, ignore.
				fmt.Println("Skipping non-literal expression.")
				continue
			}

			lit := exp.Literal

			if isSelectFromAllExpression(exp) {
				return selectStar(table)
			}

			if lit.
				tokenType == IdentifierType {
				found := false
				for i, tableCol := range table.Columns {
					if tableCol == lit.value {
						if isFirstRow {
							Columns = append(Columns, ResultColumn{
								Type: table.columnTypes[i],
								Name: table.Columns[i],
							})
						}

						result = append(result, row[i])
						found = true
						break
					}
				}

				if !found {
					return nil, ErrcolumnDoesNotExist
				}

				continue
			}

			return nil, ErrcolumnDoesNotExist
		}

		results = append(results, result)
	}

	return &QueryResults{
		Columns: Columns,
		Rows:    results,
	}, nil
}

func isSelectFromAllExpression(exp *expression) bool {
	return exp.Literal.
		tokenType == SymbolType && exp.Literal.value == tokenFromSymbol(AsteriskSymbol).value
}

func selectStar(table *dbTable) (*QueryResults, error) {
	var results [][]Cell
	var Columns []ResultColumn

	for i, row := range table.rows {
		var result []Cell
		isFirstRow := i == 0

		for i, tableCol := range table.Columns {
			if isFirstRow {
				Columns = append(Columns, ResultColumn{
					Type: table.columnTypes[i],
					Name: tableCol,
				})
			}

			result = append(result, row[i])
		}

		results = append(results, result)
	}

	return &QueryResults{
		Columns: Columns,
		Rows:    results,
	}, nil
}
