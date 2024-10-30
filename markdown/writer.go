package markdown

import (
	"fmt"
	"io"
	"strings"

	"github.com/roblillack/ftml"
)

func Write(w io.Writer, d *ftml.Document) error {
	return writeParagraphs(w, d.Paragraphs, "", "")
}

var styleTags = map[ftml.InlineStyle]struct{ B, E string }{
	ftml.StyleBold:      {"**", "**"},
	ftml.StyleItalic:    {"_", "_"},
	ftml.StyleHighlight: {"<mark>", "</mark>"},
	ftml.StyleUnderline: {"<u>", "</u>"},
	ftml.StyleCode:      {"`", "`"},
	ftml.StyleStrike:    {"~~", "~~"},
}

func writeSpan(w io.Writer, span ftml.Span) error {
	tag := styleTags[span.Style]
	if tag.B != "" {
		if _, err := io.WriteString(w, tag.B); err != nil {
			return err
		}
	}

	if span.Style == ftml.StyleNone {
		if _, err := io.WriteString(w, span.Text); err != nil {
			return err
		}
	} else {
		for _, child := range span.Children {
			if err := writeSpan(w, child); err != nil {
				return err
			}
		}
	}

	if tag.E != "" {
		if _, err := io.WriteString(w, tag.E); err != nil {
			return err
		}
	}

	return nil
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

		for _, c := range p.Content {
			if c.Text == "\n" {
				if _, err := io.WriteString(w, " \\\n"+prefix); err != nil {
					return err
				}
			}
			if err := writeSpan(w, c); err != nil {
				return err
			}
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
			if err := writeParagraphs(w, entry, linePrefix+"- ", followPrefix+"  "); err != nil {
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
		if err := writeParagraphs(w, p.Children, linePrefix+"> ", followPrefix+"> "); err != nil {
			return err
		}

		return nil
	}

	return writeParagraphs(w, p.Children, linePrefix, followPrefix)
}
