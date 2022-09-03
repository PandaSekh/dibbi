package data_structures

import (
	"fmt"
)

const (
	defaultSize         = 32
	loadFactorThreshold = 0.5
)

type HashTable struct {
	size    int
	buckets [][]HashTableEntry
}

func (ht *HashTable) String() string {
	return fmt.Sprintf("Size: %d, Buckets: %v", ht.size, ht.buckets)
}

type HashTableEntry struct {
	key   string
	value interface{}
}

// NewSized generates a HashTable with the provided buckets size
func NewSized(initialSize int) *HashTable {
	return &HashTable{
		size:    initialSize,
		buckets: make([][]HashTableEntry, initialSize),
	}
}

// New generates a HashTable with the default size for buckets (16)
func New() *HashTable {
	return NewSized(defaultSize)
}

// hashKey returns the hash of the provided StringKey capped by limit.
func hashKey(key string, limit int) int {
	return int(FnvHash(key) % uint64(limit))
}

// loadFactor returns the current loadFactor of the table
func (ht *HashTable) loadFactor() float32 {
	return float32(ht.size) / float32(len(ht.buckets))
}

// Get returns the value corresponding to the provided StringKey and true if found
func (ht *HashTable) Get(key string) (interface{}, bool) {
	hash := hashKey(key, len(ht.buckets))

	for _, value := range ht.buckets[hash] {
		if value.key == key {
			return value.value, true
		}
	}
	return nil, false
}

func (ht *HashTable) Set(key string, value interface{}) {
	hash := hashKey(key, len(ht.buckets))

	for i, el := range ht.buckets[hash] {
		if el.key == key {
			// if key is already present, overwrite
			ht.buckets[hash][i].value = value
			return
		}
	}

	ht.buckets[hash] = append(ht.buckets[hash], HashTableEntry{key, value})
	ht.size += 1
	if ht.loadFactor() > loadFactorThreshold {
		err := ht.expandTable()
		if err != nil {
			fmt.Printf("Error while setting key: %s value: %v. Error: %v", key, value, err)
		}
	}
}

func (ht *HashTable) Remove(key string) bool {
	hash := hashKey(key, len(ht.buckets))

	for index, value := range ht.buckets[hash] {
		if value.key == key {
			ret := make([]HashTableEntry, 0)
			ret = append(ret, ht.buckets[hash][:index]...)
			ht.buckets[hash] = ret
			ht.size -= 1
			return true
		}
	}
	return false
}

func (ht *HashTable) expandTable() error {
	newTable := make([][]HashTableEntry, len(ht.buckets)*2)
	for _, bucket := range ht.buckets {
		for _, e := range bucket {
			newHash := hashKey(e.key, len(newTable))
			newTable[newHash] = append(newTable[newHash], HashTableEntry{e.key, e.value})
		}
	}
	ht.buckets = newTable
	return nil
}
