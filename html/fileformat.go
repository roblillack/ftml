package html

import "github.com/roblillack/ftml"

const LineBreakElementName = "br"
const LineBreakElement = "<" + LineBreakElementName + " />"

// Tags that should be skipped completely when parsing the HTML,
// all other unknown tags ignored but their content is parsed.
var skipTags = map[string]struct{}{
	"title":  {},
	"style":  {},
	"script": {},
}

var paragraphElement = map[string]ftml.ParagraphType{
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
	"b":      ftml.StyleBold,
	"strong": ftml.StyleBold,
	"i":      ftml.StyleItalic,
	"em":     ftml.StyleItalic,
	"u":      ftml.StyleUnderline,
	"s":      ftml.StyleStrike,
	"mark":   ftml.StyleHighlight,
	"code":   ftml.StyleCode,
}

var blockLevelElements = map[string]struct{}{
	"p":          {},
	"h1":         {},
	"h2":         {},
	"h3":         {},
	"blockquote": {},
	"ul":         {},
	"ol":         {},
	"hr":         {},
	"div":        {},
	"tr":         {},
	"table":      {},
	"li":         {},
}
