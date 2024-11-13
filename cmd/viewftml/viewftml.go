package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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

func runPager() (*exec.Cmd, io.Writer, error) {
	cmdline := []string{"less"}
	if pager := os.Getenv("PAGER"); pager != "" {
		cmdline = strings.Split(pager, " ")
	}

	if cmdline[0] == "" {
		cmdline[0] = "less"
	}

	if prog := filepath.Base(cmdline[0]); prog == "less" || prog == "more" {
		cmdline = append(cmdline, "-R")
	}

	r, w := io.Pipe()
	cmd := exec.Command(cmdline[0], cmdline[1:]...)
	cmd.Stdin = r
	cmd.Stdout = os.Stdout

	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}

	return cmd, w, nil
}

func terminalWidth(defaultWidth int) int {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		log.Println(err)
		return defaultWidth
	}

	s := string(out)
	s = strings.TrimSpace(s)
	sArr := strings.Split(s, " ")
	width, err := strconv.Atoi(sArr[1])
	if err != nil {
		log.Println(err)
		return defaultWidth
	}
	return width
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

	var cmd *exec.Cmd
	var enc *formatter.Formatter
	if o, _ := os.Stdout.Stat(); (o.Mode()&os.ModeCharDevice) != os.ModeCharDevice || *disableANSI {
		enc = formatter.NewASCII(os.Stdout)
	} else {
		c, pager, err := runPager()
		if err != nil {
			Errorf("Unable to run pager: %s", err)
		}
		enc = formatter.NewANSI(pager)
		if w := terminalWidth(80); w < 60 {
			enc.Style.WrapWidth = w
			enc.Style.LeftPadding = 0
		} else if w < 100 {
			enc.Style.WrapWidth = w - 2
			enc.Style.LeftPadding = 2
		} else {
			padding := (w-100)/2 + 4
			enc.Style.WrapWidth = w - padding
			enc.Style.LeftPadding = padding
		}
		cmd = c
	}

	if err := enc.WriteDocument(doc); err != nil {
		Errorf("Unable to write document to: %s", err)
	}

	if cmd != nil {
		cmd.Wait()
	}
}
