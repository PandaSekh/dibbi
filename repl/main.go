package main

import "github.com/PandaSekh/dibbi"

func main() {
	mb := dibbi.NewMemoryBackend()

	migrate(mb)
	startRepl(mb)
}
