package html

import (
	"strings"

	"github.com/roblillack/ftml"
)

func p__(s string) *ftml.Paragraph {
	return &ftml.Paragraph{
		Type: ftml.TextParagraph,
		Content: []ftml.Span{
			{Text: s},
		},
	}
}

func p_(content ...ftml.Span) *ftml.Paragraph {
	return &ftml.Paragraph{
		Type:    ftml.TextParagraph,
		Content: content,
	}
}

func li_(entries ...*ftml.Paragraph) []*ftml.Paragraph {
	return entries
}

func ul_(entries ...[]*ftml.Paragraph) *ftml.Paragraph {
	return &ftml.Paragraph{
		Type:    ftml.UnorderedListParagraph,
		Entries: entries,
	}
}

func ol_(entries ...[]*ftml.Paragraph) *ftml.Paragraph {
	return &ftml.Paragraph{
		Type:    ftml.OrderedListParagraph,
		Entries: entries,
	}
}

func h1_(s string) *ftml.Paragraph {
	return &ftml.Paragraph{
		Type: ftml.Header1Paragraph,
		Content: []ftml.Span{
			{Text: s},
		}}
}

func h2_(s string) *ftml.Paragraph {
	return &ftml.Paragraph{
		Type: ftml.Header2Paragraph,
		Content: []ftml.Span{
			{Text: s},
		}}
}

func h3_(s string) *ftml.Paragraph {
	return &ftml.Paragraph{
		Type: ftml.Header3Paragraph,
		Content: []ftml.Span{
			{Text: s},
		}}
}

func quote_(children ...*ftml.Paragraph) *ftml.Paragraph {
	return &ftml.Paragraph{
		Type:     ftml.QuoteParagraph,
		Children: children,
	}
}

func doc(children ...*ftml.Paragraph) *ftml.Document {
	return &ftml.Document{Paragraphs: children}
}

type ParagraphList []*ftml.Paragraph

func (l ParagraphList) String() string {
	res := make([]string, len(l))
	for i, v := range l {
		res[i] = v.Type.String()
	}
	return strings.Join(res, " --> ")
}

func span(txt string) ftml.Span {
	return ftml.Span{Text: txt}
}

func spans(txt string) []ftml.Span {
	return []ftml.Span{{Text: txt}}
}

func b_(args ...ftml.Span) ftml.Span {
	return ftml.Span{Style: ftml.StyleBold, Children: args}
}

func i_(args ...ftml.Span) ftml.Span {
	return ftml.Span{Style: ftml.StyleItalic, Children: args}
}

func b__(txt string) ftml.Span {
	return ftml.Span{Style: ftml.StyleBold, Children: spans(txt)}
}

func i__(txt string) ftml.Span {
	return ftml.Span{Style: ftml.StyleItalic, Children: spans(txt)}
}
