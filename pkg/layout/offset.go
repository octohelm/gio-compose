package layout

import (
	"gioui.org/unit"
	"github.com/octohelm/gio-compose/pkg/cmp"
	"github.com/octohelm/gio-compose/pkg/layout/direction"
)

type OffsetSetter interface {
	SetOffset(dp unit.Dp, directions ...direction.Direction)
}

var _ OffsetSetter = &Offset{}

type Offset struct {
	X unit.Dp
	Y unit.Dp
}

func (s *Offset) Eq(v *Offset) cmp.Result {
	return func() bool {
		return *s == *v
	}
}

func (c *Offset) IsZero() bool {
	return c.X == 0 && c.Y == 0
}

func (c *Offset) SetOffset(dp unit.Dp, directions ...direction.Direction) {
	for _, p := range directions {
		switch p {
		case direction.X:
			c.X = dp
		case direction.Y:
			c.Y = dp
		}
	}
}
