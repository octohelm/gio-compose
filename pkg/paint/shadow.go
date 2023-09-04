package paint

import (
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/octohelm/gio-compose/pkg/cmp"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"github.com/octohelm/gio-compose/pkg/bezier"
)

type ShadowSetter interface {
	SetShadow(elevation unit.Dp)
}

var _ ShadowSetter = &Shadow{}

type Shadow struct {
	elevation unit.Dp

	ambient  Fill
	penumbra Fill
	umbra    Fill
}

func (s *Shadow) Eq(v *Shadow) cmp.Result {
	return func() bool {
		return s.elevation == v.elevation
	}
}

func (s *Shadow) SetShadow(elevation unit.Dp) {
	s.elevation = elevation
}

func (s *Shadow) Paint(gtx layout.Context, shape SizedShape) {
	if s.elevation <= 0 {
		return
	}

	if s.ambient.Transparent() {
		s.ambient.Color = color.NRGBA{A: 0x10}
	}
	if s.umbra.Transparent() {
		s.umbra.Color = color.NRGBA{A: 0x20}
	}
	if s.penumbra.Transparent() {
		s.penumbra.Color = color.NRGBA{A: 0x30}
	}

	s.paint(gtx, shape, s.ambient.Color, ambientShadow)
	s.paint(gtx, shape, s.umbra.Color, umbraShadow)
	s.paint(gtx, shape, s.penumbra.Color, penumbraShadow)
}

const pi float32 = math.Pi * 2
const shadowDirections float32 = 8 // more is better
const shadowQuality float32 = 3    // more is better
const deltaD = pi / shadowDirections
const deltaR float32 = 1.0 / shadowQuality

func (s *Shadow) paint(gtx layout.Context, shape SizedShape, c color.NRGBA, se *shadowElevation) {
	offsetY, blurRadius, spread := se.Calc(int(s.elevation))
	defer op.Offset(image.Pt(0, gtx.Dp(unit.Dp(offsetY)))).Push(gtx.Ops).Pop()

	rrect := s.rrectWithSpread(gtx, shape, unit.Dp(spread))

	br := gtx.Dp(unit.Dp(blurRadius))

	// inspired by https://www.shadertoy.com/view/Xltfzj
	// the version in comment
	factor := shadowQuality*shadowDirections + (shadowDirections/2 - 1)

	for d := float32(0.0); d < pi; d += deltaD {
		for r := deltaR; r <= 1.0; r += deltaR {
			func() {
				// add some random offset to be smoothly
				randomD := deltaD / float32(rand.Intn(int(shadowDirections)))

				scaledColor := c

				scaledColor.A = uint8(float32(c.A) / factor)
				if scaledColor.A == 0 {
					// FIXME have to set the min alpha, until fixed https://todo.sr.ht/~eliasnaur/gio/532
					scaledColor.A = 0x02
				}

				offset := f32.Pt(
					float32(br)*r*float32(math.Cos(float64(d+randomD))),
					float32(br)*r*float32(math.Sin(float64(d+randomD))),
				)

				defer op.Affine(f32.Affine2D{}.Offset(offset)).Push(gtx.Ops).Pop()
				paint.FillShape(gtx.Ops, scaledColor, rrect.Op(gtx.Ops))
			}()
		}
	}
}

func (s *Shadow) rrectWithSpread(gtx layout.Context, shape Shape, spread unit.Dp) *clip.RRect {
	rrect := &clip.RRect{
		Rect: outset(shape.Rectangle(gtx), gtx.Dp(spread)),
	}

	if g, ok := shape.(CornerRadiusGetter); ok {
		cr := g.CornerRadius()

		rrect.NE = gtx.Dp(cr.TopRight) + gtx.Dp(spread/2)
		rrect.NW = gtx.Dp(cr.TopLeft) + gtx.Dp(spread/2)
		rrect.SW = gtx.Dp(cr.BottomLeft) + gtx.Dp(spread/2)
		rrect.SE = gtx.Dp(cr.BottomRight) + gtx.Dp(spread/2)

		if rrect.NE < 0 {
			rrect.NE = 0
		}
		if rrect.NW < 0 {
			rrect.NW = 0
		}
		if rrect.SW < 0 {
			rrect.SW = 0
		}
		if rrect.SE < 0 {
			rrect.SE = 0
		}
	}

	return rrect
}

func outset(r image.Rectangle, rr int) image.Rectangle {
	r.Min.X -= rr
	r.Min.Y -= rr
	r.Max.X += rr
	r.Max.Y += rr
	return r
}

var (
	ambientShadow = &shadowElevation{
		maxElevation:      24,
		bezierCurveY:      bezier.Easing(0.4, 0, 1, 0.8),
		bezierCurveBlur:   bezier.Easing(0.4, 0, 1, 0.8),
		bezierCurveSpread: bezier.Easing(0.4, 0, 1, 0.8),
		yBoundaries:       [2]int{0, 4},
		blurBoundaries:    [2]int{0, 64},
		spreadBoundaries:  [2]int{1, 4},
	}

	umbraShadow = &shadowElevation{
		maxElevation:      24,
		bezierCurveY:      bezier.Easing(0.4, 0, 1, 0.8),
		bezierCurveBlur:   bezier.Easing(0.4, 0, 1, 0.8),
		bezierCurveSpread: bezier.Easing(0.4, 0, 1, 0.8),
		yBoundaries:       [2]int{1, 16},
		blurBoundaries:    [2]int{2, 16},
		spreadBoundaries:  [2]int{1, 2},
		negativeSpread:    true,
	}

	penumbraShadow = &shadowElevation{
		maxElevation:      24,
		bezierCurveY:      bezier.Easing(0.4, 0, 1, 0.8),
		bezierCurveBlur:   bezier.Easing(0.4, 0, 1, 0.8),
		bezierCurveSpread: bezier.Easing(0.4, 0, 1, 0.8),
		yBoundaries:       [2]int{0, 32},
		blurBoundaries:    [2]int{0, 64},
		spreadBoundaries:  [2]int{1, 5},
	}
)

type shadowElevation struct {
	maxElevation      int
	yBoundaries       boundaries
	blurBoundaries    boundaries
	spreadBoundaries  boundaries
	bezierCurveY      bezier.EasingFunc
	bezierCurveBlur   bezier.EasingFunc
	bezierCurveSpread bezier.EasingFunc
	negativeSpread    bool
}

type boundaries [2]int

func (b boundaries) At(p float64) int {
	return int(math.Round(float64(b[1]-b[0])*p)) + b[0]
}

func (se *shadowElevation) Calc(elevation int) (y int, blur int, spread int) {
	p := float64(elevation) * 1 / (float64(se.maxElevation) - 1)

	y = se.yBoundaries.At(se.bezierCurveY(p))
	blur = se.blurBoundaries.At(se.bezierCurveBlur(p))
	spread = se.spreadBoundaries.At(se.bezierCurveSpread(p))

	if se.negativeSpread {
		spread = -spread
	}

	return
}
