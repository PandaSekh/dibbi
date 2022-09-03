package main

import (
	"bufio"
	"dibbi/dibbi_kv"
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
	cluster := dibbi_kv.NewSized(5)
	fmt.Println("HTable repl started.")
	for {
		fmt.Print("# ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}

		text = strings.Replace(text, "\n", "", -1)
		processInput(text, *cluster)
	}
}

func processInput(input string, cluster dibbi_kv.DibbiKvCluster) {
	inputs := strings.Split(input, " ")
	switch inputs[0] {
	case "get":
		if r, ok := cluster.Get(inputs[1]); ok {
			fmt.Printf("GET: %v\n", r)
		}
	case "set":
		cluster.Set(inputs[1], inputs[2])
	case "rem":
		cluster.Remove(inputs[1])
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
