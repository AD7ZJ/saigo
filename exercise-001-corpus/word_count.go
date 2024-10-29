package main

import (
	"corpus/corpus"
	"fmt"
	"os"
)

func main() {
	// Check if the file name is provided as an argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <filename>")
		return
	}

	// Get the file name from command-line arguments
	filename := os.Args[1]

	//open file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Reads the entire file content into a []byte
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	histogram := corpus.Analysis(string(content))
	//histogram = corpus.Analysis("Are you serious? I knew you were.")

	for _, keyVal := range histogram {
		if len(keyVal.Word) >= 8 {
			fmt.Printf("%s\t%d\n", keyVal.Word, keyVal.Count)
		} else {
			fmt.Printf("%s\t\t%d\n", keyVal.Word, keyVal.Count)
		}

	}
}
