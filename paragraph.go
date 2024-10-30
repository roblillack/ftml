package ftml

type ParagraphType uint8

const (
	TextParagraph ParagraphType = iota
	Header1Paragraph
	Header2Paragraph
	Header3Paragraph
	OrderedListParagraph
	UnorderedListParagraph
	QuoteParagraph
)

func (t ParagraphType) String() string {
	switch t {
	case TextParagraph:
		return "Text"
	case Header1Paragraph:
		return "Header Lvl 1"
	case Header2Paragraph:
		return "Header Lvl 2"
	case Header3Paragraph:
		return "Header Lvl 3"
	case OrderedListParagraph:
		return "Ordered List"
	case UnorderedListParagraph:
		return "Unordered List"
	case QuoteParagraph:
		return "Quote"
	}

	panic("Unknown Paragraph Type")
}

type Paragraph struct {
	Type     ParagraphType
	Children []*Paragraph
	Content  []Span
	Entries  [][]*Paragraph
}

func (p *Paragraph) Leaf() bool {
	return p.Type == TextParagraph || p.Type == Header1Paragraph || p.Type == Header2Paragraph || p.Type == Header3Paragraph
}
