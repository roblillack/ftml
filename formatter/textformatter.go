package formatter

import (
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/roblillack/ftml"
)

const DefaultWrapWidth = 72
const DefaultQuotePrefix = "| "
const DefaultUnorderedListItemPrefix = " • "

type FormattingStyle struct {
	ResetStyles             string
	TextStyles              map[ftml.InlineStyle]StyleTags
	EscapeText              func(string) string
	QuotePrefix             string
	UnorderedListItemPrefix string
	WrapWidth               int
	LeftPadding             int
}

type StyleTags struct {
	B string
	E string
}

func DefaultFormattingStyle() FormattingStyle {
	return FormattingStyle{
		// No reset or style escape sequences defined ==> just bare text
		QuotePrefix:             DefaultQuotePrefix,
		UnorderedListItemPrefix: DefaultUnorderedListItemPrefix,
		LeftPadding:             0,
		WrapWidth:               DefaultWrapWidth,
	}
}

type Formatter struct {
	Style FormattingStyle

	writer             io.Writer
	currentLineSpacing int
}

func New(w io.Writer, f FormattingStyle) *Formatter {
	return &Formatter{
		writer: w,
		Style:  f,
	}
}

func NewASCII(w io.Writer) *Formatter {
	return New(w, DefaultFormattingStyle())
}

func NewANSI(w io.Writer) *Formatter {
	// As per ECMA-48, 5th edition:
	// https://www.ecma-international.org/wp-content/uploads/ECMA-48_5th_edition_june_1991.pdf

	return New(w, FormattingStyle{
		ResetStyles: "\033[0m",
		TextStyles: map[ftml.InlineStyle]StyleTags{
			ftml.StyleBold:      {"\033[1m", "\033[22m"},
			ftml.StyleItalic:    {"\033[3m", "\033[23m"},
			ftml.StyleHighlight: {"\033[7m", "\033[27m"},
			ftml.StyleUnderline: {"\033[4m", "\033[24m"},
			// ftml.StyleCode:      {"", ""},
			ftml.StyleStrike: {"\033[9m", "\033[29m"},
		},
		QuotePrefix:             DefaultQuotePrefix,
		UnorderedListItemPrefix: DefaultUnorderedListItemPrefix,
		LeftPadding:             0,
		WrapWidth:               DefaultWrapWidth,
	})
}

func (f *Formatter) WriteRaw(s string) error {
	_, err := io.WriteString(f.writer, s)
	return err
}

func (f *Formatter) WriteString(s string) error {
	if f.Style.EscapeText != nil {
		s = f.Style.EscapeText(s)
	}
	_, err := io.WriteString(f.writer, s)
	return err
}

func (f *Formatter) WriteDocument(d *ftml.Document) error {
	indent := strings.Repeat(" ", f.Style.LeftPadding)
	return f.WriteParagraphs(d.Paragraphs, indent, indent)
}

func (f *Formatter) WriteCentered(span ftml.Span, length int, followPrefix string) (int, error) {
	l := span.Width()

	for i := math.Floor(float64(f.Style.WrapWidth-len(followPrefix))/2 - float64(l)/2); i > 0; i-- {
		if _, err := io.WriteString(f.writer, " "); err != nil {
			return length, err
		}
		length++
	}

	return f.WriteSpan(span, length, followPrefix, StyleSet{})
}

func (f *Formatter) EmitLineBreak(followPrefix string, currentStyles StyleSet) (int, error) {
	modes := currentStyles.All()
	if len(modes) > 0 {
		if err := f.WriteRaw(f.Style.ResetStyles); err != nil {
			return 0, err
		}
	}
	if err := f.WriteRaw("\n" + followPrefix); err != nil {
		return 0, err
	}
	if len(modes) > 0 {
		for _, i := range modes {
			if err := f.WriteRaw(f.Style.TextStyles[i].B); err != nil {
				return 0, err
			}
		}
	}
	return len([]rune(followPrefix)), nil
}

func (f *Formatter) WriteSpan(span ftml.Span, length int, followPrefix string, outerStyles StyleSet) (int, error) {
	currentStyles := outerStyles.Add(span.Style)
	tag := f.Style.TextStyles[span.Style]
	if tag.B != "" {
		if err := f.WriteRaw(tag.B); err != nil {
			return 0, err
		}
	}

	if span.Style == ftml.StyleNone && len(span.Children) == 0 {
		for pos := 0; pos < len(span.Text); {
			if span.Text[pos] == '\n' {
				if _, err := f.EmitLineBreak(followPrefix, currentStyles); err != nil {
					return 0, err
				}
				length = len([]rune(followPrefix))
				pos++
				continue
			}

			// we keep the whitespace, because we might need to disregard it in
			// case we need to wrap the line
			whitespace := strings.Builder{}
			for ws := pos; ws < len(span.Text); ws++ {
				if !strings.ContainsRune(" \t\n", rune(span.Text[ws])) {
					break
				}
				whitespace.WriteByte(span.Text[ws])
				length++
				pos++
			}
			if pos >= len(span.Text) && length < f.Style.WrapWidth {
				// ok span ends here and we might put some further word on the line
				// in a later wrap? write the space out ...
				if err := f.WriteRaw(whitespace.String()); err != nil {
					return length, err
				}
				break
			}

			nextWs := strings.IndexAny(span.Text[pos:], " \t\n")
			if nextWs == -1 {
				nextWs = len(span.Text) - pos
			}

			word := span.Text[pos : pos+nextWs]
			wordLen := len([]rune(word))
			if length+wordLen > f.Style.WrapWidth {
				if l, err := f.EmitLineBreak(followPrefix, currentStyles); err != nil {
					return 0, err
				} else {
					length = l
				}
				whitespace.Reset()
			}

			if err := f.WriteRaw(whitespace.String()); err != nil {
				return length, err
			}
			if err := f.WriteString(word); err != nil {
				return length, err
			}
			length += wordLen
			pos += nextWs
		}
	} else {
		for _, child := range span.Children {
			l, err := f.WriteSpan(child, length, followPrefix, currentStyles)
			if err != nil {
				return 0, err
			}
			length = l
		}
	}

	if tag.E != "" {
		if err := f.WriteRaw(tag.E); err != nil {
			return length, err
		}
	}

	return length, nil
}

func (f *Formatter) WriteParagraphs(paragraphs []*ftml.Paragraph, linePrefix, followPrefix string) error {
	for idx, c := range paragraphs {
		if idx > 0 {
			if err := f.WriteRaw(followPrefix + "\n"); err != nil {
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
		if err := f.WriteRaw(linePrefix); err != nil {
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
			if err := f.WriteRaw("\n" + followPrefix); err != nil {
				return err
			}
		}

		if err := f.WriteRaw(prefix); err != nil {
			return err
		}
		length += len([]rune(prefix))

		for _, c := range p.Content {
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

		if err := f.WriteRaw("\n"); err != nil {
			return err
		}

		if underlineChar != "" {
			l := len([]rune(prefix))
			for _, i := range p.Content {
				l += i.Width()
			}

			if err := f.WriteRaw(followPrefix); err != nil {
				return err
			}

			for i := 0; i < l; i++ {
				if err := f.WriteRaw(underlineChar); err != nil {
					return err
				}
			}

			if err := f.WriteRaw("\n"); err != nil {
				return err
			}
		}

		for i := 0; i < next; i++ {
			if err := f.WriteRaw(followPrefix + "\n"); err != nil {
				return err
			}
		}

		f.currentLineSpacing = next
		return nil
	}

	if p.Type == ftml.UnorderedListParagraph {
		for idx, entry := range p.Entries {
			if idx > 0 {
				if err := f.WriteRaw(followPrefix + "\n"); err != nil {
					return err
				}
			}
			if err := f.WriteParagraphs(entry, linePrefix+f.Style.UnorderedListItemPrefix, followPrefix+strings.Repeat(" ", len([]rune(f.Style.UnorderedListItemPrefix)))); err != nil {
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
				if err := f.WriteRaw(followPrefix + "\n"); err != nil {
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
		if err := f.WriteParagraphs(p.Children, linePrefix+f.Style.QuotePrefix, followPrefix+f.Style.QuotePrefix); err != nil {
			return err
		}

		return nil
	}

	return f.WriteParagraphs(p.Children, linePrefix, followPrefix)
}
