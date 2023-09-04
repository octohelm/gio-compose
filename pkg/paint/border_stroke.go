package paint

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"github.com/octohelm/gio-compose/pkg/cmp"
	"github.com/octohelm/gio-compose/pkg/layout/position"
	"github.com/octohelm/gio-compose/pkg/paint/canvas"
)

type BorderStrokeSetter interface {
	SetBorderStroke(styles ...StrokeStyleWithPosition)
}

var _ BorderStrokeSetter = &BorderStroke{}

type BorderStroke struct {
	Top    StrokeStyle
	Right  StrokeStyle
	Bottom StrokeStyle
	Left   StrokeStyle
}

func (f *BorderStroke) Eq(v *BorderStroke) cmp.Result {
	return func() bool {
		return *f == *v
	}
}

func (f *BorderStroke) SetBorderStroke(styles ...StrokeStyleWithPosition) {
	for _, s := range styles {
		switch s.Position {
		case position.Top:
			f.Top = s.StrokeStyle
		case position.Right:
			f.Right = s.StrokeStyle
		case position.Bottom:
			f.Bottom = s.StrokeStyle
		case position.Left:
			f.Left = s.StrokeStyle
		}
	}
}

func (f *BorderStroke) PaintShape(gtx layout.Context, s SizedShape) {
	styles := []StrokeStyleWithPosition{
		{Position: position.Top, StrokeStyle: f.Top},
		{Position: position.Right, StrokeStyle: f.Right},
		{Position: position.Bottom, StrokeStyle: f.Bottom},
		{Position: position.Left, StrokeStyle: f.Left},
	}

	positions := make([]position.Position, 0, 4)
	style := StrokeStyle{}

	tryDraw := func() {
		if style.Width != 0 && len(positions) > 0 {
			p := s.Path(gtx, positions...)

			path := canvas.PathSpec(
				gtx.Ops,
				p.Stroke(float64(gtx.Dp(style.Width)), canvas.ButtCap, canvas.MiterJoin, 0.1))

			paint.FillShape(gtx.Ops, style.Color, ClipOutlineOf(path).Op(gtx.Ops))

			style = StrokeStyle{}
			positions = make([]position.Position, 0, 4)
		}
	}

	for _, s := range styles {
		if s.Width == 0 {
			tryDraw()
			continue
		}

		if style.Width != 0 {
			if s.StrokeStyle != style {
				tryDraw()
				continue
			}
		}

		positions = append(positions, s.Position)
		style = s.StrokeStyle
	}

	tryDraw()
}

type StrokeStyle struct {
	Width unit.Dp
	Color color.NRGBA
}

type StrokeStyleWithPosition struct {
	Position position.Position
	StrokeStyle
}
