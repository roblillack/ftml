package formatter

import (
	"github.com/roblillack/ftml"
)

type StyleSet map[ftml.InlineStyle]int

// func (s StyleSet) Has(style ftml.InlineStyle) bool {
// 	_, ok := s[style]
// 	return ok
// }

func (s StyleSet) Add(style ftml.InlineStyle) StyleSet {
	res := StyleSet{}
	for i, c := range s {
		res[i] = c
	}
	if style != ftml.StyleNone {
		res[style] = s[style] + 1
	}
	return res
}

// func (s StyleSet) Remove(style ftml.InlineStyle) {
// 	c, ok := s[style]
// 	if !ok {
// 		return
// 	}
// 	if c == 1 {
// 		delete(s, style)
// 		return
// 	}
// 	s[style] = c - 1
// }

func (s StyleSet) Empty() bool {
	return len(s) == 0
}

func (s StyleSet) All() []ftml.InlineStyle {
	all := make([]ftml.InlineStyle, 0, len(s))
	for i := range s {
		all = append(all, i)
	}
	return all
}
