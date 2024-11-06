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

func mustOpen(fn string) *os.File {
	if fn == "-" {
		return os.Stdin
	}

	f, err := os.Open(fn)
	if err != nil {
		Errorf("Unable to open %s for reading: %s", fn, err)
	}
	return f
}

func mustCreate(fn string) *os.File {
	if fn == "-" {
		return os.Stdout
	}

	f, err := os.Create(fn)
	if err != nil {
		Errorf("Unable to open %s for writing: %s", fn, err)
	}
	return f
}

func main() {
	inputFile := "-"
	outputFile := "-"

	switch len(os.Args) {
	case 3:
		outputFile = os.Args[2]
		fallthrough
	case 2:
		inputFile = os.Args[1]
	case 1:
		// use defaults
	default:
		Errorf("Syntax: fmtftml [INPUT] [OUTPUT]")
		return
	}

	f := mustOpen(inputFile)
	defer f.Close()

	doc, err := ftml.Parse(f)
	if err != nil {
		Errorf("Unable to parse %s: %s", inputFile, err)
	}
	if err := f.Close(); err != nil {
		Errorf("Unable to close %s after reading: %s", inputFile, err)
	}

	f = mustCreate(outputFile)
	if err := markdown.Write(f, doc); err != nil {
		Errorf("Unable to write document to %s: %s", outputFile, err)
	}
	if err := f.Close(); err != nil {
		Errorf("Unable to close %s after writing: %s", outputFile, err)
	}
}
