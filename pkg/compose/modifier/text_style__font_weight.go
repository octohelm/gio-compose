package modifier

import (
	"github.com/octohelm/gio-compose/pkg/modifier"
	"github.com/octohelm/gio-compose/pkg/text"
	"github.com/octohelm/gio-compose/pkg/text/weight"
)

func FontWeight(w weight.Weight) modifier.Modifier[*text.Style] {
	return &fontWeightModifier{
		w: w,
	}
}

type fontWeightModifier struct {
	w weight.Weight
}

func (m *fontWeightModifier) Modify(s *text.Style) {
	s.SetFontWeight(m.w)
}
