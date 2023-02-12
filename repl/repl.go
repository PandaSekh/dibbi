package main

import (
	"bufio"
	"fmt"
	"github.com/PandaSekh/dibbi"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
	"strconv"
	"strings"
)

func startRepl(database dibbi.Database) {
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
		result, present, errQuery := dibbi.Query(text, &database)

		if errQuery != nil {
			fmt.Println(err)
			continue
		}

		if present && result != nil {
			printTable(result)
			continue
		}
	}
}

func printTable(results *dibbi.Results) {
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
			case dibbi.IntType:
				i := cell.AsInt()
				if *i != 0 {
					r = fmt.Sprintf("%d", i)
				}
			case dibbi.TextType:
				s := cell.AsText()
				if *s != "" {
					r = *s
				}
			case dibbi.BoolType:
				s := cell.AsBool()
				r = strconv.FormatBool(*s)
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
