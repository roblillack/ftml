package ftml

const LineBreakElementName = "br"
const LineBreakElement = "<" + LineBreakElementName + " />"

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
