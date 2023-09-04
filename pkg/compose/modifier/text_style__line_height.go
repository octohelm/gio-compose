package modifier

import (
	"github.com/octohelm/gio-compose/pkg/modifier"
	"github.com/octohelm/gio-compose/pkg/text"
	"github.com/octohelm/gio-compose/pkg/unit"
)

func LineHeight(s unit.Sp) modifier.Modifier[*text.Style] {
	return &lineHeightModifier{
		s: s,
	}
}

type lineHeightModifier struct {
	s unit.Sp
}

func (m *lineHeightModifier) Modify(s *text.Style) {
	s.SetLineHeight(m.s)
}
