package main

import (
	"bufio"
	"dibbi/dibbi_kv"
	"fmt"
	"net"
	"os"
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
	cluster := dibbi_kv.NewEmpty()
	fmt.Println("Repl started.")
	for {
		fmt.Print("# ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}

		text = strings.Replace(text, "\n", "", -1)

		if text == "start server" {
			cluster = dibbi_kv.NewSized(5)
			fmt.Println("Server started")
		} else if strings.Contains(text, "add server") {
			inputs := strings.Split(text, " ")
			_ = cluster.AddServerToCluster(inputs[2], inputs[3])
		} else {
			channel := make(chan string, 1)
			go callServer(cluster, text, channel)
			value := <-channel
			fmt.Println(value)
		}
	}
}

func callServer(c *dibbi_kv.DibbiKvCluster, input string, channel chan string) {
	//establish connection
	connection, err := net.Dial("tcp", c.Host+":"+c.Port)
	if err != nil {
		panic(err)
	}
	///send some data
	_, err = connection.Write([]byte(input))
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	channel <- string(buffer[:mLen])
	defer func(connection net.Conn) {
		_ = connection.Close()
	}(connection)
}

//func processInput(input string, cluster dibbi_kv.DibbiKvCluster) {
//	inputs := strings.Split(input, " ")
//	switch inputs[0] {
//	case "get":
//		if r, ok := cluster.Get(inputs[1]); ok {
//			fmt.Printf("GET: %v\n", r)
//		}
//	case "set":
//		cluster.Set(inputs[1], inputs[2])
//	case "rem":
//		cluster.Remove(inputs[1])
//	}
//}
//
//func loadMigrations(mb *MemoryBackend) {
//	files, err := ioutil.ReadDir("./migrations/")
//	if err != nil {
//		fmt.Println("No migration found.")
//	}
//
//	for _, file := range files {
//		filepath := path.Join("migrations", file.Name())
//		dat, err := os.ReadFile(filepath)
//		if err != nil {
//			fmt.Printf("Reading migration file failed: %s\n", filepath)
//		}
//		ProcessInput(string(dat), mb)
//		fmt.Printf("Migration applied: %s", file.Name())
//	}
//}
