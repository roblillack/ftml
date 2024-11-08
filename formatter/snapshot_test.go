package formatter_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/roblillack/ftml"
	"github.com/roblillack/ftml/formatter"
)

func CompareDoc(t *testing.T, doc *ftml.Document, snapshotFile string) {
	buf := &bytes.Buffer{}
	if err := formatter.Write(buf, doc, false); err != nil {
		t.Fatal(err)
	}
	actual := buf.Bytes()

	if err := os.Remove(snapshotFile + ".new"); err != nil && !os.IsNotExist(err) {
		t.Fatalf("Unable to remove old snapshot: %s", err)
	}

	haveError := func(layout string, args ...interface{}) {
		out, err := os.Create(snapshotFile + ".new")
		if err == nil {
			defer out.Close()
			_, err = out.Write(actual)
		}
		if err != nil {
			t.Errorf("%s\nUnable to write new snapshot: %s", fmt.Sprintf(layout, args...), err)
		}
		t.Errorf("%s\nNew snapshot written to %s", fmt.Sprintf(layout, args...), snapshotFile+".new")
	}

	expected, err := os.ReadFile(snapshotFile)
	if err != nil {
		haveError("Unable to read snapshot: %s", err)
	}

	e := strings.Split(string(expected), "\n")
	a := strings.Split(string(actual), "\n")
	if len(e) != len(a) {
		haveError("Number of lines differ: %d != %d", len(e), len(a))
		return

	}
	for i := 0; i < len(e) && i < len(a); i++ {
		if e[i] != a[i] {
			haveError("Line %d differs:\nExpected: %s\nGot: %s", i, e[i], a[i])
			return
		}
	}
}

func TestExportingFiles(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)

	testdocs, err := filepath.Glob(filepath.Join(filepath.Dir(filename), "..", "examples", "*.ftml"))
	if err != nil {
		t.Fatal(err)
	}

	for _, exampleFn := range testdocs {
		outputFile := filepath.Join(filepath.Dir(filename), "testdata",
			strings.TrimSuffix(filepath.Base(exampleFn), filepath.Ext(exampleFn))+".txt")
		f, err := os.Open(exampleFn)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		doc, err := ftml.Parse(f)
		if err != nil {
			t.Fatal(err)
		}

		CompareDoc(t, doc, outputFile)
	}
}
