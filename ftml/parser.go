package ftml

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
	"unicode"

	"github.com/roblillack/gockl"
)

type parser struct {
	Tokenizer     *gockl.Tokenizer
	Doc           *Document
	Breadcrumbs   []*Paragraph
	ListItemLevel int
	StyleStack    []InlineStyle
}

func (p *parser) down(paraType ParagraphType) error {
	para := &Paragraph{Type: paraType}

	parent := p.parent()
	if parent != nil && parent.Leaf() {
		return fmt.Errorf("paragraphs not allowed inside leaf paragraph nodes when trying to add %s below %s", paraType, parent.Type)
	}

	if parent == nil {
		p.Doc.Paragraphs = append(p.Doc.Paragraphs, para)
	} else {
		if parent.Type == UnorderedListParagraph || parent.Type == OrderedListParagraph {
			pos := len(parent.Entries) - 1
			if pos < 0 {
				return fmt.Errorf("paragraph content for list without list item")
			}
			parent.Entries[pos] = append(parent.Entries[pos], para)
		} else {
			parent.Children = append(parent.Children, para)
		}
	}

	if !para.Leaf() {
		p.Breadcrumbs = append(p.Breadcrumbs, para)
		return nil
	}

	content, err := readContent(p.Tokenizer, elementTags[paraType])
	if err != nil {
		return err
	}
	para.Content = content

	return nil
}

func (p *parser) up(t ParagraphType) error {
	current := p.parent()
	if current == nil {
		return fmt.Errorf("closing unopened paragraph of type %v", t)
	}
	if current.Type != t {
		return fmt.Errorf("cannot close %v with %v", current.Type, t)
	}

	p.Breadcrumbs = p.Breadcrumbs[0 : len(p.Breadcrumbs)-1]
	return nil
}

func (p *parser) parent() *Paragraph {
	if len(p.Breadcrumbs) == 0 {
		return nil
	}

	return p.Breadcrumbs[len(p.Breadcrumbs)-1]
}

func (p *parser) ProcessToken(token gockl.Token) error {
	if t, ok := token.(gockl.StartElementToken); ok {
		if t.Name() == "li" {
			parent := p.parent()
			if parent == nil ||
				(parent.Type != UnorderedListParagraph && parent.Type != OrderedListParagraph) {
				return fmt.Errorf("unexpected list item, parent: %v", parent)
			}
			parent.Entries = append(parent.Entries, []*Paragraph{})
			p.ListItemLevel++
			return nil
		}

		if paraType, ok := wrapperElements[t.Name()]; ok {
			return p.down(paraType)
		}
	} else if t, ok := token.(gockl.EndElementToken); ok {
		if t.Name() == "li" {
			if p.ListItemLevel < 1 {
				return fmt.Errorf("unexpected closing tag for list item")
			}
			p.ListItemLevel--
			return nil
		}

		if paraType, ok := wrapperElements[t.Name()]; ok {
			return p.up(paraType)
		}
	} else if t, ok := token.(gockl.TextToken); ok {
		txt := strings.TrimSpace(t.Raw())
		if txt == "" {
			return nil
		}
		return fmt.Errorf("unexpected text content: %s", txt)
	}

	return nil
}

var space = regexp.MustCompile(`\s+`)

func collapseWhitespace(s string, first, last bool) string {
	if first {
		s = strings.TrimLeftFunc(s, unicode.IsSpace)
	}
	if last {
		s = strings.TrimRightFunc(s, unicode.IsSpace)
	}
	return space.ReplaceAllString(s, " ")
}

func decodeEntities(s string) string {
	return strings.ReplaceAll(s, NonCollapsibleSpaceEntity, " ")
}

func readText(z *gockl.Tokenizer) (string, gockl.Token, error) {
	res := ""

	for {
		token, err := z.Next()
		if err != nil || token == nil {
			if err == io.EOF {
				err = io.ErrUnexpectedEOF
			}
			return res, token, err
		}

		if t, ok := token.(gockl.TextToken); ok {
			res += t.Raw()
		} else {
			return res, token, nil
		}
	}
}

func readSpan(z *gockl.Tokenizer, style InlineStyle) (Span, error) {
	res := Span{Style: style, Children: []Span{}}

	for {
		str, token, err := readText(z)
		if err != nil {
			return res, err
		}
		if token == nil {
			return res, fmt.Errorf("no closing tag for %s", style)
		}
		str = decodeEntities(collapseWhitespace(str, false, false))
		if str != "" {
			res.Children = append(res.Children, Span{Text: str})
		}

		if t, ok := token.(gockl.StartElementToken); ok {
			if st, ok := inlineElements[t.Name()]; ok {
				span, err := readSpan(z, st)
				if err != nil {
					return res, err
				}
				res.Children = append(res.Children, span)
				continue
			}
		}

		if t, ok := token.(gockl.EndElementToken); ok && inlineElements[t.Name()] == style {
			return res, nil
		}

		return res, fmt.Errorf("unexpected token: %v", token)
	}
}

func trimWhiteSpace(spans []Span) []Span {
	res := spans

	for idx, i := range res {
		if i.Style != StyleNone {
			continue
		}
		txt := strings.TrimLeftFunc(i.Text, unicode.IsSpace)
		if txt == "" {
			n := append([]Span{}, res[0:idx]...)
			res = append(n, res[idx+1:]...)
		} else if txt != i.Text {
			i.Text = txt
			res[idx] = i
		}
		break
	}

	for idx := len(res) - 1; idx >= 0; idx-- {
		i := res[idx]
		if i.Style != StyleNone {
			continue
		}
		txt := strings.TrimRightFunc(i.Text, unicode.IsSpace)
		if txt == "" {
			n := append([]Span{}, res[0:idx]...)
			res = append(n, res[idx+1:]...)
		} else if txt != i.Text {
			i.Text = txt
			res[idx] = i
		}
		break
	}

	return res
}

type bufferedSpanList struct {
	Spans   []Span
	First   bool
	TrimEnd bool
	Buffer  string
}

func newBufferedSpanList() *bufferedSpanList {
	return &bufferedSpanList{
		Spans:   []Span{},
		First:   true,
		TrimEnd: false,
		Buffer:  "",
	}
}

func (b *bufferedSpanList) flush() {
	if b.Buffer != "" {
		b.Spans = append(b.Spans, Span{
			Text: decodeEntities(collapseWhitespace(b.Buffer, b.First, b.TrimEnd)),
		})
		b.Buffer = ""
		b.First = false
	}
}

func (b *bufferedSpanList) AddLineBreak() {
	b.TrimEnd = true
	b.flush()
	b.Spans = append(b.Spans, Span{Text: "\n"})
	b.First = true
}

func (b *bufferedSpanList) Add(span Span) {
	b.flush()
	b.First = false
	b.TrimEnd = false
	b.Spans = append(b.Spans, span)
}

func (b *bufferedSpanList) AddText(txt string) {
	b.TrimEnd = false
	b.flush()
	b.Buffer = txt
}

func (b *bufferedSpanList) Close() []Span {
	b.TrimEnd = true
	b.flush()
	return b.Spans
}

func readContent(z *gockl.Tokenizer, endToken string) ([]Span, error) {
	res := newBufferedSpanList()

	for {
		str, token, err := readText(z)
		if str != "" {
			res.AddText(str)
		}
		if err != nil || token == nil {
			return res.Close(), err
		}

		if t, ok := token.(gockl.EndElementToken); ok && t.Name() == endToken {
			return res.Close(), nil
		}

		if t, ok := token.(gockl.EmptyElementToken); ok && t.Name() == LineBreakElementName {
			res.AddLineBreak()
			continue
		}

		t, ok := token.(gockl.StartElementToken)
		if !ok {
			return res.Close(), fmt.Errorf("unexpected token type: %v", token)
		}

		st, ok := inlineElements[t.Name()]
		if !ok {
			return res.Close(), fmt.Errorf("non-inline token: %v", token)
		}

		span, err := readSpan(z, st)
		if err != nil {
			return res.Close(), err
		}

		res.Add(span)
	}
}

func Parse(r io.Reader) (*Document, error) {

	buf := &bytes.Buffer{}
	buf.ReadFrom(r)
	p := parser{Doc: &Document{}, Tokenizer: gockl.New(buf.String())}

	for {
		t, err := p.Tokenizer.Next()
		if t == nil || err != nil {
			break
		}
		if err := p.ProcessToken(t); err != nil {
			return nil, err
		}
	}

	return p.Doc, nil
}
