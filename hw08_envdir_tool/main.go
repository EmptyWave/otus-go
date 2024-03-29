package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Expected at least 2 arguments")
	}

	dir := os.Args[1]
	command := os.Args[2:]

	env, err := ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(RunCmd(command, env))
}
