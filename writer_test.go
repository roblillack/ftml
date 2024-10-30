package ftml

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleParagraph(t *testing.T) {
	res := "<p>This is a test.</p>\n"
	buf := &bytes.Buffer{}
	if assert.NoError(t, Write(buf, doc(p__("This is a test.")))) {
		assert.Equal(t, res, buf.String())
	}
}

func TestParagraphWithStyles(t *testing.T) {
	res := `<p>This is <i>a little more complex</i> test.</p>
`
	buf := &bytes.Buffer{}
	if assert.NoError(t, Write(buf, doc(p_(span("This is "), i__("a little more complex"), span(" test."))))) {
		assert.Equal(t, res, buf.String())
	}
}

func TestSimpleStyles(t *testing.T) {
	tests := map[string][]Span{
		`This is a test.`:                      {span("This is a test.")},
		`This is a <b>test</b>.`:               {span("This is a "), b__("test"), span(".")},
		`This is a <b> test </b>.`:             {span("This is a "), b__(" test "), span(".")},
		`This is a <b><i>second</i> test</b>.`: {span("This is a "), b_(i__("second"), span(" test")), span(".")},
	}

	for str, spans := range tests {
		log.Printf("Running tests for “%s” ...\n", strings.TrimSpace(str))
		buf := &strings.Builder{}
		dbg := &strings.Builder{}
		mock := &output{Writer: buf, MaxWidth: 1000}
		for idx, x := range spans {
			assert.NoError(t, mock.writeSpan(x, 0, idx == 0, idx == len(spans)-1))
			dbg.WriteString(x.String())
		}
		assert.Equalf(t, str, buf.String(), "written FTML not corrent for %s", dbg.String())
	}
}

func TestWriteSpaces(t *testing.T) {
	tests := map[string][]Span{
		"<p>&emsp14;This is a test.</p>\n":    {span(" This is a test.")},
		"<p>This is a test.&emsp14;</p>\n":    {span("This is a test. ")},
		"<p>A&emsp14;&emsp14;&emsp14;B</p>\n": {span("A   B")},
		"<p>\n  A<br />\n  &emsp14;B\n</p>\n": {span("A\n B")},
		"<p>\n  A&emsp14;<br />\n  B\n</p>\n": {span("A \nB")},
	}

	for str, spans := range tests {
		log.Printf("Running tests for “%s” ...\n", strings.TrimSpace(str))

		buf := &bytes.Buffer{}
		if assert.NoError(t, Write(buf, doc(p_(spans...)))) {
			assert.Equal(t, str, buf.String())
		}
	}
}
