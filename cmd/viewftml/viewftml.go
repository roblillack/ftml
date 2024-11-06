package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/roblillack/ftml"
	"github.com/roblillack/ftml/formatter"
	"github.com/roblillack/ftml/html"
)

func Errorf(layout string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, layout+"\n", args...)
	os.Exit(1)
}

func createReader(inputFile string) (bool, io.ReadCloser, error) {
	if parsedURL, err := url.Parse(inputFile); err == nil && (parsedURL.Scheme == "http" || parsedURL.Scheme == "https") {
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Get(inputFile)
		if err != nil {
			return false, nil, err
		}
		return true, resp.Body, nil
	}

	expectHTML := strings.ToLower(filepath.Ext(inputFile)) != ".ftml"
	r, err := os.Open(inputFile)
	return expectHTML, r, err
}

func main() {
	disableANSI := flag.Bool("n", false, "Disable use of ANSI escape sequences.")
	saveFTML := flag.Bool("s", false, "Save the formatted FTML to standard out.")
	flag.Parse()

	if len(flag.Args()) != 1 {
		Errorf("Syntax: viewftml [-n] [-s] INPUT")
	}

	inputFile := flag.Arg(0)
	expectHTML, f, err := createReader(inputFile)
	if err != nil {
		Errorf("Unable to read %s: %s", inputFile, err)
	}
	defer f.Close()

	var doc *ftml.Document
	if expectHTML {
		doc, err = html.Parse(f)
	} else {
		doc, err = ftml.Parse(f)
	}
	if err != nil {
		Errorf("Unable to parse %s: %s", inputFile, err)
	}
	if err := f.Close(); err != nil {
		Errorf("Unable to close %s after reading: %s", inputFile, err)
	}

	if *saveFTML {
		if err := ftml.Write(os.Stdout, doc); err != nil {
			Errorf("Unable to write FTML document: %s", err)
		}
		return
	}

	if err := formatter.Write(os.Stdout, doc, !*disableANSI); err != nil {
		Errorf("Unable to write document to: %s", err)
	}
}
