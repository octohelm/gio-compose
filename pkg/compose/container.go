package compose

import (
	"image"

	"github.com/octohelm/gio-compose/pkg/event/textinput"

	"gioui.org/op/clip"

	"github.com/octohelm/gio-compose/pkg/paint/canvas"

	"gioui.org/op"

	"github.com/octohelm/gio-compose/pkg/cmp"
	"github.com/octohelm/gio-compose/pkg/event/gesture"
	"github.com/octohelm/gio-compose/pkg/layout"
	"github.com/octohelm/gio-compose/pkg/layout/position"
	"github.com/octohelm/gio-compose/pkg/node"
	"github.com/octohelm/gio-compose/pkg/paint"
	"github.com/octohelm/gio-compose/pkg/unit"
)

func newContainer() *container {
	return &container{
		shape: &shape{SizedShape: paint.NewRoundedRect()},
	}
}

type ParentInsetSetter interface {
	SetParentInset(x unit.Dp, y unit.Dp)
}

type container struct {
	*shape

	layout.Clip
	layout.Offset
	layout.FlexWeight
	layout.EdgeInset
	layout.Scrollable

	gesture.PointerEventDetector
	textinput.InputEventDetector

	layout.PhaseRecorder
}

func (c *container) Eq(v *container) cmp.Result {
	if c == nil {
		return func() bool {
			return false
		}
	}

	return cmp.All(
		c.shape.Eq(v.shape),
		c.FlexWeight.Eq(&v.FlexWeight),
		c.Offset.Eq(&v.Offset),
		c.EdgeInset.Eq(&v.EdgeInset),
	)
}

type WidgetWithPositionBy interface {
	Element
	PositionBy(calc func() (x unit.Dp, y unit.Dp))
}

func (c *container) PositionBy(calc func() (x unit.Dp, y unit.Dp)) {
	c.PhaseRecorder.PositionBy(func() (x unit.Dp, y unit.Dp) {
		x, y = calc()
		return c.Offset.X + x, c.Offset.Y + y
	})
}

func positionChild(parent Element, child Element, position func() (x, y unit.Dp)) {
	if setter, ok := child.(WidgetWithPositionBy); ok {
		setter.PositionBy(func() (unit.Dp, unit.Dp) {
			x, y := position()
			insetOffset := layout.Offset{}
			if getter, ok := parent.(layout.EdgeInsetOffsetGetter); ok {
				insetOffset = getter.EdgeInsetOffset()
			}
			return insetOffset.X + x, insetOffset.Y + y
		})
		return
	} else {
		// calc without recording
		_, _ = position()
	}
}

func (c *container) Layout(gtx layout.Context, n node.Node, w ElementPainter) (dims layout.Dimensions) {
	child := paint.Group(gtx.Ops, func() {
		defer c.PhaseRecorder.Trigger(layout.PhaseBeforeSize, nil)
		defer func() {
			c.PhaseRecorder.RecordSize(gtx.Metric.PxToDp(dims.Size.X), gtx.Metric.PxToDp(dims.Size.Y))
			c.PhaseRecorder.Trigger(layout.PhaseDidSize, nil)
		}()

		dims = c.PointerEventDetector.LayoutChild(
			gtx,
			n,
			func(gtx layout.Context) (childDims layout.Dimensions) {
				// set min to zero
				gtx.Constraints.Min = image.Point{}
				dd := c.shape.Layout(gtx)

				if dd.Size.X > 0 {
					gtx.Constraints.Min.X = dd.Size.X
					gtx.Constraints.Max.X = dd.Size.X

					defer func() {
						childDims.Size.X = dd.Size.X
					}()
				}

				if dd.Size.Y > 0 {
					gtx.Constraints.Min.Y = dd.Size.Y
					gtx.Constraints.Max.Y = dd.Size.Y

					defer func() {
						childDims.Size.Y = dd.Size.Y
					}()
				}

				insetX := gtx.Dp(c.EdgeInset.Left + c.EdgeInset.Right)
				insetY := gtx.Dp(c.EdgeInset.Top + c.EdgeInset.Bottom)

				childCore := paint.Group(gtx.Ops, func() {
					gtx.Constraints.Max = gtx.Constraints.Max.Sub(image.Pt(insetX, insetY))

					if gtx.Constraints.Min.X > gtx.Constraints.Max.X {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
					}
					if gtx.Constraints.Min.Y > gtx.Constraints.Max.Y {
						gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
					}

					childDims = w.Layout(gtx)
				})

				defer op.Offset(image.Pt(gtx.Dp(c.EdgeInset.Left), gtx.Dp(c.EdgeInset.Top))).Push(gtx.Ops).Pop()
				childCore.Add(gtx.Ops)

				return layout.Dimensions{
					Size:     childDims.Size.Add(image.Pt(insetX, insetY)),
					Baseline: childDims.Baseline + gtx.Dp(c.EdgeInset.Bottom),
				}
			})

		c.InputEventDetector.Layout(gtx)
	})

	if !c.Offset.IsZero() {
		defer op.Offset(image.Pt(gtx.Dp(c.Offset.X), gtx.Dp(c.Offset.Y))).Push(gtx.Ops).Pop()
	}

	// fit child size if required
	gtx.Constraints.Min, gtx.Constraints.Max = dims.Size, dims.Size
	// paint shape before gestures layout to avoid to clipped
	c.shape.Paint(gtx)

	if c.Clip.Enabled {
		defer clip.Outline{Path: canvas.PathSpec(gtx.Ops, c.shape.Path(gtx))}.Op().Push(gtx.Ops).Pop()
	}
	// paint child
	child.Add(gtx.Ops)

	return dims
}

type shape struct {
	paint.SizedShape

	paint.Shadow
	paint.Fill
	paint.BorderStroke
}

var _ paint.CornerRadiusSetter = &shape{}

func (s *shape) SetCornerRadius(v unit.Dp, positions ...position.Position) {
	if setter, ok := s.SizedShape.(paint.CornerRadiusSetter); ok {
		setter.SetCornerRadius(v, positions...)
	}
}

func (s *shape) Eq(t *shape) cmp.Result {
	return cmp.All(
		s.SizedShape.Equals(t.SizedShape),
		s.Shadow.Eq(&t.Shadow),
		s.Fill.Eq(&t.Fill),
		s.BorderStroke.Eq(&t.BorderStroke),
	)
}

func (s *shape) Paint(gtx layout.Context) {
	// paint shadow if exists
	s.Shadow.Paint(gtx, s.SizedShape)
	// paint background if exists
	s.Fill.PaintClip(gtx, paint.ClipOutlineOf(canvas.PathSpec(gtx.Ops, s.Path(gtx))))
	// paint then border if exists
	s.BorderStroke.PaintShape(gtx, s)
}

func (s *shape) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Dimensions{
		Size: s.Rectangle(gtx).Max,
	}
}
