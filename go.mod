module main

go 1.18

replace go_dibbi/lexer => ./lexer

replace go_dibbi/ast => ./ast

require (
	go_dibbi/ast v0.0.0-00010101000000-000000000000
	go_dibbi/lexer v0.0.0-00010101000000-000000000000
)
