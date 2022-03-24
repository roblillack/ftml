package main

import (
	"flag"
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
	disableANSI := flag.Bool("n", false, "Disable use of ANSI escape sequences.")
	flag.Parse()

	if len(flag.Args()) != 1 {
		Errorf("Syntax: viewftml [-n] INPUT")
	}

	inputFile := flag.Arg(0)

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

	if err := formatter.Write(os.Stdout, doc, !*disableANSI); err != nil {
		Errorf("Unable to write document to: %s", err)
	}
}
