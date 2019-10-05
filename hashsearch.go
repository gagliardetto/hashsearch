package hashsearch

import (
	"hash"
	"hash/fnv"
	"sort"
	"sync"
)

type IntArr struct {
	mu  *sync.RWMutex
	arr []int
}

func NewIntArr() *IntArr {
	return &IntArr{
		mu: &sync.RWMutex{},
	}
}

func (ia *IntArr) Sort() {
	ia.mu.Lock()
	sort.Ints(ia.arr)
	ia.mu.Unlock()
}

func (ia *IntArr) OrderedAppend(hash int) {
	ia.mu.Lock()
	i := sort.SearchInts(ia.arr, hash)
	ia.arr = append(ia.arr, 0 /* use the zero value of the element type */)
	copy(ia.arr[i+1:], ia.arr[i:])
	ia.arr[i] = hash
	ia.mu.Unlock()
}
func (ia *IntArr) WarningUnorderedAppend(hash int) {
	ia.mu.Lock()
	ia.arr = append(ia.arr, hash)
	ia.mu.Unlock()
}
func (ia *IntArr) Has(hash int) bool {
	ia.mu.RLock()
	index := sort.SearchInts(ia.arr, int(hash))
	if index == len(ia.arr) {
		ia.mu.RUnlock()
		return false
	}
	has := int(hash) == ia.arr[index]
	ia.mu.RUnlock()
	return has
}

//////////////

type HashSearch struct {
	mu     *sync.RWMutex
	hashes []int
}

func New() *HashSearch {
	return &HashSearch{
		mu: &sync.RWMutex{},
	}
}
func (hs *HashSearch) WarningUnorderedAppend(item string) {
	hs.WarningUnorderedAppendBytes([]byte(item))
}
func (hs *HashSearch) WarningUnorderedAppendBytes(item []byte) {
	hs.mu.Lock()
	hs.hashes = append(hs.hashes, int(HashBytes(item)))
	hs.mu.Unlock()
}
func (hs *HashSearch) Sort() {
	hs.mu.Lock()
	sort.Ints(hs.hashes)
	hs.mu.Unlock()
}

func (hs *HashSearch) OrderedAppend(item string) {
	hs.OrderedAppendBytes([]byte(item))
}
func (hs *HashSearch) OrderedAppendBytes(item []byte) {
	numericHash := HashBytes(item)

	hs.mu.Lock()
	i := sort.SearchInts(hs.hashes, int(numericHash))
	hs.hashes = append(hs.hashes, 0 /* use the zero value of the element type */)
	copy(hs.hashes[i+1:], hs.hashes[i:])
	hs.hashes[i] = int(numericHash)
	hs.mu.Unlock()
}
func (hs *HashSearch) Has(item string) bool {
	return hs.HasBytes([]byte(item))
}
func (hs *HashSearch) HasBytes(item []byte) bool {
	numericHash := HashBytes(item)

	hs.mu.RLock()
	index := sort.SearchInts(hs.hashes, int(numericHash))
	if index == len(hs.hashes) {
		hs.mu.RUnlock()
		return false
	}
	has := int(numericHash) == hs.hashes[index]
	hs.mu.RUnlock()
	return has
}

var hasherPool *sync.Pool

func init() {
	hasherPool = &sync.Pool{
		New: func() interface{} {
			return fnv.New64a()
		},
	}
}
func HashString(s string) uint64 {
	return HashBytes([]byte(s))
}
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
