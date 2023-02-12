package main

import (
	"fmt"
	"github.com/PandaSekh/dibbi"
	"os"
	"path"
)

func migrate(db dibbi.Database) {
	files, err := os.ReadDir("./migrations/")
	if err != nil {
		fmt.Println("No migration found.")
	}

	for _, file := range files {
		filepath := path.Join("migrations", file.Name())
		dat, err := os.ReadFile(filepath)
		if err != nil {
			fmt.Printf("Reading migration file failed: %s\n", filepath)
		}
		dibbi.Query(string(dat), &db)
		fmt.Printf("Migration applied: %s\n", file.Name())
	}
}
