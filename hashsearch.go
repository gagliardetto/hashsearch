package hashsearch

import (
	"hash"
	"hash/fnv"
	"sync"
)

var hasherPool *sync.Pool

func init() {
	hasherPool = &sync.Pool{
		New: func() interface{} {
			return fnv.New64a()
		},
	}
}

// HashString returns a new `hash.Hash64` sum
// of the provided string.
func HashString(s string) uint64 {
	return HashBytes([]byte(s))
}

// HashBytes returns a new `hash.Hash64` sum
// of the provided []byte.
func HashBytes(b []byte) uint64 {
	h := hasherPool.Get().(hash.Hash64)
	defer hasherPool.Put(h)

	h.Reset()
	_, err := h.Write(b)
	if err != nil {
		panic(err)
	}

	return h.Sum64()
}

// HashSearch is a utility to track hashes of string and []byte values.
type HashSearch struct {
	hashes *Uint64
}

func New() *HashSearch {
	return &HashSearch{
		hashes: NewUint64(),
	}
}

// Sort sorts the index of hashes.
// This should be called after using `WARNING_Unsorted*`.
func (hs *HashSearch) Sort() {
	hs.hashes.Sort()
}

// Has returns true if the index contains the hash of the provided string.
func (hs *HashSearch) Has(item string) bool {
	return hs.hashes.Has(HashString(item))
}

// HasBytes returns true if the index contains the hash of the provided []byte.
func (hs *HashSearch) HasBytes(item []byte) bool {
	return hs.hashes.Has(HashBytes(item))
}

// HasOrAddBytes checks whether the hash of the provided item is already present;
// If the item is already present, the function returns true.
// If it is not present, the hash of the item is added
// and the function returns false.
func (hs *HashSearch) HasOrAddBytes(item []byte) bool {
	return hs.hashes.HasOrAdd(HashBytes(item))
}

// HasOrAdd checks whether the hash of the item is already present;
// If it is not present, the hash of the item is added
// and the function returns false.
// If the item is already present, the function returns true.
func (hs *HashSearch) HasOrAdd(item string) bool {
	return hs.hashes.HasOrAdd(HashString(item))
}

// Add add the hash of the provided value to the index.
func (hs *HashSearch) Add(item string) {
	hs.hashes.Add(HashString(item))
}

// AddFromBytes add the hash of the provided value to the index.
func (hs *HashSearch) AddFromBytes(item []byte) {
	hs.hashes.Add(HashBytes(item))
}

// WARNING_UnsortedAppend appends (without order considerations)
// the hash of the provided value to the index.
// WARNING: this function shouldn't be normally used.
// If this function is used, then you should call the `Sort()` method
// so that the index will be able to work as expected.
func (hs *HashSearch) WARNING_UnsortedAppend(item string) {
	hs.hashes.WARNING_UnsortedAppend(HashString(item))
}

// WARNING_UnsortedAppendBytes appends (without order considerations)
// the hash of the provided value to the index.
// WARNING: this function shouldn't be normally used.
// If this function is used, then you should call the `Sort()` method
// so that the index will be able to work as expected.
func (hs *HashSearch) WARNING_UnsortedAppendBytes(item []byte) {
	hs.hashes.WARNING_UnsortedAppend(HashBytes(item))
}
