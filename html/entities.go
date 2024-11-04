package html

import (
	"strings"
)

// The very same entity is used by HTML5 to signify unicode code point U+2005:
// A 1/4 em space which is usually around the same width as a normal ASCII space,
// is breakable and will not collapse with other spaces in HTML.
// Using this to mark breaking spaces, we're able to stay compatible with HTML5.
const NonCollapsibleSpaceEntity = "&emsp14;"

func decodeEntities(s string) string {
	s = strings.ReplaceAll(s, NonCollapsibleSpaceEntity, " ")
	s = strings.ReplaceAll(s, "&nbsp;", " ")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&quot;", "\"")
	s = strings.ReplaceAll(s, "&apos;", "'")

	return s
}
