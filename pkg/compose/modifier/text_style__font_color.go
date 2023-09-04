package modifier

import (
	"image/color"

	"github.com/octohelm/gio-compose/pkg/modifier"

	"github.com/octohelm/gio-compose/pkg/text"
)

func TextColor(c color.Color) modifier.Modifier[*text.Style] {
	return &textColorModifier{
		c: c,
	}
}

type textColorModifier struct {
	c color.Color
}

func (m textColorModifier) Modify(s *text.Style) {
	s.SetTextColor(m.c)
}
