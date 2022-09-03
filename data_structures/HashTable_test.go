package data_structures

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"strconv"
	"testing"
)

func TestHashTable_Get(t *testing.T) {
	type fields struct {
		size    int
		buckets [][]HashTableEntry
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
		want1  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ht := &HashTable{
				size:    tt.fields.size,
				buckets: tt.fields.buckets,
			}
			got, got1 := ht.Get(tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestHashTable_Set(t *testing.T) {
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			args: args{
				key:   "key",
				value: 1920,
			},
			name: "integer",
		},
		{
			args: args{
				key:   "key_2",
				value: 1922331231210,
			},
			name: "big integer",
		},
		{
			args: args{
				key:   "a_long_key_value",
				value: "hello",
			},
			name: "long key",
		},
		{
			args: args{
				key:   "a long key value",
				value: []int{1, 2, 3, 4, 5},
			},
			name: "long key with space and slice as value",
		},
	}
	size := 5
	ht := NewSized(size)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ht.Set(tt.args.key, tt.args.value)
		})
	}
	assert.Equal(t, ht.size, size+len(tests))
}

//
//func TestHashTable_expandTable(t *testing.T) {
//	type fields struct {
//		size    int
//		buckets [][]HashTableEntry
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ht := &HashTable{
//				size:    tt.fields.size,
//				buckets: tt.fields.buckets,
//			}
//			if err := ht.expandTable(); (err != nil) != tt.wantErr {
//				t.Errorf("expandTable() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestStringKey_Equal(t *testing.T) {
//	type args struct {
//		other StringKey
//	}
//	tests := []struct {
//		name string
//		str  StringKey
//		args args
//		want bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := tt.str.Equal(tt.args.other); got != tt.want {
//				t.Errorf("Equal() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestStringKey_HashBytes(t *testing.T) {
//	tests := []struct {
//		name string
//		str  StringKey
//		want []byte
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := tt.str.HashBytes(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("HashBytes() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func Test_hashKey(t *testing.T) {
//	type args struct {
//		key   StringKey
//		limit int
//	}
//	tests := []struct {
//		name string
//		args args
//		want int
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := hashKey(tt.args.key, tt.args.limit); got != tt.want {
//				t.Errorf("hashKey() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

// Benchmarks
//func BenchmarkStringStringTableSize(b *testing.B) {
//	benchmarks := []struct {
//		name      string
//		tableSize int
//	}{
//		{
//			name:      "16",
//			tableSize: 16,
//		},
//		{
//			name:      "32",
//			tableSize: 32,
//		},
//		{
//			name:      "64",
//			tableSize: 64,
//		},
//		{
//			name:      "128",
//			tableSize: 128,
//		},
//		{
//			name:      "256",
//			tableSize: 256,
//		},
//		{
//			name:      "512",
//			tableSize: 512,
//		},
//		{
//			name:      "1024",
//			tableSize: 1024,
//		},
//		{
//			name:      "2048",
//			tableSize: 2048,
//		},
//		//{
//		//	name:      "4096",
//		//	tableSize: 4096,
//		//},
//		//{
//		//	name:      "8192",
//		//	tableSize: 8192,
//		//},
//		//{
//		//	name:      "16384",
//		//	tableSize: 16384,
//		//},
//	}
//	for _, bm := range benchmarks {
//		hTable := NewSized(bm.tableSize)
//		b.Run(bm.name, func(b *testing.B) {
//			for i := 0; i < b.N; i++ {
//				s := strconv.Itoa(i)
//				hTable.Set(s, s)
//			}
//		})
//	}
//}
//
func BenchmarkSetStringInteger(b *testing.B) {
	hTable := New()
	for i := 0; i < b.N; i++ {
		hTable.Set(strconv.Itoa(i), i)
	}
}

func BenchmarkSetStringConverterString(b *testing.B) {
	hTable := New()
	for i := 0; i < b.N; i++ {
		hTable.Set(strconv.Itoa(i), "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
}
