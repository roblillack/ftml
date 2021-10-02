package formatter

import (
	"fmt"
	"io"
	"strings"

	"github.com/roblillack/pure/ftml"
)

type Breadcrumbs []ftml.ParagraphType

func Write(w io.Writer, d *ftml.Document) error {
	return writeParagraphs(w, d.Paragraphs, "", "")
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

func writeSpan(w io.Writer, span ftml.Span, length int, followPrefix string) (int, error) {
	tag := styleTags[span.Style]
	if tag.B != "" {
		if _, err := io.WriteString(w, tag.B); err != nil {
			return length, err
		}
	}

	if span.Style == ftml.StyleNone {
		for pos := 0; pos < len(span.Text); {
			for ws := pos; ws < len(span.Text); ws++ {
				if !strings.ContainsRune(" \t\n", rune(span.Text[ws])) {
					break
				}
				if _, err := io.WriteString(w, string(span.Text[ws])); err != nil {
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
				if _, err := io.WriteString(w, span.Text[pos:]); err != nil {
					return length, err
				}
				length += len([]rune(span.Text[pos:]))
				break
			}

			word := span.Text[pos : pos+nextWs]
			wordLen := len([]rune(word))
			if length+wordLen > WrapWidth {
				if _, err := io.WriteString(w, "\n"+followPrefix); err != nil {
					return length, err
				}
				length = len([]rune(followPrefix))
			}

			if _, err := io.WriteString(w, word); err != nil {
				return length, err
			}
			length += wordLen
			pos += nextWs
		}

		// words := strings.Fields(span.Text)
		// for _, word := range words {
		// 	wordLen := len([]rune(word))
		// 	if wordLen+length > WrapWidth {
		// 		if _, err := io.WriteString(w, "\n"+followPrefix); err != nil {
		// 			return length, err
		// 		}
		// 		length = len([]rune(followPrefix))
		// 	}
		// 	if _, err := io.WriteString(w, word+" "); err != nil {
		// 		return length, err
		// 	}
		// 	length += wordLen + 1
		// }
	} else {
		for _, child := range span.Children {
			l, err := writeSpan(w, child, length, followPrefix)
			if err != nil {
				return length, err
			}
			length = l
		}
	}

	if tag.E != "" {
		if _, err := io.WriteString(w, tag.E); err != nil {
			return length, err
		}
	}

	return length, nil
}

func writeParagraphs(w io.Writer, paragraphs []*ftml.Paragraph, linePrefix, followPrefix string) error {
	for idx, c := range paragraphs {
		if idx > 0 {
			if _, err := io.WriteString(w, followPrefix+"\n"); err != nil {
				return err
			}
		}

		if err := writeParagraph(w, c, linePrefix, followPrefix); err != nil {
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

func writeParagraph(w io.Writer, p *ftml.Paragraph, linePrefix string, followPrefix string) error {
	if p.Leaf() {
		if _, err := io.WriteString(w, linePrefix); err != nil {
			return err
		}

		length := len([]rune(linePrefix))

		prefix := ""
		switch p.Type {
		case ftml.Header1Paragraph:
			prefix = "# "
		case ftml.Header2Paragraph:
			prefix = "## "
		case ftml.Header3Paragraph:
			prefix = "### "
		}
		if _, err := io.WriteString(w, prefix); err != nil {
			return err
		}
		length += len([]rune(prefix))

		for _, c := range p.Content {
			if c.Text == "\n" {
				if _, err := io.WriteString(w, " \n"+followPrefix+prefix); err != nil {
					return err
				}
				length = len([]rune(followPrefix + prefix))
				continue
			}
			l, err := writeSpan(w, c, length, followPrefix)
			if err != nil {
				return err
			}
			length = l
		}
		if _, err := io.WriteString(w, "\n"); err != nil {
			return err
		}

		return nil
	}

	if p.Type == ftml.UnorderedListParagraph {
		for idx, entry := range p.Entries {
			if idx > 0 {
				if _, err := io.WriteString(w, followPrefix+"\n"); err != nil {
					return err
				}
			}
			if err := writeParagraphs(w, entry, linePrefix+" - ", followPrefix+"   "); err != nil {
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
				if _, err := io.WriteString(w, followPrefix+"\n"); err != nil {
					return err
				}
			}
			if err := writeParagraphs(w, entry, linePrefix+pad(fmt.Sprintf("%d. ", idx+1), digits+2), followPrefix+strings.Repeat(" ", digits+2)); err != nil {
				return err
			}
			linePrefix = followPrefix
		}

		return nil
	}

	if p.Type == ftml.QuoteParagraph {
		if err := writeParagraphs(w, p.Children, linePrefix+"| ", followPrefix+"| "); err != nil {
			return err
		}

		return nil
	}

	return writeParagraphs(w, p.Children, linePrefix, followPrefix)
}
