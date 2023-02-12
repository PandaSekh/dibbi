package dibbi

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"
	"strconv"
)

type memoryCell []byte

func (mc memoryCell) AsInt() *int32 {
	if len(mc) == 0 {
		return nil
	}

	var i int32
	err := binary.Read(bytes.NewBuffer(mc), binary.BigEndian, &i)
	if err != nil {
		fmt.Printf("Corrupted data [%s]: %s\n", mc, err)
		return nil
	}

	return &i
}

func (mc memoryCell) AsText() *string {
	if len(mc) == 0 {
		return nil
	}

	s := string(mc)
	return &s
}

func (mc memoryCell) AsBool() *bool {
	b, err := strconv.ParseBool(string(mc))
	if err != nil {
		f := false
		return &f
	}
	return &b
}

func (mc memoryCell) equals(b memoryCell) bool {
	if mc == nil || b == nil {
		return mc == nil && b == nil
	}

	return bytes.Equal(mc, b)
}

// literalToMemoryCell maps left Go value into left memory Cell
func literalToMemoryCell(t *token) memoryCell {
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
		return memoryCell(t.value)
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

type table struct {
	name        string
	columns     []string
	columnTypes []ColumnType
	rows        [][]memoryCell
}

type InMemoryBackend struct {
	tables map[string]*table
}

func NewMemoryBackend() *InMemoryBackend {
	fmt.Println("Created new memory backend")
	return &InMemoryBackend{
		tables: map[string]*table{},
	}
}

func (mb *InMemoryBackend) CreateTable(crt *createTableStatement) error {
	t := table{}
	mb.tables[crt.Name.value] = &t
	if crt.Columns == nil {

		return nil
	}

	for _, col := range *crt.Columns {
		t.columns = append(t.columns, col.Name.value)

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

func (mb *InMemoryBackend) Insert(inst *insertStatement) error {
	table, ok := mb.tables[inst.Table.value]
	if !ok {
		return ErrTableDoesNotExist
	}

	if inst.Values == nil {
		return nil
	}

	var row []memoryCell

	if len(*inst.Values) != len(table.columns) {
		return ErrMissingValues
	}

	for _, value := range *inst.Values {
		if value.expressionType != literalType {
			fmt.Println("Skipping non-literal.")
			continue
		}

		row = append(row, tokenToCell(value.literal))
	}

	table.rows = append(table.rows, row)
	return nil
}

func tokenToCell(t *token) memoryCell {
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
		return memoryCell(t.value)
	}

	if t.
		tokenType == BooleanType {
		return literalToMemoryCell(t)
	}

	return nil
}

func (mb *InMemoryBackend) Select(selectStatement *selectStatement) (*Results, error) {
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
			if exp.expressionType != literalType {
				// Unsupported, doesn't currently exist, ignore.
				fmt.Println("Skipping non-literal expression.")
				continue
			}

			lit := exp.literal

			if isSelectFromAllExpression(exp) {
				return selectStar(table)
			}

			if lit.
				tokenType == IdentifierType {
				found := false
				for i, tableCol := range table.columns {
					if tableCol == lit.value {
						if isFirstRow {
							Columns = append(Columns, ResultColumn{
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
		Columns: Columns,
		Rows:    results,
	}, nil
}

func isSelectFromAllExpression(exp *expression) bool {
	return exp.literal.
		tokenType == symbolType && exp.literal.value == tokenFromSymbol(AsteriskSymbol).value
}

func selectStar(table *table) (*Results, error) {
	var results [][]Cell
	var Columns []ResultColumn

	for i, row := range table.rows {
		var result []Cell
		isFirstRow := i == 0

		for i, tableCol := range table.columns {
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

	return &Results{
		Columns: Columns,
		Rows:    results,
	}, nil
}
