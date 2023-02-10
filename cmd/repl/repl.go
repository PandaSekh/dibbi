package repl

import (
	"bufio"
	"dibbi/internal"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

func StartRepl() {
	mb := internal.NewMemoryBackend()
	migrate(mb)
	start(mb)
}

func start(mb *internal.MemoryBackend) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("dibbi started.")
	for {
		fmt.Print("# ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}

		text = strings.Replace(text, "\n", "", -1)
		ProcessInput(text, mb)
	}
}

func ProcessInput(text string, mb *internal.MemoryBackend) {
	ast, err := internal.Parse(text)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, stmt := range ast.Statements {
		switch stmt.Type {
		case internal.CreateTableType:
			err = mb.CreateTable(ast.Statements[0].CreateTableStatement)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case internal.InsertType:
			err = mb.Insert(stmt.InsertStatement)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case internal.SelectType:
			results, err := mb.Select(stmt.SelectStatement)
			if err != nil {
				fmt.Println(err)
				continue
			}

			printTable(results)
		}
	}
}

func printTable(results *internal.Results) {
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
			case internal.IntType:
				i := cell.AsInt()
				if i != 0 {
					r = fmt.Sprintf("%d", i)
				}
			case internal.TextType:
				s := cell.AsText()
				if s != "" {
					r = s
				}
			case internal.BoolType:
				s := cell.AsBool()
				r = strconv.FormatBool(s)
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
