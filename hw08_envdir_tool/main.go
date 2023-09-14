package main

import (
	"fmt"
	"os"
)

func main() {
	// Place your code here.
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) < 2 {
		fmt.Println("Not enough arguments")
		os.Exit(1)
	}
	env, err := ReadDir(argsWithoutProg[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(RunCmd(argsWithoutProg[1:], env))
}
