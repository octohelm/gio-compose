package layout

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/octohelm/gio-compose/pkg/cmp"
	"github.com/octohelm/gio-compose/pkg/layout/position"
)

type EdgeInsetSetter interface {
	SetEdgeInset(dp unit.Dp, positions ...position.Position)
}

type EdgeInsetOffsetGetter interface {
	EdgeInsetOffset() Offset
}

type EdgeInset layout.Inset

func (r *EdgeInset) Eq(v *EdgeInset) cmp.Result {
	return func() bool {
		return *r == *v
	}
}

var _ EdgeInsetOffsetGetter = &EdgeInset{}

func (c *EdgeInset) EdgeInsetOffset() Offset {
	return Offset{
		X: c.Left,
		Y: c.Top,
	}
}

func (c *EdgeInset) SetEdgeInset(dp unit.Dp, positions ...position.Position) {
	for _, p := range positions {
		switch p {
		case position.Right:
			c.Right = dp
		case position.Top:
			c.Top = dp
		case position.Bottom:
			c.Bottom = dp
		case position.Left:
			c.Left = dp
		}
	}
}
