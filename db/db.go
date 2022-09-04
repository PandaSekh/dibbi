package db

import "dibbi/data_structures"

type Db interface {
	Get(string) (interface{}, bool)
	Set(string, interface{}) bool
	Remove(string) bool
	Contains(string) bool
	GetHost() string
	GetPort() string
	GetTable() data_structures.HashTable
}
