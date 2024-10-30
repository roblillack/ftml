package main

import (
	"fmt"
	"os"

	"github.com/roblillack/ftml"
	"github.com/roblillack/ftml/markdown"
)

func Errorf(layout string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, layout+"\n", args...)
	os.Exit(1)
}

func main() {
	if len(os.Args) != 3 {
		Errorf("Syntax: fmtftml INPUT OUTPUT")
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

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

	f, err = os.Create(outputFile)
	if err != nil {
		Errorf("Unable to open %s for writing: %s", outputFile, err)
	}
	if err := markdown.Write(f, doc); err != nil {
		Errorf("Unable to write document to %s: %s", outputFile, err)
	}
	if err := f.Close(); err != nil {
		Errorf("Unable to close %s after writing: %s", outputFile, err)
	}
}
