package main

import (
	"fmt"
	"os"

	"github.com/roblillack/pure/formatter"
	"github.com/roblillack/pure/ftml"
)

func Errorf(layout string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, layout+"\n", args...)
	os.Exit(1)
}

func main() {
	if len(os.Args) != 2 {
		Errorf("Syntax: viewftml INPUT")
	}

	inputFile := os.Args[1]

	f, err := os.Open(inputFile)
	if err != nil {
		Errorf("Unable to read %s: %s", inputFile, err)
	}
	defer f.Close()

	doc, err := ftml.Parse(f)
	if err != nil {
		Errorf("Unable to parse %s: %s", inputFile, err)
	}
	if err := f.Close(); err != nil {
		Errorf("Unable to close %s after reading: %s", inputFile, err)
	}

	if err := formatter.Write(os.Stdout, doc); err != nil {
		Errorf("Unable to write document to: %s", err)
	}
}
