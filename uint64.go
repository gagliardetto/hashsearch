package hashsearch

import (
	"sort"
	"sync"
)

// Uint64s sorts a slice of uint64 in increasing order.
func Uint64s(x []uint64) { sort.Sort(Uint64Slice(x)) }

// Uint64Slice attaches the methods of Interface to []uint64, sorting in increasing order.
type Uint64Slice []uint64

func (x Uint64Slice) Len() int           { return len(x) }
func (x Uint64Slice) Less(i, j int) bool { return x[i] < x[j] }
func (x Uint64Slice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// Sort is a convenience method: x.Sort() calls Sort(x).
func (x Uint64Slice) Sort() { sort.Sort(x) }

func SearchUint64(a []uint64, x uint64) int {
	return sort.Search(len(a), func(i int) bool { return a[i] >= x })
}

// Uint64 is a utility to track uint64 values.
type Uint64 struct {
	mu     *sync.RWMutex
	values []uint64
}

func NewUint64() *Uint64 {
	return &Uint64{
		mu:     &sync.RWMutex{},
		values: make([]uint64, 0),
	}
}

// Sort sorts the index of values.
// This should be called after using `WARNING_UnsortedAppend`.
func (ia *Uint64) Sort() {
	ia.mu.Lock()
	Uint64s(ia.values)
	ia.mu.Unlock()
}

// Add adds the provided value to the index inserting it in a sorted way.
// NOTE: it accepts duplicate values.
func (ia *Uint64) Add(val uint64) {
	ia.mu.Lock()
	ia.noMutexAddSorted(val)
	ia.mu.Unlock()
}

// WARNING_UnsortedAppend appends (without order considerations)
// the provided value to the index.
// WARNING: this function shouldn't be normally used.
// If this function is used, then you should call the `Sort()` method
// so that the index will be able to work as expected.
func (ia *Uint64) WARNING_UnsortedAppend(val uint64) {
	ia.mu.Lock()
	ia.values = append(ia.values, val)
	ia.mu.Unlock()
}

// Has returns true if the index contains the provided value.
func (ia *Uint64) Has(val uint64) bool {
	ia.mu.RLock()
	has := ia.noMutexHas(val)
	ia.mu.RUnlock()
	return has
}

// HasOrAdd checks whether the provided value is already present;
// If the item is already present, the function returns true.
// If it is not present, the provided value is added
// and the function returns false.
func (ia *Uint64) HasOrAdd(val uint64) bool {
	ia.mu.Lock()
	has := ia.noMutexHas(val)
	if !has {
		ia.noMutexAddSorted(val)
	}
	ia.mu.Unlock()
	return has
}

func (ia *Uint64) noMutexAddSorted(val uint64) {
	i := SearchUint64(ia.values, val)
	ia.values = append(ia.values, 0 /* use the zero value of the element type */)
	copy(ia.values[i+1:], ia.values[i:])
	ia.values[i] = val
}

func (ia *Uint64) noMutexHas(val uint64) bool {
	index := SearchUint64(ia.values, val)
	if index == len(ia.values) {
		return false
	}
	has := val == ia.values[index]
	return has
}
