package hashsearch

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUint64(t *testing.T) {
	{
		h1 := HashBytes([]byte("hello world"))
		require.Equal(t, uint64(0x779a65e7023cd2e7), h1)

		h2 := HashString("hello world")
		assert.Equal(t, uint64(0x779a65e7023cd2e7), h2)

		require.Equal(t, h1, h2)
	}

	{
		hs := New()
		require.NotNil(t, hs)
		require.NotNil(t, hs.hashes)
		require.NotNil(t, hs.hashes.values)
		require.NotNil(t, hs.hashes.mu)
		require.Len(t, hs.hashes.values, 0)

		require.False(t, hs.Has("hello world"))
		require.False(t, hs.HasBytes([]byte("hello world")))
		require.False(t, hs.HasOrAddBytes([]byte("hello world")))
		require.True(t, hs.HasBytes([]byte("hello world")))
		require.True(t, hs.Has("hello world"))

		hs.AddFromBytes([]byte("foo"))
		require.True(t, hs.HasBytes([]byte("foo")))
		require.True(t, hs.Has("foo"))

		hs.Add("bar")
		require.True(t, hs.HasBytes([]byte("bar")))
		require.True(t, hs.Has("bar"))

		expected := []uint64{
			// NOTE: the order is by hash value, not alphabetical by source string.
			HashBytes([]byte("bar")),
			HashBytes([]byte("hello world")),
			HashBytes([]byte("foo")),
		}

		require.Equal(t, expected, hs.hashes.values)

		hs.WARNING_UnsortedAppend("this will be last")
		expected = append(expected, HashBytes([]byte("this will be last")))
		require.Equal(t, expected, hs.hashes.values)

		hs.WARNING_UnsortedAppendBytes([]byte("last from bytes"))
		expected = append(expected, HashBytes([]byte("last from bytes")))
		require.Equal(t, expected, hs.hashes.values)

		hs.Sort()
		expectedSorted := []uint64{0x3934191339461a, 0x779a65e7023cd2e7, 0x9a51a095cb21ff14, 0xa88edf609886ba58, 0xdcb27518fed9d577}
		require.Equal(t, expectedSorted, hs.hashes.values)
	}
}

func TestInt(t *testing.T) {
	hs := NewInt()
	require.NotNil(t, hs)
	require.NotNil(t, hs.values)
	require.Len(t, hs.values, 0)
	require.NotNil(t, hs.mu)

	require.False(t, hs.Has(33))
	hs.Add(33)
	require.True(t, hs.Has(33))
	require.Equal(t, hs.values, []int{33})
	hs.Add(33)
	require.Equal(t, hs.values, []int{33, 33})
	hs.Add(33)
	require.Equal(t, hs.values, []int{33, 33, 33})
	hs.Add(1)
	require.Equal(t, hs.values, []int{1, 33, 33, 33})
	hs.Add(44)
	require.Equal(t, hs.values, []int{1, 33, 33, 33, 44})
	hs.WARNING_UnsortedAppend(-1)
	require.Equal(t, hs.values, []int{1, 33, 33, 33, 44, -1})
	hs.Sort()
	require.Equal(t, hs.values, []int{-1, 1, 33, 33, 33, 44})
	require.False(t, hs.HasOrAdd(997))
	require.True(t, hs.HasOrAdd(997))
	require.Equal(t, hs.values, []int{-1, 1, 33, 33, 33, 44, 997})
}
