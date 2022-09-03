package main

import (
	"bufio"
	"dibbi/data_structures"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func main() {
	//mb := NewMemoryBackend()
	//loadMigrations(mb)
	//startRepl(mb)
	startReplHTable()
}

func startReplHTable() {
	reader := bufio.NewReader(os.Stdin)
	hTable := data_structures.NewSized(3)
	fmt.Println("HTable repl started.")
	for {
		fmt.Print("# ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}

		text = strings.Replace(text, "\n", "", -1)
		processInput(text, *hTable)
	}
}

func processInput(input string, htable data_structures.HashTable) {
	inputs := strings.Split(input, " ")
	switch inputs[0] {
	case "get":
		if r, ok := htable.Get(inputs[1]); ok {
			fmt.Printf("GET: %v\n", r)
		}
	case "set":
		htable.Set(inputs[1], inputs[2])
	}
}

func loadMigrations(mb *MemoryBackend) {
	files, err := ioutil.ReadDir("./migrations/")
	if err != nil {
		fmt.Println("No migration found.")
	}

	for _, file := range files {
		filepath := path.Join("migrations", file.Name())
		dat, err := os.ReadFile(filepath)
		if err != nil {
			fmt.Printf("Reading migration file failed: %s\n", filepath)
		}
		ProcessInput(string(dat), mb)
		fmt.Printf("Migration applied: %s", file.Name())
	}
}
