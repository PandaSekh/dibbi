package dibbi_kv

import (
	"dibbi/data_structures"
	"fmt"
)

const (
	defaultSize = 3
)

type DibbiKvCluster struct {
	size    int
	buckets []DibbiKv
}

func (c *DibbiKvCluster) getBucketIndex(key string) int {
	return int(data_structures.FnvHash(key) % uint64(c.size))
}

// NewSized generates a Cluster with the provided buckets size
func NewSized(initialSize int) *DibbiKvCluster {
	c := &DibbiKvCluster{
		size:    initialSize,
		buckets: make([]DibbiKv, initialSize),
	}

	for i := range c.buckets {
		c.buckets[i] = *NewDibbiKv()
	}
	return c
}

func (c *DibbiKvCluster) String() string {
	return fmt.Sprintf("Size: %d - Buckets: %v", c.size, c.buckets)
}

// New generates a Cluster with the default size for buckets (3)
func New() *DibbiKvCluster {
	return NewSized(defaultSize)
}

func (c *DibbiKvCluster) Get(key string) (interface{}, bool) {
	i := c.getBucketIndex(key)
	channel := make(chan interface{}, 1)
	go c.buckets[i].GetAsync(key, channel)
	val, open := <-channel
	if !open && val == nil {
		return nil, false
	}

	return val, true
}

func (c *DibbiKvCluster) Contains(key string) bool {
	i := c.getBucketIndex(key)
	return c.buckets[i].Contains(key)
}

func (c *DibbiKvCluster) Set(key string, value interface{}) bool {
	i := c.getBucketIndex(key)

	return c.buckets[i].Set(key, value)
}

func (c *DibbiKvCluster) Remove(key string) bool {
	i := c.getBucketIndex(key)
	return c.buckets[i].Remove(key)
}
