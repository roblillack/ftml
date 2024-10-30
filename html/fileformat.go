package html

import "github.com/roblillack/ftml"

const LineBreakElementName = "br"
const LineBreakElement = "<" + LineBreakElementName + " />"

var skipTags = map[string]struct{}{
	"title":  {},
	"style":  {},
	"script": {},
}

var wrapperElements = map[string]ftml.ParagraphType{
	// Contains text
	"p":  ftml.TextParagraph,
	"h1": ftml.Header1Paragraph,
	"h2": ftml.Header2Paragraph,
	"h3": ftml.Header3Paragraph,
	// Contains child paragraphs
	"blockquote": ftml.QuoteParagraph,
	// Contains items each of which contain child paragraphs
	"ul": ftml.UnorderedListParagraph,
	"ol": ftml.OrderedListParagraph,
}

var inlineElements = map[string]ftml.InlineStyle{
	"b":    ftml.StyleBold,
	"i":    ftml.StyleItalic,
	"u":    ftml.StyleUnderline,
	"s":    ftml.StyleStrike,
	"mark": ftml.StyleHighlight,
	"code": ftml.StyleCode,
}

var elementTags = map[ftml.ParagraphType]string{}
var styleTags = map[ftml.InlineStyle]string{}

func init() {
	for tag, t := range wrapperElements {
		elementTags[t] = tag
	}

	for tag, s := range inlineElements {
		styleTags[s] = tag
	}
}
