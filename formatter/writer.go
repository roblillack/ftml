package formatter

import (
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/roblillack/pure/ftml"
)

type Breadcrumbs []ftml.ParagraphType

type Formatter struct {
	Document *ftml.Document
	Writer   io.Writer
	ANSI     bool

	currentLineSpacing int
}

func Write(w io.Writer, d *ftml.Document, ansiCodes bool) error {
	f := &Formatter{
		Writer:   w,
		Document: d,
		ANSI:     ansiCodes,
	}

	return f.WriteParagraphs(d.Paragraphs, "", "")
}

const WrapWidth = 72

// As per ECMA-48, 5th edition:
// https://www.ecma-international.org/wp-content/uploads/ECMA-48_5th_edition_june_1991.pdf
var styleTags = map[ftml.InlineStyle]struct{ B, E string }{
	ftml.StyleBold:      {"\033[1m", "\033[22m"},
	ftml.StyleItalic:    {"\033[3m", "\033[23m"},
	ftml.StyleHighlight: {"\033[7m", "\033[27m"},
	ftml.StyleUnderline: {"\033[4m", "\033[24m"},
	ftml.StyleCode:      {"", ""},
	ftml.StyleStrike:    {"\033[9m", "\033[29m"},
}

const resetAllModes = "\033[0m"

func (f *Formatter) WriteCentered(span ftml.Span, length int, followPrefix string) (int, error) {
	l := span.Width()

	for i := math.Floor(float64(WrapWidth-len(followPrefix))/2 - float64(l)/2); i > 0; i-- {
		if _, err := io.WriteString(f.Writer, " "); err != nil {
			return length, err
		}
		length++
	}

	return f.WriteSpan(span, length, followPrefix, StyleSet{})
}

func (f *Formatter) WriteSpan(span ftml.Span, length int, followPrefix string, outerStyles StyleSet) (int, error) {
	currentStyles := outerStyles.Add(span.Style)
	tag := styleTags[span.Style]
	if f.ANSI && tag.B != "" {
		if _, err := io.WriteString(f.Writer, tag.B); err != nil {
			return length, err
		}
	}

	if span.Style == ftml.StyleNone {
		for pos := 0; pos < len(span.Text); {
			for ws := pos; ws < len(span.Text); ws++ {
				if !strings.ContainsRune(" \t\n", rune(span.Text[ws])) {
					break
				}
				if _, err := io.WriteString(f.Writer, string(span.Text[ws])); err != nil {
					return length, err
				}
				length++
				pos++
			}
			if pos >= len(span.Text) {
				break
			}

			nextWs := strings.IndexAny(span.Text[pos:], " \t\n")
			if nextWs == -1 {
				if _, err := io.WriteString(f.Writer, span.Text[pos:]); err != nil {
					return length, err
				}
				length += len([]rune(span.Text[pos:]))
				break
			}

			word := span.Text[pos : pos+nextWs]
			wordLen := len([]rune(word))
			if length+wordLen > WrapWidth {
				modes := currentStyles.All()
				if f.ANSI && len(modes) > 0 {
					if _, err := io.WriteString(f.Writer, resetAllModes); err != nil {
						return length, err
					}
				}
				if _, err := io.WriteString(f.Writer, "\n"+followPrefix); err != nil {
					return length, err
				}
				if f.ANSI && len(modes) > 0 {
					for _, i := range modes {
						if _, err := io.WriteString(f.Writer, styleTags[i].B); err != nil {
							return length, err
						}
					}
				}
				length = len([]rune(followPrefix))
			}

			if _, err := io.WriteString(f.Writer, word); err != nil {
				return length, err
			}
			length += wordLen
			pos += nextWs
		}
	} else {
		for _, child := range span.Children {
			l, err := f.WriteSpan(child, length, followPrefix, currentStyles)
			if err != nil {
				return length, err
			}
			length = l
		}
	}

	if f.ANSI && tag.E != "" {
		if _, err := io.WriteString(f.Writer, tag.E); err != nil {
			return length, err
		}
	}

	return length, nil
}

func (f *Formatter) WriteParagraphs(paragraphs []*ftml.Paragraph, linePrefix, followPrefix string) error {
	for idx, c := range paragraphs {
		if idx > 0 {
			if _, err := io.WriteString(f.Writer, followPrefix+"\n"); err != nil {
				return err
			}
		}

		if err := f.WriteParagraph(c, linePrefix, followPrefix); err != nil {
			return err
		}

		linePrefix = followPrefix
	}

	return nil
}

func pad(s string, l int) string {
	for ; len(s) < l; s += " " {
	}
	return s
}

func (f *Formatter) WriteParagraph(p *ftml.Paragraph, linePrefix string, followPrefix string) error {
	if p.Leaf() {
		if _, err := io.WriteString(f.Writer, linePrefix); err != nil {
			return err
		}

		length := len([]rune(linePrefix))

		prev := 0
		next := 0
		boldHeaders := false
		underlineChar := ""
		prefix := ""

		switch p.Type {
		case ftml.Header1Paragraph:
			boldHeaders = true
			prev = 3
			next = 2
		case ftml.Header2Paragraph:
			underlineChar = "="
			boldHeaders = true
			prev = 2
			next = 1
		case ftml.Header3Paragraph:
			boldHeaders = true
			underlineChar = "-"
			prev = 1
		}

		if prev > f.currentLineSpacing {
			prev -= f.currentLineSpacing
		} else {
			prev = 0
		}
		f.currentLineSpacing = 0

		for i := 0; i < prev; i++ {
			if _, err := io.WriteString(f.Writer, "\n"+followPrefix); err != nil {
				return err
			}
		}

		if _, err := io.WriteString(f.Writer, prefix); err != nil {
			return err
		}
		length += len([]rune(prefix))

		for _, c := range p.Content {
			if c.Text == "\n" {
				if _, err := io.WriteString(f.Writer, " \n"+followPrefix+prefix); err != nil {
					return err
				}
				length = len([]rune(followPrefix + prefix))
				continue
			}
			var l int
			var err error

			if boldHeaders {
				c = ftml.Span{
					Style:    ftml.StyleBold,
					Children: []ftml.Span{c},
				}
			}
			if p.Type == ftml.Header1Paragraph {
				l, err = f.WriteCentered(c, length, followPrefix)
			} else {
				l, err = f.WriteSpan(c, length, followPrefix, StyleSet{})
			}
			if err != nil {
				return err
			}
			length = l
		}

		if _, err := io.WriteString(f.Writer, "\n"); err != nil {
			return err
		}

		if underlineChar != "" {
			l := len([]rune(prefix))
			for _, i := range p.Content {
				l += i.Width()
			}

			if _, err := io.WriteString(f.Writer, followPrefix); err != nil {
				return err
			}

			for i := 0; i < l; i++ {
				if _, err := io.WriteString(f.Writer, underlineChar); err != nil {
					return err
				}
			}

			if _, err := io.WriteString(f.Writer, "\n"); err != nil {
				return err
			}
		}

		for i := 0; i < next; i++ {
			if _, err := io.WriteString(f.Writer, followPrefix+"\n"); err != nil {
				return err
			}
		}

		f.currentLineSpacing = next
		return nil
	}

	if p.Type == ftml.UnorderedListParagraph {
		for idx, entry := range p.Entries {
			if idx > 0 {
				if _, err := io.WriteString(f.Writer, followPrefix+"\n"); err != nil {
					return err
				}
			}
			if err := f.WriteParagraphs(entry, linePrefix+" • ", followPrefix+"   "); err != nil {
				return err
			}
			linePrefix = followPrefix
		}

		return nil
	}

	if p.Type == ftml.OrderedListParagraph {
		digits := len(fmt.Sprintf("%d", len(p.Entries)))

		for idx, entry := range p.Entries {
			if idx > 0 {
				if _, err := io.WriteString(f.Writer, followPrefix+"\n"); err != nil {
					return err
				}
			}
			if err := f.WriteParagraphs(entry, linePrefix+pad(fmt.Sprintf("%d. ", idx+1), digits+2), followPrefix+strings.Repeat(" ", digits+2)); err != nil {
				return err
			}
			linePrefix = followPrefix
		}

		return nil
	}

	if p.Type == ftml.QuoteParagraph {
		if err := f.WriteParagraphs(p.Children, linePrefix+"| ", followPrefix+"| "); err != nil {
			return err
		}

		return nil
	}

	return f.WriteParagraphs(p.Children, linePrefix, followPrefix)
}
