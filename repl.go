package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

func startRepl() {
	mb := NewMemoryBackend()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("go_dibbi started.")
	for {
		fmt.Print("# ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}

		text = strings.Replace(text, "\n", "", -1)

		ast, err := parse(text)
		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, stmt := range ast.Statements {
			switch stmt.Type {
			case CreateTableType:
				err = mb.CreateTable(ast.Statements[0].CreateTableStatement)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("ok")
			case InsertType:
				err = mb.Insert(stmt.InsertStatement)
				if err != nil {
					fmt.Println(err)
					continue
				}

				fmt.Println("ok")
			case SelectType:
				results, err := mb.Select(stmt.SelectStatement)
				if err != nil {
					fmt.Println(err)
					continue
				}

				printTable(results)

				fmt.Println("ok")
			}
		}
	}
}

func printTable(results *Results) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	header := table.Row{}
	for _, col := range results.Columns {
		header = append(header, col.Name)
	}
	t.AppendHeader(header)

	for _, result := range results.Rows {
		row := table.Row{}
		for i, cell := range result {
			typ := results.Columns[i].Type
			r := ""
			switch typ {
			case IntType:
				i := cell.AsInt()
				if i != 0 {
					r = fmt.Sprintf("%d", i)
				}
			case TextType:
				s := cell.AsText()
				if s != "" {
					r = s
				}
			}

			row = append(row, r)
		}

		t.AppendRow(row)
	}

	t.SetAutoIndex(true)
	t.SetCaption(fmt.Sprintf("\nResults: %d\n", len(results.Rows)))
	t.SetStyle(table.StyleColoredYellowWhiteOnBlack)
	t.Render()
}
