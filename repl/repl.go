package main

import (
	"fmt"
	"github.com/PandaSekh/dibbi"
	"github.com/chzyer/readline"
	"github.com/jedib0t/go-pretty/v6/table"
	"io"
	"os"
	"strconv"
	"strings"
)

func startRepl(database dibbi.Database) {
	l, err := readline.NewEx(&readline.Config{
		Prompt:          "# ",
		HistoryFile:     "/tmp/tmp",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		panic(err)
	}
	defer func(l *readline.Instance) {
		err := l.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(l)

	fmt.Println("dibbi started.")

repl:
	for {
		fmt.Print("# ")
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue repl
			}
		} else if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error while reading line:", err)
			continue repl
		}

		trimmed := strings.TrimSpace(line)
		if trimmed == "quit" || trimmed == "exit" || trimmed == "\\q" {
			break
		}

		result, queryError := dibbi.Query(line, &database)

		if queryError != nil {
			fmt.Println(queryError)
			continue
		}

		if result != nil {
			printResults(result)
			continue
		}
	}
}

func printResults(results *dibbi.Results) {
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
