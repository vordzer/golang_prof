package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

func exists(path string) (bool, error) {
	if len(path) == 0 {
		return false, nil
	}
	sf, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if !sf.IsDir() || os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func StringNormalize(s string) string {
	s = strings.TrimRight(s, "\t ")
	return s
}

func FirstString(filename string) (string, error) {
	inFile, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	_, scanner, err := bufio.ScanLines(inFile, true)
	if err != nil {
		return "", err
	}
	output := bytes.ReplaceAll(scanner, []byte("\x00"), []byte("\n"))
	return StringNormalize(string(output)), nil
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	// Place your code here
	if res, err := exists(dir); !res || err != nil {
		fmt.Println("Wrong path")
		return nil, err
	}
	if err := os.Chdir(dir); err != nil {
		fmt.Println("Can't go to directory")
		return nil, err
	}
	rd, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Can't read directory")
		return nil, err
	}
	for _, el := range rd {
		if el.IsDir() || strings.Contains(el.Name(), "=") {
			continue
		}
		os.Unsetenv(el.Name())
		if value, err := FirstString(el.Name()); err == nil {
			if len(value) != 0 {
				os.Setenv(el.Name(), value)
			}
		}
	}
	env := Environment{}
	for _, key := range os.Environ() {
		e := strings.Split(key, "=")
		env[e[0]] = EnvValue{e[1], false}
	}

	return env, nil
}
