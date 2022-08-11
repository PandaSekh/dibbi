package main

import (
	"fmt"
	"go_dibbi/lexer"
)

func main() {
	tokens, err := lexer.Lex("SELECT * FROM my_table WHERE name = 'hello_world'")
	//tokens, err := lexer.Lex("SELECT * FROM my_table WHERE name = hello_world")
	//tokens, err := lexer.Lex("SELECT a")
	if err != nil {
		fmt.Println(err)
	}

	for _, token := range tokens {
		fmt.Println(token.Value)
	}
}
