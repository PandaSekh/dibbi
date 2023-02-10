package repl

import (
	"dibbi/internal"
	"fmt"
	"os"
	"path"
)

func migrate(mb *internal.MemoryBackend) {
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
		ProcessInput(string(dat), mb)
		fmt.Printf("Migration applied: %s\n", file.Name())
	}
}
