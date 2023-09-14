package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	// Place your code here
	t.Run("StringNormalize", func(t *testing.T) {
		testStr := "test \t"
		result := StringNormalize(testStr)
		require.Equal(t, "test", result)
	})
	t.Run("exists fail", func(t *testing.T) {
		result, _ := exists("")
		require.Equal(t, false, result)
	})
	t.Run("exists success", func(t *testing.T) {
		fIn, _ := os.CreateTemp("", "tmpfile-")
		defer func() {
			fIn.Close()
			os.Remove(fIn.Name())
		}()
		result, err := exists(fIn.Name())
		require.True(t, result)
		require.NoError(t, err)
	})
	t.Run("first string", func(t *testing.T) {
		fIn, _ := os.CreateTemp("", "tmpfile-")
		defer func() {
			fIn.Close()
			os.Remove(fIn.Name())
		}()
		fIn.Write([]byte("test\nnot need"))
		result, err := FirstString(fIn.Name())
		require.Equal(t, "test", result)
		require.NoError(t, err)
	})
}
