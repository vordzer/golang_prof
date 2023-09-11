package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func makeTestFiles() (*os.File, *os.File, error) {
	fIn, err := os.CreateTemp("", "tmpfile-")
	if err != nil {
		fmt.Printf("Can't create for in")
		return nil, nil, err
	}
	fOut, err := os.CreateTemp("", "tmpfile-")
	if err != nil {
		fmt.Printf("Can't create for in")
		return nil, nil, err
	}
	return fIn, fOut, nil
}

func CloseRemoveTestFiles(fIn, fOut *os.File) {
	fIn.Close()
	fOut.Close()
	os.Remove(fIn.Name())
	os.Remove(fOut.Name())
}

func TestCopy(t *testing.T) {
	// Place your code here.
	t.Run("invalid offset", func(t *testing.T) {
		err := Copy("", "", -1, 0)

		require.True(t, errors.Is(err, ErrInvalidParam))
	})
	t.Run("invalid limit", func(t *testing.T) {
		fIn, fOut, err := makeTestFiles()
		if err != nil {
			fmt.Printf("Can't create tmpfile")
		}
		defer CloseRemoveTestFiles(fIn, fOut)
		err = Copy(fIn.Name(), fOut.Name(), 0, -1)

		require.True(t, errors.Is(err, ErrInvalidParam))
	})
	t.Run("invalid filename", func(t *testing.T) {
		err := Copy("fakename1", "fakename2", 0, 0)

		require.True(t, errors.Is(err, ErrFromFileNotExist))
	})
	t.Run("Same from and to", func(t *testing.T) {
		fIn, _, err := makeTestFiles()
		if err != nil {
			fmt.Printf("Can't create tmpfile")
		}
		err = Copy(fIn.Name(), fIn.Name(), 0, 0)

		require.True(t, errors.Is(err, ErrInvalidParam))
	})
	t.Run("file less offset", func(t *testing.T) {
		fIn, fOut, err := makeTestFiles()
		if err != nil {
			fmt.Printf("Can't create tmpfile")
		}
		defer CloseRemoveTestFiles(fIn, fOut)

		err = Copy(fIn.Name(), fOut.Name(), 20, 0)
		require.True(t, errors.Is(err, ErrOffsetExceedsFileSize))
	})
	t.Run("empty copy", func(t *testing.T) {
		fIn, fOut, err := makeTestFiles()
		if err != nil {
			fmt.Printf("Can't create tmpfile")
		}
		defer CloseRemoveTestFiles(fIn, fOut)

		err = Copy(fIn.Name(), fOut.Name(), 0, 0)
		require.True(t, errors.Is(err, nil))
	})
	t.Run("endless copy", func(t *testing.T) {
		fIn, fOut, err := makeTestFiles()
		if err != nil {
			fmt.Printf("Can't create tmpfile")
		}
		defer CloseRemoveTestFiles(fIn, fOut)

		err = Copy("/dev/urandom", fOut.Name(), 0, 0)
		require.True(t, errors.Is(err, ErrUnsupportedFile))
	})
	t.Run("Success copy limit file", func(t *testing.T) {
		fIn, fOut, err := makeTestFiles()
		if err != nil {
			fmt.Printf("Can't create tmpfile")
		}
		defer CloseRemoveTestFiles(fIn, fOut)

		buf := make([]byte, 200)
		rand.Read(buf)
		fIn.Write(buf)
		resBuf := make([]byte, 200)

		err = Copy(fIn.Name(), fOut.Name(), 0, 50)
		n, _ := fOut.Read(resBuf)

		require.True(t, errors.Is(err, nil))
		require.Equal(t, buf[0:50], resBuf[0:n])
	})
	t.Run("Success copy full file", func(t *testing.T) {
		fIn, fOut, err := makeTestFiles()
		if err != nil {
			fmt.Printf("Can't create tmpfile")
		}
		defer CloseRemoveTestFiles(fIn, fOut)

		buf := make([]byte, 200)
		rand.Read(buf)
		fIn.Write(buf)
		resBuf := make([]byte, 200)

		err = Copy(fIn.Name(), fOut.Name(), 0, 0)
		fOut.Read(resBuf)

		require.True(t, errors.Is(err, nil))
		require.Equal(t, buf, resBuf)
	})
	t.Run("Success copy full file with offset", func(t *testing.T) {
		fIn, fOut, err := makeTestFiles()
		if err != nil {
			fmt.Printf("Can't create tmpfile")
		}
		defer CloseRemoveTestFiles(fIn, fOut)

		buf := make([]byte, 200)
		rand.Read(buf)
		fIn.Write(buf)
		resBuf := make([]byte, 200)

		err = Copy(fIn.Name(), fOut.Name(), 50, 0)
		fOut.Read(resBuf)

		require.True(t, errors.Is(err, nil))
		require.Equal(t, buf[50:], resBuf[:150])
	})
	t.Run("Success copy full file with offset and limit", func(t *testing.T) {
		fIn, fOut, err := makeTestFiles()
		if err != nil {
			fmt.Printf("Can't create tmpfile")
		}
		defer CloseRemoveTestFiles(fIn, fOut)

		buf := make([]byte, 200)
		rand.Read(buf)
		fIn.Write(buf)
		resBuf := make([]byte, 200)

		err = Copy(fIn.Name(), fOut.Name(), 50, 100)
		fOut.Read(resBuf)

		require.True(t, errors.Is(err, nil))
		require.Equal(t, buf[50:150], resBuf[:100])
	})
}
