package layout

import (
	"gioui.org/unit"
	"github.com/octohelm/gio-compose/pkg/cmp"
)

type SpacingSetter interface {
	SetSpacing(dp unit.Dp)
}

var _ SpacingSetter = &Spacer{}

type Spacer struct {
	Spacing unit.Dp
}

func (s *Spacer) Eq(v *Spacer) cmp.Result {
	return func() bool {
		return s.Spacing == v.Spacing
	}
}

func (s *Spacer) SetSpacing(dp unit.Dp) {
	s.Spacing = dp
}
