package formatter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/roblillack/ftml"
	"github.com/stretchr/testify/assert"
)

func TestHardNewlines(t *testing.T) {
	for _, testCase := range []struct {
		input  []string
		output string
	}{
		{[]string{"Line A\n", "Line B"}, "| Line A\n| Line B\n"},
		{[]string{"Line A", "\n", "Line B"}, "| Line A\n| Line B\n"},
		{[]string{"Line A", "\nLine B"}, "| Line A\n| Line B\n"},
		{[]string{"This is a paragraph that contains a very long line of highlighted text to force the formatter to break"}, "| This is a paragraph that contains a very long line of highlighted text\n| to force the formatter to break\n"},
	} {
		spans := make([]ftml.Span, len(testCase.input))
		for idx, txt := range testCase.input {
			spans[idx] = ftml.Span{Text: txt}
		}
		inputDoc := ftml.Document{
			Paragraphs: []*ftml.Paragraph{
				{
					Type: ftml.QuoteParagraph,
					Children: []*ftml.Paragraph{
						{
							Type:    ftml.TextParagraph,
							Content: spans,
						},
					},
				},
			},
		}
		buf := &bytes.Buffer{}
		if assert.NoError(t, Write(buf, &inputDoc, false)) {
			assert.Equal(t, testCase.output, buf.String(), "Output incorrect for input: %+q", testCase.input)
		}
	}
}

func TestSimpleParagraph(t *testing.T) {
	doc, err := ftml.Parse(strings.NewReader("<p>This is a paragraph.</p><p>And this is the second one.</p>"))
	if err != nil {
		t.Error(err)
	}

	res := `This is a paragraph.

And this is the second one.
`

	buf := &bytes.Buffer{}
	if assert.NoError(t, Write(buf, doc, false)) {
		assert.Equal(t, res, buf.String())
	}
}

func TestWordWrap(t *testing.T) {
	doc, err := ftml.Parse(strings.NewReader("<p>FTML is text markup language which is designed to offer humans a way to express themselves better than plain text, but without the complexity of HTML, Markdown, or other document formats. An FTML file makes no assumptions about how the rendered text will look like, but only about the structure of the text. The format covers typical simple text documents, such as emails, memos, notes, online help, and is specifically suitable for copy-pasting text from one document to another.</p>"))
	if err != nil {
		t.Error(err)
	}

	//   ------------ 72 characters --------------------------------------------v
	res := "" +
		"FTML is text markup language which is designed to offer humans a way to\n" +
		"express themselves better than plain text, but without the complexity of\n" +
		"HTML, Markdown, or other document formats. An FTML file makes no\n" +
		"assumptions about how the rendered text will look like, but only about\n" +
		"the structure of the text. The format covers typical simple text\n" +
		"documents, such as emails, memos, notes, online help, and is\n" +
		"specifically suitable for copy-pasting text from one document to\n" +
		"another.\n"
	//   ------------ 72 characters --------------------------------------------^
	buf := &bytes.Buffer{}
	if assert.NoError(t, Write(buf, doc, false)) {
		assert.Equal(t, res, buf.String())
	}
}
