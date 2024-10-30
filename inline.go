package ftml

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type InlineStyle uint8

const (
	StyleNone InlineStyle = iota
	StyleBold
	StyleItalic
	StyleHighlight
	StyleUnderline
	StyleStrike
	StyleLink
	StyleCode
)

func (s InlineStyle) String() string {
	switch s {
	case StyleNone:
		return "text"
	case StyleBold:
		return "bold"
	case StyleItalic:
		return "italic"
	case StyleUnderline:
		return "underline"
	case StyleStrike:
		return "striked"
	case StyleHighlight:
		return "highlight"
	case StyleLink:
		return "link"
	case StyleCode:
		return "code"
	}

	panic("Unknown Inline Style")
}

type Span struct {
	Style      InlineStyle
	Text       string
	LinkTarget string
	Children   []Span
}

func (s *Span) Width() int {
	l := utf8.RuneCountInString(s.Text)
	for _, i := range s.Children {
		l += i.Width()
	}
	return l
}

func (s *Span) String() string {
	b := &strings.Builder{}

	if len(s.Children) > 0 {
		b.WriteString(fmt.Sprintf("[%s:", s.Style))
		for _, i := range s.Children {
			b.WriteString(i.String())
		}
		b.WriteString("]")
	} else {
		b.WriteString("‘")
		b.WriteString(s.Text)
		b.WriteString("’")
	}

	return b.String()
}
