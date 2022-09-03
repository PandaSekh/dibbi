package dibbi_kv

import (
	"dibbi/data_structures"
	"fmt"
	"sync"
)

// DibbiKv is a key-value database
type DibbiKv struct {
	table data_structures.HashTable
	mu    sync.Mutex
}

func (d *DibbiKv) String() string {
	return fmt.Sprintf("%v", d.table)
}

func NewDibbiKv() *DibbiKv {
	c := DibbiKv{}
	c.table = *data_structures.NewSized(4000)

	return &c
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
