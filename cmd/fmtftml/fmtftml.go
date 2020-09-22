package main

import (
	"fmt"
	"os"

	"github.com/roblillack/pure/ftml"
)

func Errorf(layout string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, layout+"\n", args...)
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		Errorf("Syntax: fmtftml [FILE...]")
	}
	for _, fn := range os.Args[1:] {
		f, err := os.Open(fn)
		if err != nil {
			Errorf("Unable to read %s: %s", fn, err)
		}
		defer f.Close()

		doc, err := ftml.Parse(f)
		if err != nil {
			Errorf("Unable to parse %s: %s", fn, err)
		}
		if err := f.Close(); err != nil {
			Errorf("Unable to close %s after reading: %s", fn, err)
		}
		if err := os.Rename(fn, fn+".bak"); err != nil {
			Errorf("Unable to rename %s: %s", fn, err)
		}

		f, err = os.Create(fn)
		if err != nil {
			Errorf("Unable to open %s for writing: %s", fn, err)
		}
		if err := ftml.Write(f, doc); err != nil {
			Errorf("Unable to write document to %s: %s", fn, err)
		}
		if err := f.Close(); err != nil {
			Errorf("Unable to close %s after writing: %s", fn, err)
		}
		if err := os.Remove(fn + ".bak"); err != nil {
			Errorf("Unable to remove backup %s.bak: %s", fn, err)
		}
	}
}
