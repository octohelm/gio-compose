package modifier

import (
	"github.com/octohelm/gio-compose/pkg/modifier"
	"github.com/octohelm/gio-compose/pkg/text"
)

func FontFamily(typeface string) modifier.Modifier[*text.Style] {
	return &fontFamilyModifier{
		typeface: typeface,
	}
}

type fontFamilyModifier struct {
	typeface string
}

func (m fontFamilyModifier) Modify(s *text.Style) {
	s.SetFontFamily(m.typeface)
}
