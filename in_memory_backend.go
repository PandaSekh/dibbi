package dibbi

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"
	"strconv"
)

var (
	trueToken  = token{tokenType: booleanType, value: "true"}
	falseToken = token{tokenType: booleanType, value: "false"}

	trueMemoryCell  = literalToMemoryCell(&trueToken)
	falseMemoryCell = literalToMemoryCell(&falseToken)
	nullMemoryCell  = literalToMemoryCell(&token{tokenType: nullType})
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

// literalToMemoryCell maps a Go value into a memory Cell
func literalToMemoryCell(t *token) memoryCell {
	if t.tokenType == numericType {
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
		tokenType == stringType {
		return memoryCell(t.value)
	}

	if t.
		tokenType == booleanType {
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

func (mb *InMemoryBackend) Select(selectStatement *selectStatement) (*Results, error) {
	table := &table{}

	if selectStatement.from != nil && selectStatement.from.value != "" {
		var ok bool
		table, ok = mb.tables[selectStatement.from.value]
		if !ok {
			return nil, ErrTableDoesNotExist
		}
	}

	if selectStatement.items == nil || len(*selectStatement.items) == 0 {
		return &Results{}, nil
	}

	var results [][]Cell
	var columns []ResultColumn

	if selectStatement.from == nil {
		table.rows = [][]memoryCell{{}}
	}

	for row := range table.rows {
		var result []Cell
		isFirstRow := len(results) == 0

		if selectStatement.where != nil {
			value, _, _, err := table.evaluateCell(uint(row), *selectStatement.where)
			if err != nil {
				return nil, err
			}

			if !*value.AsBool() {
				continue
			}
		}

		for _, column := range *selectStatement.items {
			if column.asterisk {
				// TODO improve
				return selectStar(table)
			}
			value, columnName, columnType, err := table.evaluateCell(uint(row), *column.exp)
			if err != nil {
				return nil, err
			}

			if isFirstRow {
				columns = append(columns, ResultColumn{
					Type: columnType,
					Name: columnName,
				})
			}

			result = append(result, value)
		}

		results = append(results, result)
	}

	return &Results{
		Columns: columns,
		Rows:    results,
	}, nil
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

func tokenToCell(t *token) memoryCell {
	if t.
		tokenType == numericType {
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
		tokenType == stringType {
		return memoryCell(t.value)
	}

	if t.
		tokenType == booleanType {
		return literalToMemoryCell(t)
	}

	return nil
}

// evaluateCell evaluates the expression in the cell
func (t *table) evaluateCell(rowIndex uint, exp expression) (memoryCell, string, ColumnType, error) {
	switch exp.expressionType {
	case literalType:
		return t.evaluateLiteralCell(rowIndex, exp)
	case binaryType:
		return t.evaluateBinaryCell(rowIndex, exp)
	default:
		return nil, "", 0, ErrInvalidCell
	}
}

func (t *table) evaluateLiteralCell(rowIndex uint, exp expression) (memoryCell, string, ColumnType, error) {
	if exp.expressionType != literalType {
		return nil, "", 0, ErrInvalidCell
	}

	lit := exp.literal
	if lit.tokenType == identifierType {
		for i, tableCol := range t.columns {
			if tableCol == lit.value {
				return t.rows[rowIndex][i], tableCol, t.columnTypes[i], nil
			}
		}

		return nil, "", 0, ErrColumnDoesNotExist
	}

	columnType := IntType
	if lit.tokenType == stringType {
		columnType = TextType
	} else if lit.tokenType == booleanType {
		columnType = BoolType
	}

	return literalToMemoryCell(lit), "?column?", columnType, nil
}

func (t *table) evaluateBinaryCell(rowIndex uint, exp expression) (memoryCell, string, ColumnType, error) {
	if exp.expressionType != binaryType {
		return nil, "", 0, ErrInvalidCell
	}

	binaryExp := exp.binary

	left, _, leftType, err := t.evaluateCell(rowIndex, binaryExp.left)
	if err != nil {
		return nil, "", 0, err
	}

	right, _, rightType, err := t.evaluateCell(rowIndex, binaryExp.right)
	if err != nil {
		return nil, "", 0, err
	}

	switch binaryExp.operator.tokenType {
	case symbolType:
		switch symbol(binaryExp.operator.value) {
		case equalsSymbol:
			areOperandsEquals := left.equals(right)
			if leftType == TextType && rightType == TextType && areOperandsEquals {
				return trueMemoryCell, "?column?", BoolType, nil
			}

			if leftType == IntType && rightType == IntType && areOperandsEquals {
				return trueMemoryCell, "?column?", BoolType, nil
			}

			if leftType == BoolType && rightType == BoolType && areOperandsEquals {
				return trueMemoryCell, "?column?", BoolType, nil
			}

			return falseMemoryCell, "?column?", BoolType, nil
		case notEqualSymbol:
			if leftType != rightType || !left.equals(right) {
				return trueMemoryCell, "?column?", BoolType, nil
			}

			return falseMemoryCell, "?column?", BoolType, nil
		case concatSymbol:
			if leftType != TextType || rightType != TextType {
				return nil, "", 0, ErrInvalidOperands
			}

			return literalToMemoryCell(&token{tokenType: stringType, value: *left.AsText() + *right.AsText()}), "?column?", TextType, nil
		case plusSymbol:
			if leftType != IntType || rightType != IntType {
				return nil, "", 0, ErrInvalidOperands
			}

			iValue := int(*(left.AsInt()) + *(right.AsInt()))
			return literalToMemoryCell(&token{tokenType: numericType, value: strconv.Itoa(iValue)}), "?column?", IntType, nil
		default:
			// TODO
			break
		}

	case keywordType:
		switch keyword(binaryExp.operator.value) {
		case andKeyword:
			if leftType != BoolType || rightType != BoolType {
				return nil, "", 0, ErrInvalidOperands
			}

			res := falseMemoryCell
			if *left.AsBool() && *right.AsBool() {
				res = trueMemoryCell
			}

			return res, "?column?", BoolType, nil
		case orKeyword:
			if leftType != BoolType || rightType != BoolType {
				return nil, "", 0, ErrInvalidOperands
			}

			res := falseMemoryCell
			if *left.AsBool() || *right.AsBool() {
				res = trueMemoryCell
			}

			return res, "?column?", BoolType, nil
		default:
			// TODO
			break
		}
	}

	return nil, "", 0, ErrInvalidCell
}
