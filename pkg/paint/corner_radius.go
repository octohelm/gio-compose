package paint

import (
	"gioui.org/unit"
	"github.com/octohelm/gio-compose/pkg/cmp"
	"github.com/octohelm/gio-compose/pkg/layout/position"
)

type CornerRadiusSetter interface {
	SetCornerRadius(v unit.Dp, positions ...position.Position)
}

type CornerRadiusGetter interface {
	CornerRadius() CornerRadiusValues
}

type CornerRadiusValues struct {
	TopLeft     unit.Dp
	TopRight    unit.Dp
	BottomLeft  unit.Dp
	BottomRight unit.Dp
}

func (f *CornerRadiusValues) Eq(v *CornerRadiusValues) cmp.Result {
	return func() bool {
		return *f == *v
	}
}

func (r CornerRadiusValues) CornerRadius() CornerRadiusValues {
	return r
}

func (r *CornerRadiusValues) SetCornerRadius(dp unit.Dp, positions ...position.Position) {
	for _, p := range positions {
		switch p {
		case position.TopLeft:
			r.TopLeft = dp
		case position.BottomLeft:
			r.BottomLeft = dp
		case position.TopRight:
			r.TopRight = dp
		case position.BottomRight:
			r.BottomRight = dp
		case position.Top:
			r.TopLeft = dp
			r.TopRight = dp
		case position.Bottom:
			r.BottomLeft = dp
			r.BottomRight = dp
		case position.Left:
			r.TopLeft = dp
			r.BottomLeft = dp
		case position.Right:
			r.TopRight = dp
			r.BottomRight = dp
		}
	}
}
