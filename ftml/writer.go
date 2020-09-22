package ftml

import (
	"fmt"
	"io"
)

func Write(w io.Writer, d *Document) error {
	level := 0
	first := true

	for _, p := range d.Paragraphs {
		if first {
			first = false
		} else {
			if _, err := io.WriteString(w, "\n"); err != nil {
				return err
			}
		}
		if err := writeParagraph(w, p, level); err != nil {
			return err
		}
	}
	return nil
}

func writeIndent(w io.Writer, level int, tag string) error {
	for i := 0; i < level; i++ {
		if _, err := io.WriteString(w, "  "); err != nil {
			return err
		}
	}

	if tag != "" {
		if _, err := fmt.Fprintf(w, "%s\n", tag); err != nil {
			return err
		}
	}

	return nil
}

func writeSpan(w io.Writer, span Span) error {
	if span.Style == StyleNone {
		_, err := io.WriteString(w, span.Text)
		return err
	}

	if tag, ok := styleTags[span.Style]; ok {
		if _, err := io.WriteString(w, "<"+tag+">"); err != nil {
			return err
		}

		for _, child := range span.Children {
			if err := writeSpan(w, child); err != nil {
				return err
			}
		}

		if _, err := io.WriteString(w, "</"+tag+">"); err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("Unknown span: %s", span.String())
}

func writeLeafParagraph(w io.Writer, content []Span, tag string, level int) error {

	if err := writeIndent(w, level, ""); err != nil {
		return err
	}
	if _, err := io.WriteString(w, "<"+tag+">"); err != nil {
		return err
	}

	for _, c := range content {
		if err := writeSpan(w, c); err != nil {
			return err
		}
	}

	if _, err := io.WriteString(w, "</"+tag+">\n"); err != nil {
		return err
	}

	return nil
}

func writeParagraph(w io.Writer, p *Paragraph, level int) error {
	tag := elementTags[p.Type]
	if tag == "" {
		return fmt.Errorf("Unknown paragraph type: %s", p.Type)
	}

	if p.Leaf() {
		return writeLeafParagraph(w, p.Content, tag, level)
	}

	if err := writeIndent(w, level, "<"+tag+">"); err != nil {
		return err
	}

	if p.Type == UnorderedListParagraph || p.Type == OrderedListParagraph {
		first := true
		for _, li := range p.Entries {
			if first {
				first = false
			} else {
				if _, err := io.WriteString(w, "\n"); err != nil {
					return err
				}
			}
			if err := writeIndent(w, level+1, "<li>"); err != nil {
				return err
			}
			for _, c := range li {
				if err := writeParagraph(w, c, level+2); err != nil {
					return err
				}
			}
			if err := writeIndent(w, level+1, "</li>"); err != nil {
				return err
			}
		}
	} else {
		first := true
		for _, c := range p.Children {
			if first {
				first = false
			} else {
				if _, err := io.WriteString(w, "\n"); err != nil {
					return err
				}
			}
			if err := writeParagraph(w, c, level+1); err != nil {
				return err
			}
		}
	}

	if err := writeIndent(w, level, "</"+tag+">"); err != nil {
		return err
	}

	return nil
}
