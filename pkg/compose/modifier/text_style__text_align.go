package modifier

import (
	"github.com/octohelm/gio-compose/pkg/modifier"
	"github.com/octohelm/gio-compose/pkg/text"
)

func TextAlign(align text.Alignment) modifier.Modifier[*text.Style] {
	return &textAlignModifier{
		align: align,
	}
}

type textAlignModifier struct {
	align text.Alignment
}

func (t textAlignModifier) Modify(w *text.Style) {
	w.SetTextAlign(t.align)
}
