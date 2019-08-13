package hashsearch

import (
	"hash"
	"hash/fnv"
	"sort"
	"sync"
)

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
	hs.mu.Lock()
	hs.hashes = append(hs.hashes, int(doHash(item)))
	hs.mu.Unlock()
}
func (hs *HashSearch) Sort() {
	hs.mu.Lock()
	sort.Ints(hs.hashes)
	hs.mu.Unlock()
}

func (hs *HashSearch) OrderedAppend(item string) {
	numericHash := doHash(item)

	hs.mu.Lock()
	i := sort.SearchInts(hs.hashes, int(numericHash))
	hs.hashes = append(hs.hashes[:i], append([]int{int(numericHash)}, hs.hashes[i:]...)...)
	hs.mu.Unlock()
}

func (hs *HashSearch) Has(item string) bool {
	numericHash := doHash(item)

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
func doHash(s string) uint64 {
	h := hasherPool.Get().(hash.Hash64)

	defer hasherPool.Put(h)
	h.Reset()
	_, err := h.Write([]byte(s))
	if err != nil {
		panic(err)
	}

	return h.Sum64()
}
