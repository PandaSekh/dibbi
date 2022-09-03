package dibbi_kv

import (
	"testing"
)

func TestDibbiKv_Contains(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "SET",
			args: args{key: "my_key", value: "value"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDibbiKv[string]()
			d.Set(tt.args.key, tt.args.value)
			if got := d.Contains(tt.args.key); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func TestDibbiKv_Get(t *testing.T) {
//	type fields struct {
//		table map[string]V
//	}
//	type args struct {
//		key string
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//		want   V
//		want1  bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			d := &DibbiKv{
//				table: tt.fields.table,
//			}
//			got, got1 := d.Get(tt.args.key)
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Get() got = %v, want %v", got, tt.want)
//			}
//			if got1 != tt.want1 {
//				t.Errorf("Get() got1 = %v, want %v", got1, tt.want1)
//			}
//		})
//	}
//}
//
//func TestDibbiKv_Remove(t *testing.T) {
//	type fields struct {
//		table map[string]V
//	}
//	type args struct {
//		key string
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//		want   bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			d := &DibbiKv{
//				table: tt.fields.table,
//			}
//			if got := d.Remove(tt.args.key); got != tt.want {
//				t.Errorf("Remove() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestDibbiKv_Set(t *testing.T) {
//	type fields struct {
//		table map[string]V
//	}
//	type args struct {
//		key   string
//		value V
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//		want   bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			d := &DibbiKv{
//				table: tt.fields.table,
//			}
//			if got := d.Set(tt.args.key, tt.args.value); got != tt.want {
//				t.Errorf("Set() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

///////////////////////////
/// Benchmarks
///////////////////////////
func BenchmarkSetString(b *testing.B) {
	kv := NewDibbiKv[int]()
	for i := 0; i < b.N; i++ {
		kv.Set(string(rune(i)), i)
	}
}
func BenchmarkSetStringHashTable(b *testing.B) {
	kv := NewDibbiKvHashTable()
	for i := 0; i < b.N; i++ {
		kv.HSet(string(rune(i)), i)
	}
}

func BenchmarkGetString(b *testing.B) {
	kv := NewDibbiKv[int]()
	// insert 1mln values
	for i := 0; i < 1_000_000; i++ {
		kv.Set(string(rune(i)), i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		kv.Get(string(rune(i)))
	}
}

func BenchmarkGetStringHashTable(b *testing.B) {
	kv := NewDibbiKvHashTable()
	// insert 1mln values
	for i := 0; i < 1_000_000; i++ {
		kv.HSet(string(rune(i)), i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		kv.HGet(string(rune(i)))
	}
}
