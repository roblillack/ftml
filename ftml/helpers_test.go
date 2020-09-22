package ftml

import "strings"

func p__(s string) *Paragraph {
	return &Paragraph{
		Type: TextParagraph,
		Content: []Span{
			{Text: s},
		},
	}
}

func p_(content ...Span) *Paragraph {
	return &Paragraph{
		Type:    TextParagraph,
		Content: content,
	}
}

func li_(entries ...*Paragraph) []*Paragraph {
	return entries
}

func ul_(entries ...[]*Paragraph) *Paragraph {
	return &Paragraph{
		Type:    UnorderedListParagraph,
		Entries: entries,
	}
}

func ol_(entries ...[]*Paragraph) *Paragraph {
	return &Paragraph{
		Type:    OrderedListParagraph,
		Entries: entries,
	}
}

func h1_(s string) *Paragraph {
	return &Paragraph{
		Type: Header1Paragraph,
		Content: []Span{
			{Text: s},
		}}
}

func h2_(s string) *Paragraph {
	return &Paragraph{
		Type: Header2Paragraph,
		Content: []Span{
			{Text: s},
		}}
}

func h3_(s string) *Paragraph {
	return &Paragraph{
		Type: Header3Paragraph,
		Content: []Span{
			{Text: s},
		}}
}

func quote_(children ...*Paragraph) *Paragraph {
	return &Paragraph{
		Type:     QuoteParagraph,
		Children: children,
	}
}

func doc(children ...*Paragraph) *Document {
	return &Document{Paragraphs: children}
}

type ParagraphList []*Paragraph

func (l ParagraphList) String() string {
	res := make([]string, len(l))
	for i, v := range l {
		res[i] = v.Type.String()
	}
	return strings.Join(res, " --> ")
}

func span(txt string) Span {
	return Span{Text: txt}
}

func spans(txt string) []Span {
	return []Span{{Text: txt}}
}

func b_(args ...Span) Span {
	return Span{Style: StyleBold, Children: args}
}

func i_(args ...Span) Span {
	return Span{Style: StyleItalic, Children: args}
}

func b__(txt string) Span {
	return Span{Style: StyleBold, Children: spans(txt)}
}

func i__(txt string) Span {
	return Span{Style: StyleItalic, Children: spans(txt)}
}
