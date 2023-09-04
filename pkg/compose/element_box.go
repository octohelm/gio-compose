package compose

import (
	"context"
	"image"

	"gioui.org/op"

	"github.com/octohelm/gio-compose/pkg/cmp"
	"github.com/octohelm/gio-compose/pkg/compose/internal"
	"github.com/octohelm/gio-compose/pkg/layout"
	"github.com/octohelm/gio-compose/pkg/layout/alignment"
	"github.com/octohelm/gio-compose/pkg/modifier"
	"github.com/octohelm/gio-compose/pkg/node"
	"github.com/octohelm/gio-compose/pkg/paint"
	"github.com/octohelm/gio-compose/pkg/paint/size"
	"github.com/octohelm/gio-compose/pkg/unit"
)

func Box(modifiers ...modifier.Modifier[any]) VNode {
	return H(
		&boxElement{
			Element: node.Element{
				Name: "Box",
			},
		},
		modifiers...,
	)
}

var _ Element = &boxElement{}

type boxElement struct {
	internal.ElementComponent
	node.Element
	boxWidgetAttrs
}

func (b boxElement) New(ctx context.Context) internal.Element {
	return &boxElement{
		Element: node.Element{
			Name: b.Name,
		},
	}
}

func (b *boxElement) Update(ctx context.Context, modifiers ...modifier.Modifier[any]) bool {
	attrs := &boxWidgetAttrs{
		container: newContainer(),
	}

	modifier.Modify[any](attrs, modifiers...)

	return cmp.UpdateWhen(
		cmp.Not(b.boxWidgetAttrs.Eq(attrs)),
		&b.boxWidgetAttrs, attrs,
	)
}

type boxWidgetAttrs struct {
	*container
	layout.Aligner
}

func (attrs *boxWidgetAttrs) Eq(v *boxWidgetAttrs) cmp.Result {
	return cmp.All(
		attrs.container.Eq(v.container),
		attrs.Aligner.Eq(&v.Aligner),
	)
}

func (b *boxElement) Layout(gtx layout.Context) (d layout.Dimensions) {
	return b.container.Layout(gtx, b, ElementPainterFunc(func(gtx layout.Context) layout.Dimensions {
		sizeMax := gtx.Constraints.Min
		children := make([]*boxChild, 0)

		for child := range node.IterChildElement(context.Background(), b) {
			if w, ok := child.(Element); ok {
				bc := &boxChild{widget: w}

				noAutoWidth := bc.Sized(size.Width) != size.MatchParent
				noAutoHeight := bc.Sized(size.Height) != size.MatchParent

				if noAutoWidth || noAutoHeight {
					bc.Paint(gtx)

					if noAutoWidth {
						if w := bc.dims.Size.X; w > sizeMax.X {
							sizeMax.X = w
						}
					}

					if noAutoHeight {
						if h := bc.dims.Size.Y; h > sizeMax.Y {
							sizeMax.Y = h
						}
					}
				}

				children = append(children, bc)
			}
		}

		for i := range children {
			children[i].DrawTo(gtx, b, sizeMax, len(children))
		}

		return layout.Dimensions{Size: sizeMax}
	}))
}

type boxChild struct {
	widget Element

	painted bool
	group   op.CallOp
	dims    layout.Dimensions
}

func (bc *boxChild) Sized(t size.Type) size.SizingType {
	c, ok := bc.widget.(paint.SizedChecker)
	if ok {
		return c.Sized(t)
	}
	return size.WrapContent
}

func (bc *boxChild) Paint(gtx layout.Context) {
	if bc.painted {
		return
	}

	bc.painted = true
	bc.group = paint.Group(gtx.Ops, func() {
		bc.dims = bc.widget.Layout(gtx)
	})
}

func (bc *boxChild) DrawTo(gtx layout.Context, parent Element, maxSize image.Point, n int) {
	positionChild(parent, bc.widget, func() (x, y unit.Dp) {
		if bc.Sized(size.Height) == size.MatchParent || bc.Sized(size.Width) == size.MatchParent {
			if n > 1 {
				gtx.Constraints.Max = maxSize
			}
			bc.Paint(gtx)
			if n == 1 {
				maxSize = bc.dims.Size
			}
		}

		align := alignment.Center
		if alignGetter, ok := bc.widget.(layout.AlignGetter); ok {
			align = alignGetter.Align()
		}

		x, y = calcPosition(
			align,
			gtx.Metric.PxToDp(maxSize.X), gtx.Metric.PxToDp(maxSize.Y),
			gtx.Metric.PxToDp(bc.dims.Size.X), gtx.Metric.PxToDp(bc.dims.Size.Y),
		)

		defer op.Offset(image.Pt(gtx.Dp(x), gtx.Dp(y))).Push(gtx.Ops).Pop()
		bc.group.Add(gtx.Ops)
		return
	})
}

func calcPosition(align alignment.Alignment, w, h unit.Dp, cw, ch unit.Dp) (unit.Dp, unit.Dp) {
	switch align {
	case alignment.TopStart:
		return 0, 0
	case alignment.Top:
		return (w - cw) / 2, 0
	case alignment.TopEnd:
		return w - cw, 0
	case alignment.Start:
		return 0, (h - ch) / 2
	case alignment.Center:
		return (w - cw) / 2, (h - ch) / 2
	case alignment.End:
		return w - cw, (h - ch) / 2
	case alignment.BottomStart:
		return 0, h - cw
	case alignment.Bottom:
		return (w - cw) / 2, h - ch
	case alignment.BottomEnd:
		return w - cw, h - ch
	}
	return 0, 0
}
