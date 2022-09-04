package dibbi_kv

import (
	"dibbi/data_structures"
	"dibbi/db"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	defaultSize = 3
)

type DibbiKvCluster struct {
	size    int
	buckets []db.Db
	Host    string
	Port    string
}

func (c *DibbiKvCluster) getBucketIndex(key string) int {
	return int(data_structures.FnvHash(key) % uint64(c.size))
}

// NewSized generates a Cluster with the provided buckets size
func NewSized(initialSize int) *DibbiKvCluster {
	c := &DibbiKvCluster{
		size:    initialSize,
		buckets: make([]db.Db, initialSize),
		Host:    "localhost",
		Port:    strconv.Itoa(1111 + rand.Intn(9999-1111)),
	}

	for i := range c.buckets {
		c.buckets[i] = NewDibbiKv()
	}

	c.startClusterServer()

	return c
}

func NewEmpty() *DibbiKvCluster {
	c := &DibbiKvCluster{
		size:    0,
		buckets: make([]db.Db, 0),
		Host:    "localhost",
		Port:    strconv.Itoa(1111 + rand.Intn(9999-1111)),
	}

	for i := range c.buckets {
		c.buckets[i] = NewDibbiKv()
	}

	go c.startClusterServer()

	return c
}

func (c *DibbiKvCluster) AddServerToCluster(host string, port string) error {
	newTable := make([]db.Db, len(c.buckets)+1)
	newTable = append(newTable, c.buckets...)
	newDb := NewDibbiKvRemote(host, port)
	newTable = append(newTable, newDb)

	// todo
	//for _, bucket := range c.buckets {
	//	for _, e := range bucket.GetTable() {
	//		newIndex := c.getBucketIndex()
	//		newHash := hashKey(e.key, len(newTable))
	//		newTable[newHash] = append(newTable[newHash], HashTableEntry{e.key, e.value})
	//	}
	//}
	c.buckets = newTable
	c.size += 1
	return nil
}

func printLocalIp() {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println(err)

	}

	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)
	ipAddress := conn.LocalAddr().(*net.UDPAddr)
	fmt.Printf("Local IP is: %s\n", ipAddress)
}

func (c *DibbiKvCluster) startClusterServer() {
	fmt.Printf("DibbiKv Cluster Server Start on: %s:%s\n", c.Host, c.Port)
	server, err := net.Listen("tcp", c.Host+":"+c.Port)
	if err != nil {
		fmt.Println("Error listening: ", err.Error())
		os.Exit(1)
	}
	printLocalIp()
	defer func(server net.Listener) {
		_ = server.Close()
	}(server)
	for {
		connection, err := server.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		go c.processClusterClientRequest(connection)
	}
}

func sendRequestToDibbiKvServer(host string, port string, msg string) string {
	//establish connection
	connection, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		panic(err)
	}
	///send some data
	_, err = connection.Write([]byte(msg))
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	defer func(connection net.Conn) {
		_ = connection.Close()
	}(connection)

	return string(buffer[:mLen])
}

func (c *DibbiKvCluster) processClusterClientRequest(connection net.Conn) {
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	msg := string(buffer[:mLen])
	// get correct key
	key := strings.Split(msg, " ")[1] // todo not hardcoded -> use encoder/decoder

	// get server to send message to
	bucketIndex := c.getBucketIndex(key)
	dibbiKvServer := c.buckets[bucketIndex]

	// send msg to correct server
	dibbiResponse := sendRequestToDibbiKvServer(dibbiKvServer.GetHost(), dibbiKvServer.GetPort(), msg)

	// return response
	_, err = connection.Write([]byte(dibbiResponse))
	_ = connection.Close()
}

func (c *DibbiKvCluster) String() string {
	return fmt.Sprintf("Size: %d - Buckets: %v", c.size, c.buckets)
}

// New generates a Cluster with the default size for buckets (3)
func New() *DibbiKvCluster {
	return NewSized(defaultSize)
}

//func (c *DibbiKvCluster) Get(key string) (interface{}, bool) {
//	i := c.getBucketIndex(key)
//	channel := make(chan interface{}, 1)
//	go c.buckets[i].GetAsync(key, channel)
//	val, open := <-channel
//	if !open && val == nil {
//		return nil, false
//	}
//
//	return val, true
//}

func (c *DibbiKvCluster) Contains(key string) bool {
	i := c.getBucketIndex(key)
	return c.buckets[i].Contains(key)
}

func (c *DibbiKvCluster) Set(key string, value interface{}) {
	i := c.getBucketIndex(key)
	go c.buckets[i].Set(key, value)
}

//func (c *DibbiKvCluster) Remove(key string) bool {
//	i := c.getBucketIndex(key)
//	channel := make(chan bool, 1)
//	go c.buckets[i].RemoveAsync(key, channel)
//
//	return <-channel
//}
