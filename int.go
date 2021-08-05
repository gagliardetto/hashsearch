package hashsearch

import (
	"sort"
	"sync"
)

type Int struct {
	mu     *sync.RWMutex
	values []int
}

func NewInt() *Int {
	return &Int{
		mu:     &sync.RWMutex{},
		values: make([]int, 0),
	}
}

// Sort sorts the index of values.
// This should be called after using `WARNING_UnsortedAppend`.
func (ia *Int) Sort() {
	ia.mu.Lock()
	sort.Ints(ia.values)
	ia.mu.Unlock()
}

// Add adds the provided value to the index inserting it in a sorted way.
// NOTE: it accepts duplicate values.
func (ia *Int) Add(val int) {
	ia.mu.Lock()
	ia.noMutexAddSorted(val)
	ia.mu.Unlock()
}

// WARNING_UnsortedAppend appends (without order considerations)
// the provided value to the index.
// WARNING: this function shouldn't be normally used.
// If this function is used, then you should call the `Sort()` method
// so that the index will be able to work as expected.
func (ia *Int) WARNING_UnsortedAppend(val int) {
	ia.mu.Lock()
	ia.values = append(ia.values, val)
	ia.mu.Unlock()
}

// Has returns true if the index contains the provided value.
func (ia *Int) Has(val int) bool {
	ia.mu.RLock()
	has := ia.noMutexHas(val)
	ia.mu.RUnlock()
	return has
}

// HasOrAdd checks whether the provided value is already present;
// If the item is already present, the function returns true.
// If it is not present, the provided value is added
// and the function returns false.
func (ia *Int) HasOrAdd(val int) bool {
	ia.mu.Lock()
	has := ia.noMutexHas(val)
	if !has {
		ia.noMutexAddSorted(val)
	}
	ia.mu.Unlock()
	return has
}

func (ia *Int) noMutexAddSorted(val int) {
	i := sort.SearchInts(ia.values, val)
	ia.values = append(ia.values, 0 /* use the zero value of the element type */)
	copy(ia.values[i+1:], ia.values[i:])
	ia.values[i] = val
}

func (ia *Int) noMutexHas(val int) bool {
	index := sort.SearchInts(ia.values, val)
	if index == len(ia.values) {
		return false
	}
	has := val == ia.values[index]
	return has
}
