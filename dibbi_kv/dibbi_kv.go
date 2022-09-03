package dibbi_kv

import "dibbi/data_structures"

// DibbiKv is a key-value database
type DibbiKv[V interface{}] struct {
	table map[string]V
}

type DibbiKvHash struct {
	table data_structures.HashTable
}

func NewDibbiKv[T interface{}]() *DibbiKv[T] {
	c := DibbiKv[T]{}
	c.table = make(map[string]T)

	return &c
}

func NewDibbiKvHashTable() *DibbiKvHash {
	c := DibbiKvHash{}
	c.table = *data_structures.NewSized(4000)

	return &c
}

func (d *DibbiKv[V]) Get(key string) (V, bool) {
	v, found := d.table[key]
	return v, found
}

func (d *DibbiKvHash) HGet(key string) (interface{}, bool) {
	v, found := d.table.Get(key)
	return v, found
}

func (d *DibbiKv[V]) Contains(key string) bool {
	_, found := d.table[key]
	return found
}

func (d *DibbiKv[V]) Set(key string, value V) bool {
	d.table[key] = value
	return true
}

func (d *DibbiKvHash) HSet(key string, value interface{}) bool {
	d.table.Set(key, value)
	return true
}

func (d *DibbiKv[V]) Remove(key string) bool {
	delete(d.table, key)
	return true
}
