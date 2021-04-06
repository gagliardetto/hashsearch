package fnvsearch

import (
	"bytes"
	"hash"
	"hash/fnv"
	"io"
	"sort"
	"sync"
)

var hasherPool128a *sync.Pool

func init() {
	hasherPool128a = &sync.Pool{
		New: func() interface{} {
			return fnv.New128a()
		},
	}
}

// HashString returns a hash.Hash obtained from
// a FNV-1a 128bit hash of the provided string.
func HashString(s string) []byte {
	return HashBytes([]byte(s))
}

// HashBytes returns a hash.Hash obtained from
// a FNV-1a 128bit hash of the provided byte slice.
func HashBytes(b []byte) []byte {
	h := hasherPool128a.Get().(hash.Hash)

	defer hasherPool128a.Put(h)
	h.Reset()

	if _, err := h.Write(b); err != nil {
		panic(err)
	}

	return h.Sum(nil)
}

// HashReader returns a hash.Hash obtained from
// a FNV-1a 128bit hash of the contents of the provided reader.
func HashReader(reader io.Reader) []byte {
	h := hasherPool128a.Get().(hash.Hash)

	defer hasherPool128a.Put(h)
	h.Reset()

	if _, err := io.Copy(h, reader); err != nil {
		panic(err)
	}

	return h.Sum(nil)
}

// HashFromWriter returns a hash.Hash obtained from
// a FNV-1a 128bit hash of the contents written to w.
func HashFromWriter(callback func(w io.Writer) error) []byte {
	h := hasherPool128a.Get().(hash.Hash)

	defer hasherPool128a.Put(h)
	h.Reset()

	if err := callback(h); err != nil {
		panic(err)
	}

	return h.Sum(nil)
}

// SearchByteSlices searches for x in a sorted slice of []byte and returns the index
// as specified by Search. The return value is the index to insert x if x is not
// present (it could be len(a)).
// The slice must be sorted in ascending order.
//
func SearchByteSlices(a [][]byte, x []byte) int {
	return sort.Search(len(a), func(i int) bool { comp := bytes.Compare(a[i], x); return comp >= 0 })
}

// ByteSliceSlice attaches the methods of sort.Interface to [][]byte, sorting in increasing order.
type ByteSliceSlice [][]byte

func (p ByteSliceSlice) Len() int           { return len(p) }
func (p ByteSliceSlice) Less(i, j int) bool { return bytes.Compare(p[i], p[j]) == -1 }
func (p ByteSliceSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Tracker is used to track hashes.
type Tracker struct {
	mu     *sync.RWMutex
	hashes [][]byte
}

// New initializes and returns a new Tracker object.
func New() *Tracker {
	return &Tracker{
		mu: &sync.RWMutex{},
	}
}

// WARNING_UnsortedAppendHashFromString is a non-safe method that appends (without order considerations)
// a new hash (obtained by hashing the provided string) to the hash tracker.
func (hs *Tracker) WARNING_UnsortedAppendHashFromString(item string) {
	hs.WARNING_UnsortedAppendHashFromBytes([]byte(item))
}

// WARNING_UnsortedAppendHashFromBytes is a non-safe method that appends (without order considerations)
// a new hash (obtained by hashing the provided byte slice) to the hash tracker.
func (hs *Tracker) WARNING_UnsortedAppendHashFromBytes(item []byte) {
	hs.WARNING_UnsortedAppendHash(HashBytes(item))
}

// WARNING_UnsortedAppendHash is a non-safe method that appends (without order considerations)
// a new hash to the hash tracker.
func (hs *Tracker) WARNING_UnsortedAppendHash(hash []byte) {
	hs.mu.Lock()
	hs.hashes = append(hs.hashes, hash)
	hs.mu.Unlock()
}

// Sort sorts the hash tracker slice.
func (hs *Tracker) Sort() {
	hs.mu.Lock()
	sort.Sort(ByteSliceSlice(hs.hashes))
	hs.mu.Unlock()
}

// AddFromString adds a new hash obtained from hashing the provided string.
func (hs *Tracker) AddFromString(item string) {
	hs.AddFromBytes([]byte(item))
}

// AddFromBytes adds a new hash obtained from hashing the provided byte slice.
func (hs *Tracker) AddFromBytes(item []byte) {
	hs.AddHash(HashBytes(item))
}

// AddFromReader adds a new hash obtained from hashing the contents of the provided io.Reader.
func (hs *Tracker) AddFromReader(reader io.Reader) {
	hs.AddHash(HashReader(reader))
}

// AddHash adds the provided hash to the hash tracker.
func (hs *Tracker) AddHash(hash []byte) {
	hs.mu.Lock()
	hs.noMutexAddSortedHash(hash)
	hs.mu.Unlock()
}

// HasString returns true if the hash tracker contains the hash of the provided string.
func (hs *Tracker) HasString(item string) bool {
	return hs.HasBytes([]byte(item))
}

// HasBytes returns true if the hash tracker contains the hash of the provided byte slice.
func (hs *Tracker) HasBytes(item []byte) bool {
	return hs.HasByHash(HashBytes(item))
}

// HasFromreader returns true if the hash tracker contains the hash of the contents of the provided io.Reader.
func (hs *Tracker) HasFromreader(reader io.Reader) bool {
	return hs.HasByHash(HashReader(reader))
}

// HasByHash returns true if the hash tracker contains the provided hash.
func (hs *Tracker) HasByHash(hash []byte) bool {
	hs.mu.RLock()
	has := hs.noMutexHas(hash)
	hs.mu.RUnlock()
	return has
}

// HasOrAddFromString returns true if the hash tracker contains the hash of the provided string;
// if it does not, then the hash will be added.
func (hs *Tracker) HasOrAddFromString(item string) bool {
	return hs.HasOrAddHash(HashString(item))
}

// HasOrAddFromBytes returns true if the hash tracker contains the hash of the provided byte slice;
// if it does not, then the hash will be added.
func (hs *Tracker) HasOrAddFromBytes(item []byte) bool {
	return hs.HasOrAddHash(HashBytes(item))
}

// HasOrAddFromReader returns true if the hash tracker contains the hash of the contents of the provided io.Reader;
// if it does not, then the hash will be added.
func (hs *Tracker) HasOrAddFromReader(reader io.Reader) bool {
	return hs.HasOrAddHash(HashReader(reader))
}

// HasOrAddHash returns true if the hash tracker contains the provided hash;
// if it does not, then the hash will be added.
func (hs *Tracker) HasOrAddHash(hash []byte) bool {
	hs.mu.Lock()
	has := hs.noMutexHas(hash)
	if !has {
		hs.noMutexAddSortedHash(hash)
	}
	hs.mu.Unlock()
	return has
}

func (hs *Tracker) noMutexAddSortedHash(hash []byte) {
	i := SearchByteSlices(hs.hashes, hash)
	hs.hashes = append(hs.hashes, nil /* use the zero value of the element type */)
	copy(hs.hashes[i+1:], hs.hashes[i:])
	hs.hashes[i] = hash
}

func (hs *Tracker) noMutexHas(hash []byte) bool {
	index := SearchByteSlices(hs.hashes, hash)
	if index == len(hs.hashes) {
		return false
	}
	// has:
	return bytes.Equal(hs.hashes[index], hash)
}
