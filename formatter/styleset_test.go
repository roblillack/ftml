package formatter

import (
	"testing"

	"github.com/roblillack/ftml"
	"github.com/stretchr/testify/assert"
)

func TestStyleSet(t *testing.T) {
	assert.True(t, StyleSet{}.Empty())
	assert.Equal(t, StyleSet{ftml.StyleBold: 1}, StyleSet{}.Add(ftml.StyleBold))
	s := StyleSet{}.Add(ftml.StyleItalic)
	assert.NotSame(t, s, s.Add(ftml.StyleItalic))
	assert.Equal(t, StyleSet{}.Add(ftml.StyleBold).All(), StyleSet{}.Add(ftml.StyleBold).Add(ftml.StyleBold).All())
	assert.Equal(t, []ftml.InlineStyle{ftml.StyleBold}, StyleSet{}.Add(ftml.StyleBold).All())
	assert.Equal(t, []ftml.InlineStyle{ftml.StyleCode}, StyleSet{}.Add(ftml.StyleCode).All())
}
