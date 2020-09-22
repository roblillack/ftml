package ftml

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/roblillack/gockl"
)

type parser struct {
	Tokenizer     *gockl.Tokenizer
	Doc           *Document
	Breadcrumbs   []*Paragraph
	ListItemLevel int
	StyleStack    []InlineStyle
}

var wrapperElements = map[string]ParagraphType{
	// Contains text
	"p":  TextParagraph,
	"h1": Header1Paragraph,
	"h2": Header2Paragraph,
	"h3": Header3Paragraph,
	// Contains child paragraphs
	"blockquote": QuoteParagraph,
	// Contains items each of which contain child paragraphs
	"ul": UnorderedListParagraph,
	"ol": OrderedListParagraph,
}

var inlineElements = map[string]InlineStyle{
	"b":    StyleBold,
	"i":    StyleItalic,
	"u":    StyleUnderline,
	"s":    StyleStrike,
	"mark": StyleHighlight,
	"code": StyleCode,
}

var elementTags = map[ParagraphType]string{}
var styleTags = map[InlineStyle]string{}

func init() {
	for tag, t := range wrapperElements {
		elementTags[t] = tag
	}

	for tag, s := range inlineElements {
		styleTags[s] = tag
	}
}

func (p *parser) down(paraType ParagraphType) error {
	para := &Paragraph{Type: paraType}

	parent := p.parent()
	if parent != nil && parent.Leaf() {
		return fmt.Errorf("paragraphs not allowed inside leaf paragraph nodes whenn trying to add %s below %s", paraType, parent.Type)
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

func readContent(z *gockl.Tokenizer, endToken string) ([]Span, error) {
	res := []Span{}

	for {
		str, token, err := readText(z)
		if str != "" {
			res = append(res, Span{Text: str})
		}
		if err != nil || token == nil {
			return res, err
		}

		if t, ok := token.(gockl.EndElementToken); ok && t.Name() == endToken {
			return res, nil
		}

		t, ok := token.(gockl.StartElementToken)
		if !ok {
			return res, fmt.Errorf("unexpected token type: %v", token)
		}

		st, ok := inlineElements[t.Name()]
		if !ok {
			return res, fmt.Errorf("non-inline token: %v", token)
		}

		span, err := readSpan(z, st)
		if err != nil {
			return res, err
		}

		res = append(res, span)
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
