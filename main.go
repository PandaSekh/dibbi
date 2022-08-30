package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func main(){
	mb := NewMemoryBackend()
	loadMigrations(mb)
	startRepl(mb)
}

func loadMigrations(mb *MemoryBackend){
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