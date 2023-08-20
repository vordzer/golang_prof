package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrFromFileNotExist      = errors.New("file for copy is not exist")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrInvalidParam          = errors.New("invalid param")
)

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, ErrFromFileNotExist
	}
	return false, err
}

func validateParam(from, _ string, offset, limit int64) (bool, error) {
	if limit < 0 || offset < 0 {
		return false, ErrInvalidParam
	}
	if res, err := exists(from); !res {
		return false, err
	}
	return true, nil
}

func NeedRead(inFile *os.File, offset, limit int64) (int64, error) {
	sf, err := inFile.Stat()
	if err != nil {
		return 0, err
	}
	if sf.Size() < offset {
		return 0, ErrOffsetExceedsFileSize
	}
	allRead := sf.Size()
	if limit != 0 && allRead > limit+offset {
		allRead = limit
	} else if offset != 0 {
		allRead -= offset
	}
	return allRead, nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Place your code here.
	if ok, err := validateParam(fromPath, toPath, offset, limit); !ok {
		return err
	}
	// Open and check file for reading
	inFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer func() {
		if err := inFile.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Open file for writing
	outFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer func() {
		if err := outFile.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	allRead, err := NeedRead(inFile, offset, limit)
	if err != nil {
		return err
	}
	if allRead == 0 {
		return nil // empty copy or endless file
	}
	percentSize := allRead / 100 // not optimized, for beauty percent
	if percentSize == 0 {
		percentSize = 1
	}

	err = nil
	curLoad := 0
	readed := 0
	step := int(100 / allRead)
	if step == 0 {
		step = 1
	}
	buf := make([]byte, percentSize)
	if offset != 0 {
		skipBuf := make([]byte, offset)
		if _, err := inFile.Read(skipBuf); err != nil {
			return err
		}
	}
	for allRead != 0 && !errors.Is(err, io.EOF) {
		readed, err = inFile.Read(buf)
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
		if readed == 0 {
			break
		}
		if readed > int(allRead) {
			readed = int(allRead)
			allRead = 0
		} else {
			allRead -= int64(readed)
		}
		_, errw := outFile.Write(buf[:readed])
		if errw != nil {
			return errw
		}
		fmt.Printf("[%v%%]\t%v -> %v\n", curLoad, inFile.Name(), outFile.Name())
		curLoad += step
	}
	fmt.Printf("[%v%%]\t%v -> %v\n", 100, inFile.Name(), outFile.Name())

	return nil
}
