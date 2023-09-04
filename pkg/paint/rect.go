package paint

import (
	"image"
	"math"
	"slices"

	"github.com/octohelm/gio-compose/pkg/layout/position"

	"github.com/octohelm/gio-compose/pkg/paint/canvas"

	"github.com/octohelm/gio-compose/pkg/cmp"
	"github.com/octohelm/gio-compose/pkg/paint/size"

	"gioui.org/layout"
	"gioui.org/unit"
)

func NewRect() *Rect {
	return &Rect{
		Width:  0,
		Height: 0,
	}
}

type Rect struct {
	Width  unit.Dp
	Height unit.Dp
}

func (r *Rect) Equals(v any) cmp.Result {
	return cmp.Cast(v, func(v *Rect) cmp.Result {
		return r.Eq(v)
	})
}

func (r *Rect) Eq(v *Rect) cmp.Result {
	return func() bool {
		return *r == *v
	}
}

func (r *Rect) Sized(t size.Type) size.SizingType {
	if t == size.Width {
		if r.Width > 0 {
			return size.Exactly
		}
		if r.Width < 0 {
			return size.MatchParent
		}
		return size.WrapContent
	}
	if r.Height > 0 {
		return size.Exactly
	}
	if r.Height < 0 {
		return size.MatchParent
	}
	return size.WrapContent
}

var _ SizeSetter = &Rect{}

func (r *Rect) SetSize(dp unit.Dp, types ...size.Type) {
	for _, s := range types {
		switch s {
		case size.Width:
			r.Width = dp
		case size.Height:
			r.Height = dp
		}
	}
}

func (r *Rect) Rectangle(gtx layout.Context) image.Rectangle {
	s := image.Point{}

	if r.Width > 0 {
		s.X = gtx.Dp(r.Width)
	} else {
		factor := math.Abs(float64(r.Width))

		if factor > 0 && factor <= 1 {
			s.X = int(math.Round(float64(gtx.Constraints.Max.X) * factor))
		} else {
			s.X = gtx.Constraints.Min.X
		}
	}

	if r.Height > 0 {
		s.Y = gtx.Dp(r.Height)
	} else {
		factor := math.Abs(float64(r.Height))
		if factor > 0 && factor <= 1 {
			s.Y = int(math.Round(float64(gtx.Constraints.Max.Y) * factor))
		} else {
			s.Y = gtx.Constraints.Min.Y
		}
	}

	return image.Rectangle{
		Max: s,
	}
}

func (r *Rect) Path(gtx layout.Context, positions ...position.Position) *canvas.Path {
	if len(positions) == 0 {
		positions = position.AllFour
	}

	w, h := gtx.Dp(r.Width), gtx.Dp(r.Height)

	p := &canvas.Path{}

	started := false

	draw := func(x0, y0, x1, y1 float64) {
		if !started {
			p.MoveTo(x0, y1)
			started = true
		}
		p.LineTo(x1, y1)
	}

	if slices.Contains(positions, position.Top) {
		draw(0, 0, float64(w), 0)
	}

	if slices.Contains(positions, position.Right) {
		draw(float64(w), 0, float64(w), float64(h))
	}

	if slices.Contains(positions, position.Bottom) {
		draw(float64(w), float64(h), 0, float64(h))
	}

	if slices.Contains(positions, position.Left) {
		draw(0, float64(h), 0, 0)
	}

	if len(positions) == 4 {
		p.Close()
	}

	return p
}

func NewRoundedRect() *RoundedRect {
	return &RoundedRect{
		Rect: NewRect(),
	}
}

type RoundedRect struct {
	*Rect
	CornerRadiusValues
}

func (r *RoundedRect) Equals(v any) cmp.Result {
	return cmp.Cast(v, func(v *RoundedRect) cmp.Result {
		return r.Eq(v)
	})
}

func (r *RoundedRect) Eq(v *RoundedRect) cmp.Result {
	return cmp.All(
		r.Rect.Eq(v.Rect),
		r.CornerRadiusValues.Eq(&v.CornerRadiusValues),
	)
}

func (r *RoundedRect) Path(gtx layout.Context, positions ...position.Position) *canvas.Path {
	if len(positions) == 0 {
		positions = position.AllFour
	}

	rect := r.Rectangle(gtx)

	maxRounded := min(gtx.Metric.PxToDp(rect.Max.X), gtx.Metric.PxToDp(rect.Max.Y)) / 2
	se, sw, nw, ne := float64(gtx.Dp(roundedFix(r.BottomRight, maxRounded))), float64(gtx.Dp(roundedFix(r.BottomLeft, maxRounded))), float64(gtx.Dp(roundedFix(r.TopLeft, maxRounded))), float64(gtx.Dp(roundedFix(r.TopRight, maxRounded)))
	w, n, e, s := float64(0), float64(0), float64(rect.Max.X), float64(rect.Max.Y)

	p := &canvas.Path{}

	started := false
	draw := func(x0, y0, r0, x1, y1, r1 float64, pos position.Position) {
		var x0r, y0r float64

		if r0 > 0 {
			r0off := r0 - (r0 * math.Sin(math.Pi/4))

			switch pos {
			case position.Top:
				x0r, y0r = x0+r0off, y0+r0off
				x0 += r0

			case position.Right:
				x0r, y0r = x0-r0off, y0+r0off
				y0 += r0

			case position.Bottom:
				x0r, y0r = x0-r0off, y0-r0off
				x0 -= r0

			case position.Left:
				x0r, y0r = x0+r0off, y0-r0off
				y0 -= r0
			}
		}

		if !started {
			if r0 > 0 {
				p.MoveTo(x0r, y0r)
			} else {
				p.MoveTo(x0, y0)
			}
			started = true
		}

		if r0 > 0 {
			// clockwise 45 deg
			p.ArcTo(r0, r0, math.Pi/4, false, true, x0, y0)
		} else {
			p.LineTo(x0, y0)
		}

		var x1r, y1r float64

		if r1 > 0 {
			r1off := r1 - (r1 * math.Sin(math.Pi/4))

			switch pos {
			case position.Top:
				x1r, y1r = x1-r1off, y1+r1off
				x1 -= r1

			case position.Right:
				x1r, y1r = x1-r1off, y1-r1off
				y1 -= r1

			case position.Bottom:
				x1r, y1r = x1+r1off, y1-r1off
				x1 += r1

			case position.Left:
				x1r, y1r = x1+r1off, y1+r1off
				y1 += r1
			}
		}

		p.LineTo(x1, y1)

		if r1 > 0 && x1r > 0 {
			// clockwise 45 deg
			p.ArcTo(r1, r1, math.Pi/4, false, true, x1r, y1r)
		}
	}

	if slices.Contains(positions, position.Top) {
		draw(w, n, nw, e, n, ne, position.Top)
	}

	if slices.Contains(positions, position.Right) {
		draw(e, n, ne, e, s, se, position.Right)
	}

	if slices.Contains(positions, position.Bottom) {
		draw(e, s, se, w, s, sw, position.Bottom)
	}

	if slices.Contains(positions, position.Left) {
		draw(w, s, sw, w, n, nw, position.Left)
	}

	if len(positions) == 4 {
		p.Close()
	}

	return p
}

func roundedFix(v unit.Dp, maxR unit.Dp) unit.Dp {
	if v >= maxR {
		return maxR
	}
	return v
}
