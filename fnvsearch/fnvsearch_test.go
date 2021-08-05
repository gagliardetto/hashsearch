package fnvsearch

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSearch(t *testing.T) {
	helloWorld := "hello world"
	foo := "foo"
	bar := "bar"
	{
		expected := []byte{0x6c, 0x15, 0x57, 0x99, 0xfd, 0xc8, 0xee, 0xc4, 0xb9, 0x15, 0x23, 0x80, 0x8e, 0x77, 0x26, 0xb7}
		h1 := HashBytes([]byte(helloWorld))
		require.True(t, bytes.Equal(expected, h1))

		h2 := HashString(helloWorld)
		require.True(t, bytes.Equal(expected, h2))

		require.Equal(t, h1, h2)

		h3 := HashReader(strings.NewReader(helloWorld))
		require.True(t, bytes.Equal(expected, h3))

		require.Equal(t, h1, h3)

		h4 := HashFromWriter(func(w io.Writer) error {
			_, err := w.Write([]byte(helloWorld))
			return err
		})
		require.True(t, bytes.Equal(expected, h4))

		require.Equal(t, h1, h4)
	}
	{
		hs := New()
		require.NotNil(t, hs)
		require.NotNil(t, hs.hashes)
		require.Len(t, hs.hashes, 0)
		require.NotNil(t, hs.mu)

		require.False(t, hs.HasString(helloWorld))
		require.False(t, hs.HasBytes([]byte(helloWorld)))
		require.False(t, hs.HasFromReader(strings.NewReader(helloWorld)))

		hs.AddFromString(helloWorld)
		require.Len(t, hs.hashes, 1)
		require.True(t, bytes.Equal(hs.hashes[0], HashBytes([]byte(helloWorld))))

		require.True(t, hs.HasOrAddFromString(helloWorld))
		require.False(t, hs.HasOrAddFromString(foo))
		require.False(t, hs.HasOrAddFromString(bar))
		require.True(t, hs.HasString(foo))
		require.True(t, hs.HasBytes([]byte(bar)))
		require.True(t, hs.HasFromReader(bytes.NewReader([]byte(bar))))
		require.False(t, hs.HasFromReader(bytes.NewReader([]byte("not found"))))
	}
}
