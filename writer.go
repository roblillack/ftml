package ftml

import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"unicode/utf8"
)

type output struct {
	Writer      io.Writer
	Document    *Document
	Indentation string
	MaxWidth    int

	width int
	level int
}

func Write(w io.Writer, d *Document) error {
	out := &output{
		Writer:      w,
		Document:    d,
		Indentation: "  ",
		MaxWidth:    80,
	}

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
		if err := out.writeParagraph(p, level); err != nil {
			return err
		}
	}
	return nil
}

func (o *output) writeIndent(level int, tag string) error {
	for i := 0; i < level; i++ {
		if _, err := io.WriteString(o.Writer, o.Indentation); err != nil {
			return err
		}
	}

	if tag != "" {
		if _, err := fmt.Fprintf(o.Writer, "%s\n", tag); err != nil {
			return err
		}
		o.width = 0
	} else {
		o.width = utf8.RuneCountInString(o.Indentation) * level
	}

	return nil
}

var anySpace = regexp.MustCompile(`\s`)

func (o *output) Emit(txt string, level int) error {
	lines := strings.Split(txt, "\n")
	for lineIdx, line := range lines {
		for idx, i := range anySpace.Split(line, -1) {
			l := utf8.RuneCountInString(i)
			if o.width+l >= o.MaxWidth {
				if _, err := io.WriteString(o.Writer, "\n"); err != nil {
					return err
				}
				for i := 0; i < level; i++ {
					if _, err := io.WriteString(o.Writer, o.Indentation); err != nil {
						return err
					}
				}
				o.width = level * utf8.RuneCountInString(o.Indentation)
			} else if idx > 0 {
				l++
				if _, err := io.WriteString(o.Writer, " "); err != nil {
					return err
				}
			}
			if _, err := io.WriteString(o.Writer, i); err != nil {
				return err
			}
			o.width += l
		}

		if lineIdx < len(lines)-1 {
			if _, err := io.WriteString(o.Writer, "\n"); err != nil {
				return err
			}
			for i := 0; i < level; i++ {
				if _, err := io.WriteString(o.Writer, o.Indentation); err != nil {
					return err
				}
			}
			o.width = level * utf8.RuneCountInString(o.Indentation)
		}
	}

	return nil

}

func (o *output) writeSpan(span Span, level int, first, last bool) error {
	indent := ""
	for i := 0; i < level; i++ {
		indent += o.Indentation
	}

	if span.Style == StyleNone && len(span.Children) == 0 {
		txt := encodeEntities(span.Text, first, last)
		txt = strings.ReplaceAll(txt, "\n", LineBreakElement+"\n")
		// if last {
		// 	txt = strings.ReplaceAll(txt, "\n", LineBreakElement+"\n")
		// } else {
		// 	txt = strings.ReplaceAll(txt, "\n", LineBreakElement+"\n")
		// }
		// txt, o.width = encodeLineBreaks(txt, level, o.Indentation, last, o.width)
		return o.Emit(txt, level)
	}

	if tag, ok := styleTags[span.Style]; ok {
		if err := o.Emit("<"+tag+">", level); err != nil {
			return err
		}

		for _, child := range span.Children {
			if err := o.writeSpan(child, level, false, false); err != nil {
				return err
			}
		}

		if err := o.Emit("</"+tag+">", level); err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("Unknown span: %s", span.String())
}

func simpleWriteSpan(o io.StringWriter, span Span, first, last bool) error {
	if len(span.Children) > 0 {
		tag, ok := styleTags[span.Style]

		if ok {
			if _, err := o.WriteString("<" + tag + ">"); err != nil {
				return err
			}
		}
		for _, child := range span.Children {
			if err := simpleWriteSpan(o, child, false, false); err != nil {
				return err
			}
		}

		if ok {
			if _, err := o.WriteString("</" + tag + ">"); err != nil {
				return err
			}
		}
		return nil
	}

	txt := encodeEntities(span.Text, first, last)
	txt = strings.ReplaceAll(txt, "\n", LineBreakElement+"\n")
	_, err := o.WriteString(txt)
	return err
}

func (o *output) writeLeafParagraph(content []Span, tag string, level int) error {

	// Let's try outputting in one line first
	b := &strings.Builder{}
	for i := 0; i < level; i++ {
		if _, err := b.WriteString(o.Indentation); err != nil {
			return err
		}
	}
	if _, err := b.WriteString("<" + tag + ">"); err != nil {
		return err
	}
	for idx, c := range content {
		if err := simpleWriteSpan(b, c, idx == 0, idx == len(content)-1); err != nil {
			return err
		}
	}
	if _, err := b.WriteString("</" + tag + ">\n"); err != nil {
		return err
	}
	line := b.String()
	if utf8.RuneCountInString(line) <= o.MaxWidth && !strings.Contains(line[0:len(line)-1], "\n") {
		_, err := io.WriteString(o.Writer, line)
		return err
	}

	// Ok, tags and content on different lines ...
	if err := o.writeIndent(level, ""); err != nil {
		return err
	}
	if _, err := io.WriteString(o.Writer, "<"+tag+">\n"); err != nil {
		return err
	}

	if err := o.writeIndent(level+1, ""); err != nil {
		return err
	}

	for idx, c := range content {
		if err := o.writeSpan(c, level+1, idx == 0, idx == len(content)-1); err != nil {
			return err
		}
	}

	if _, err := io.WriteString(o.Writer, "\n"); err != nil {
		return err
	}
	if err := o.writeIndent(level, ""); err != nil {
		return err
	}
	if _, err := io.WriteString(o.Writer, "</"+tag+">\n"); err != nil {
		return err
	}

	return nil
}

func (o *output) writeParagraph(p *Paragraph, level int) error {
	tag := elementTags[p.Type]
	if tag == "" {
		return fmt.Errorf("Unknown paragraph type: %s", p.Type)
	}

	if p.Leaf() {
		return o.writeLeafParagraph(p.Content, tag, level)
	}

	if err := o.writeIndent(level, "<"+tag+">"); err != nil {
		return err
	}

	if p.Type == UnorderedListParagraph || p.Type == OrderedListParagraph {
		first := true
		for _, li := range p.Entries {
			if first {
				first = false
			} else {
				if _, err := io.WriteString(o.Writer, "\n"); err != nil {
					return err
				}
			}
			if err := o.writeIndent(level+1, "<li>"); err != nil {
				return err
			}
			for _, c := range li {
				if err := o.writeParagraph(c, level+2); err != nil {
					return err
				}
			}
			if err := o.writeIndent(level+1, "</li>"); err != nil {
				return err
			}
		}
	} else {
		first := true
		for _, c := range p.Children {
			if first {
				first = false
			} else {
				if _, err := io.WriteString(o.Writer, "\n"); err != nil {
					return err
				}
			}
			if err := o.writeParagraph(c, level+1); err != nil {
				return err
			}
		}
	}

	if err := o.writeIndent(level, "</"+tag+">"); err != nil {
		return err
	}

	return nil
}
