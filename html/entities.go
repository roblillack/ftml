package html

import (
	"regexp"
	"strings"
)

// The very same entity is used by HTML5 to signify unicode code point U+2005:
// A 1/4 em space which is usually around the same width as a normal ASCII space,
// is breakable and will not collapse with other spaces in HTML.
// Using this to mark breaking spaces, we're able to stay compatible with HTML5.
const NonCollapsibleSpaceEntity = "&emsp14;"

const LessThanSymbolEntity = "&lt;"

func decodeEntities(s string) string {
	s = strings.ReplaceAll(s, NonCollapsibleSpaceEntity, " ")
	s = strings.ReplaceAll(s, LessThanSymbolEntity, "<")
	return s
}

func replaceSpaces(s string) string {
	b := strings.Builder{}
	for i := 0; i < len(s); i++ {
		b.WriteString(NonCollapsibleSpaceEntity)
	}
	return b.String()
}

func replaceLeadingSpaces(s string) string {
	b := strings.Builder{}
	for i := 0; i < len(s)-1; i++ {
		b.WriteString(NonCollapsibleSpaceEntity)
	}
	b.WriteString(s[len(s)-1:])
	return b.String()
}

func replaceTrailingSpaces(s string) string {
	b := strings.Builder{}
	b.WriteString(s[0:1])
	for i := 0; i < len(s)-1; i++ {
		b.WriteString(NonCollapsibleSpaceEntity)
	}
	return b.String()
}

var multipleSpaces = regexp.MustCompile(`  +`)
var trailingSpaces = regexp.MustCompile(`\s +`)
var leadingSpaces = regexp.MustCompile(` +\s`)
var spacesAtStart = regexp.MustCompile(`^ +`)
var spacesAtEnd = regexp.MustCompile(` +$`)

func encodeEntities(s string, first, last bool) string {
	if first {
		s = spacesAtStart.ReplaceAllStringFunc(s, replaceSpaces)
	}
	if last {
		s = spacesAtEnd.ReplaceAllStringFunc(s, replaceSpaces)
	}
	s = multipleSpaces.ReplaceAllStringFunc(s, replaceSpaces)
	s = trailingSpaces.ReplaceAllStringFunc(s, replaceTrailingSpaces)
	s = leadingSpaces.ReplaceAllStringFunc(s, replaceLeadingSpaces)

	s = strings.ReplaceAll(s, "<", LessThanSymbolEntity)

	return s
}
