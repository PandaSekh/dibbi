package data_structures

// https://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function#FNV_hash_parameters
const (
	fnvOffsetBasis uint64 = 14695981039346656037
	fnvPrime              = 1099511628211
)

func fnvHash(key string) int {
	hash := fnvOffsetBasis
	sBytes := []byte(key)
	for _, b := range sBytes {
		hash = hash ^ uint64(b)
		hash = hash * fnvPrime
	}
	return int(hash)
}
