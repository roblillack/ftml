package ftml

// The very same entity is used by HTML5 to signify unicode code point U+2005:
// A 1/4 em space which is usually around the same width as a normal ASCII space,
// is breakable and will not collapse with other spaces in HTML.
// Using this to mark breaking spaces, we're able to stay compatible with HTML5.
const NonCollapsibleSpaceEntity = "&emsp14;"

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
