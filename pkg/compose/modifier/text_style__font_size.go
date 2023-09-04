package modifier

import (
	"github.com/octohelm/gio-compose/pkg/modifier"
	"github.com/octohelm/gio-compose/pkg/text"
	"github.com/octohelm/gio-compose/pkg/unit"
)

func FontSize(s unit.Sp) modifier.Modifier[*text.Style] {
	return &fontSizeModifier{
		s: s,
	}
}

type fontSizeModifier struct {
	s unit.Sp
}

func (m *fontSizeModifier) Modify(s *text.Style) {
	s.SetFontSize(m.s)
}
