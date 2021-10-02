package formatter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/roblillack/pure/ftml"
	"github.com/stretchr/testify/assert"
)

func TestSimpleParagraph(t *testing.T) {
	doc, err := ftml.Parse(strings.NewReader("<p>This is a paragraph.</p><p>And this is the second one.</p>"))
	if err != nil {
		t.Error(err)
	}

	res := `This is a paragraph.

And this is the second one.
`

	buf := &bytes.Buffer{}
	if assert.NoError(t, Write(buf, doc)) {
		assert.Equal(t, res, buf.String())
	}
}
