package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
)

type MemoryCell []byte

func (mc MemoryCell) AsInt() int32 {
    var i int32
    err := binary.Read(bytes.NewBuffer(mc), binary.BigEndian, &i)
    if err != nil {
        panic(err)
    }

    return i
}

func (mc MemoryCell) AsText() string {
    return string(mc)
}

type dbtable struct {
    columns     []string
    columnTypes []ColumnType
    rows        [][]MemoryCell
}

type MemoryBackend struct {
    tables map[string]*dbtable
}

func NewMemoryBackend() *MemoryBackend {
    return &MemoryBackend{
        tables: map[string]*dbtable{},
    }
}

func (mb *MemoryBackend) CreateTable(crt *CreateTableStatement) error {
    t := dbtable{}
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

    row := []MemoryCell{}

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
        return MemoryCell(buf.Bytes())
    }

    if t.Type == StringType {
        return MemoryCell(t.Value)
    }

    return nil
}

func (mb *MemoryBackend) Select(slct *SelectStatement) (*Results, error) {
    table, ok := mb.tables[slct.From.Value]
    if !ok {
        return nil, ErrTableDoesNotExist
    }

    results := [][]Cell{}
    columns := []struct {
        Type ColumnType
        Name string
    }{}

    for i, row := range table.rows {
        result := []Cell{}
        isFirstRow := i == 0

        for _, exp := range slct.Items {
            if exp.ExpressionType != LiteralType {
                // Unsupported, doesn't currently exist, ignore.
                fmt.Println("Skipping non-literal expression.")
                continue
            }

            lit := exp.Literal
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
                                Name: lit.Value,
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