package dibbi_kv

import (
	"dibbi/data_structures"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

// DibbiKv is a key-value database
type DibbiKv struct {
	table data_structures.HashTable
	mu    *sync.Mutex
	host  string
	port  string
}

func (d *DibbiKv) String() string {
	return fmt.Sprintf("%v", d.table)
}

func (d *DibbiKv) GetHost() string {
	return d.host
}

func (d *DibbiKv) GetTable() data_structures.HashTable {
	return d.table
}

func (d *DibbiKv) GetPort() string {
	return d.port
}

func NewDibbiKv() *DibbiKv {
	return NewDibbiKvRemote("localhost", strconv.Itoa(1111+rand.Intn(9999-1111)))
}

func NewDibbiKvRemote(host string, port string) *DibbiKv {
	dkv := &DibbiKv{
		table: *data_structures.NewSized(4000),
		mu:    &sync.Mutex{},
		host:  host,
		port:  port,
	}

	go dkv.startServer()

	return dkv
}

func (d *DibbiKv) startServer() {
	fmt.Printf("DibbiKv Server Start on: %s:%s\n", d.host, d.port)
	server, err := net.Listen("tcp", d.host+":"+d.port)
	if err != nil {
		fmt.Println("Error listening: ", err.Error())
		os.Exit(1)
	}
	defer func(server net.Listener) {
		_ = server.Close()
	}(server)
	for {
		connection, err := server.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		go d.processClientRequest(connection)
	}
}

func (d *DibbiKv) processClientRequest(connection net.Conn) {
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	req := string(buffer[:mLen])
	var res string

	inputs := strings.Split(req, " ")
	switch inputs[0] {
	case "get":
		if r, ok := d.Get(inputs[1]); ok {
			res = fmt.Sprintf("%v", r)
		} else {
			res = "nil"
		}
	case "set":
		d.Set(inputs[1], inputs[2])
		res = "ok"
	case "rem":
		d.Remove(inputs[1])
		res = "ok"
	}

	_, err = connection.Write([]byte(res))
	_ = connection.Close()
}

func (d *DibbiKv) GetAsync(key string, c chan interface{}) {
	d.mu.Lock()
	defer d.mu.Unlock()
	v, found := d.table.Get(key)

	if !found {
		close(c)
	} else {
		c <- v
	}
}

func (d *DibbiKv) Get(key string) (interface{}, bool) {
	d.mu.Lock()
	v, found := d.table.Get(key)
	d.mu.Unlock()

	return v, found
}

func (d *DibbiKv) Contains(key string) bool {
	_, found := d.Get(key)
	return found
}

func (d *DibbiKv) Set(key string, value interface{}) bool {
	d.mu.Lock()
	d.table.Set(key, value)
	d.mu.Unlock()

	return true
}

func (d *DibbiKv) Remove(key string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	res := d.table.Remove(key)

	return res
}

func (d *DibbiKv) RemoveAsync(key string, c chan bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	res := d.table.Remove(key)

	c <- res
}
