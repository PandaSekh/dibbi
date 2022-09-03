package db

type Db interface {
	Get(string) (interface{}, bool)
	Set(string, interface{}) bool
	Remove(string) bool
	Contains(string) bool
}
