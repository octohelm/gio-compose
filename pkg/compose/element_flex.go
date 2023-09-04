package compose

import (
	"context"
	"image"
	"math"

	"github.com/octohelm/gio-compose/pkg/modifier"

	"github.com/octohelm/gio-compose/pkg/unit"

	"gioui.org/op"

	giolayout "gioui.org/layout"

	"github.com/octohelm/gio-compose/pkg/cmp"
	"github.com/octohelm/gio-compose/pkg/compose/internal"
	"github.com/octohelm/gio-compose/pkg/layout"
	"github.com/octohelm/gio-compose/pkg/layout/alignment"
	"github.com/octohelm/gio-compose/pkg/layout/arrangement"
	"github.com/octohelm/gio-compose/pkg/layout/direction"
	"github.com/octohelm/gio-compose/pkg/node"
	"github.com/octohelm/gio-compose/pkg/paint"
	"github.com/octohelm/gio-compose/pkg/paint/size"
)

func Column(modifiers ...modifier.Modifier[any]) VNode {
	return H(&flexElement{
		Element: node.Element{
			Name: "Column",
		},
		Axis: direction.Vertical,
	}, modifiers...)
}

func Row(modifiers ...modifier.Modifier[any]) VNode {
	return H(&flexElement{
		Element: node.Element{
			Name: "Row",
		},
		Axis: direction.Horizontal,
	}, modifiers...)
}

var _ Element = &flexElement{}

type flexElement struct {
	internal.ElementComponent
	node.Element
	Axis direction.Axis

	flexWidgetAttrs
	list *giolayout.List
}

func (fe *flexElement) Update(ctx context.Context, modifiers ...modifier.Modifier[any]) bool {
	attrs := &flexWidgetAttrs{
		container: newContainer(),
	}

	modifier.Modify[any](attrs, modifiers...)

	return cmp.UpdateWhen(
		cmp.Not(fe.flexWidgetAttrs.Eq(attrs)),
		&fe.flexWidgetAttrs, attrs,
	)
}

type flexWidgetAttrs struct {
	*container

	layout.Spacer
	layout.Aligner
	layout.Arrangementer
	layout.Scrollable
}

func (attrs *flexWidgetAttrs) Eq(v *flexWidgetAttrs) cmp.Result {
	return cmp.All(
		attrs.container.Eq(v.container),
		attrs.Scrollable.Eq(&v.Scrollable),
		attrs.Spacer.Eq(&v.Spacer),
		attrs.Aligner.Eq(&v.Aligner),
		attrs.Arrangementer.Eq(&v.Arrangementer),
	)
}

func (fe *flexElement) New(ctx context.Context) Element {
	return &flexElement{
		Element: node.Element{
			Name: fe.Name,
		},
		Axis: fe.Axis,
	}
}

func (fe *flexElement) Layout(gtx layout.Context) layout.Dimensions {
	if fe.Scrollable.Enabled {
		return fe.container.Layout(gtx, fe, ElementPainterFunc(func(gtx layout.Context) layout.Dimensions {
			children := make([]Element, 0)

			for child := range node.IterChildElement(context.Background(), fe) {
				if w, ok := child.(Element); ok {
					children = append(children, w)
				}
			}

			if fe.list == nil {
				// bind instance to store list state
				fe.list = &giolayout.List{
					Axis:      fe.Scrollable.Axis.LayoutAxis(),
					Alignment: fe.Alignment.LayoutAlignment(),
				}
			}

			maxViewSize := 0

			switch fe.Scrollable.Axis {
			case direction.Vertical:
				maxViewSize = gtx.Constraints.Max.Y
			case direction.Horizontal:
				maxViewSize = gtx.Constraints.Max.X
			}

			childrenOffsets := map[int]int{}

			isVisible := func(index int) bool {
				return index >= fe.list.Position.First && index <= fe.list.Position.First+fe.list.Position.Count
			}

			defer func() {
				for index := range children {
					positionChild(fe, children[index], func() (x, y unit.Dp) {
						childOffset := -maxViewSize

						if isVisible(index) {
							if offset, ok := childrenOffsets[index]; ok {
								childOffset = offset
							}
						} else {
							childOffset = -maxViewSize
						}

						x, y = unit.Dp(0), unit.Dp(0)

						switch fe.Scrollable.Axis {
						case direction.Vertical:
							x, y = 0, gtx.Metric.PxToDp(childOffset)
						case direction.Horizontal:
							x, y = gtx.Metric.PxToDp(childOffset), 0
						}

						return
					})
				}
			}()

			visibleOffset := -fe.list.Position.Offset

			return fe.list.Layout(gtx, len(children), func(gtx layout.Context, index int) layout.Dimensions {
				c := children[index]

				dims := c.Layout(gtx)

				if isVisible(index) {
					// when visible
					childrenOffsets[index] = visibleOffset
					visibleOffset += dims.Size.Y
				}

				return dims
			})
		}))
	}

	return fe.container.Layout(gtx, fe, ElementPainterFunc(func(gtx layout.Context) layout.Dimensions {
		children := make([]*flexChild, 0)

		idx := 0
		addFlexChild := func(child *flexChild) {
			if idx > 0 && fe.Spacing != 0 {
				children = append(children, sized(Graph(func(gtx layout.Context) layout.Dimensions {
					if fe.Axis == direction.Vertical {
						return layout.Dimensions{
							Size: image.Point{
								Y: gtx.Dp(fe.Spacing),
							},
						}
					}
					return layout.Dimensions{
						Size: image.Point{
							X: gtx.Dp(fe.Spacing),
						},
					}
				})))
			}

			children = append(children, child)
			idx++
		}

		calculator := &flexCalculator{axis: fe.Axis}
		calculator.reset(gtx.Constraints.Min)

		for child := range node.IterChildElement(context.Background(), fe) {
			if w, ok := child.(Element); ok {
				weight := float32(0)

				if getter, ok := w.(layout.WeightGetter); ok {
					if v, ok := getter.Weight(); ok {
						weight = v
					}
				}

				// force weight when arrangement is EqualWeight
				if fe.Arrangement == arrangement.EqualWeight {
					weight = 1
				}

				if weight != 0 {
					addFlexChild(flexed(weight, w))
					calculator.incrWeight(weight)
				} else {
					addFlexChild(sized(w))
				}
			}
		}

		if len(children) == 0 {
			return layout.Dimensions{}
		}

		for i := range children {
			child := children[i]

			if !child.flexed() {
				child.Paint(gtx)
				calculator.incr(child.dims.Size)
			}
		}

		remainSpacing := calculator.remainSpacing(gtx.Constraints.Max)

		for i := range children {
			child := children[i]

			if child.flexed() {
				if remainSpacing > 0 {
					gtx.Constraints = calculator.constraints(gtx.Constraints, child.weight, remainSpacing)
				}

				if sc, ok := child.widget.(paint.SizeSetter); ok {
					switch fe.Axis {
					case direction.Horizontal:
						sc.SetSize(-1, size.Width)
					case direction.Vertical:
						sc.SetSize(-1, size.Height)
					}
				}

				child.Paint(gtx)
				calculator.incr(child.dims.Size)
			}
		}

		offset := 0

		for i := range children {
			child := children[i]

			positionChild(fe, child.widget, func() (x, y unit.Dp) {
				off := image.Point{}
				off = off.Add(image.Pt(fe.offsetOfAlignment(calculator.size, child.dims.Size)))
				off = off.Add(image.Pt(fe.offsetOfArrangement(remainSpacing, len(children), i)))

				switch fe.Axis {
				case direction.Horizontal:
					off.X += offset
					offset += child.dims.Size.X
				case direction.Vertical:
					off.Y += offset
					offset += child.dims.Size.Y
				}

				defer op.Offset(image.Pt(off.X, off.Y)).Push(gtx.Ops).Pop()
				child.call.Add(gtx.Ops)

				return gtx.Metric.PxToDp(off.X), gtx.Metric.PxToDp(off.Y)
			})
		}

		return layout.Dimensions{
			Size: calculator.size,
		}
	}))
}

type flexCalculator struct {
	axis        direction.Axis
	size        image.Point
	totalWeight float32
}

func (c *flexCalculator) reset(size image.Point) {
	switch c.axis {
	case direction.Horizontal:
		c.size.Y = size.Y
	case direction.Vertical:
		c.size.X = size.X
	}
}

func (c *flexCalculator) incr(size image.Point) {
	switch c.axis {
	case direction.Horizontal:
		c.size.X += size.X

		if size.Y > c.size.Y {
			c.size.Y = size.Y
		}
	case direction.Vertical:
		c.size.Y += size.Y

		if size.X > c.size.X {
			c.size.X = size.X
		}
	}
}

func (c *flexCalculator) remainSpacing(max image.Point) int {
	switch c.axis {
	case direction.Horizontal:
		return max.X - c.size.X
	case direction.Vertical:
		return max.Y - c.size.Y
	}
	return 0
}

func (c *flexCalculator) constraints(constraints layout.Constraints, weight float32, spacing int) layout.Constraints {
	if c.totalWeight != 0 {
		switch c.axis {
		case direction.Horizontal:
			constraints.Min.X = int(math.Round(float64(float32(spacing) * weight / c.totalWeight)))
			constraints.Max.X = constraints.Min.X
		case direction.Vertical:
			constraints.Min.Y = int(math.Round(float64(float32(spacing) * weight / c.totalWeight)))
			constraints.Max.Y = constraints.Min.Y
		}
	}

	return constraints
}

func (c *flexCalculator) incrWeight(weight float32) {
	c.totalWeight += weight
}

func (fe *flexElement) offsetOfAlignment(parent image.Point, child image.Point) (int, int) {
	switch fe.Alignment {
	case alignment.End:
		if fe.Axis == direction.Horizontal {
			return 0, parent.Y - child.Y
		}
		return parent.X - child.X, 0
	case alignment.Center:
		if fe.Axis == direction.Horizontal {
			return 0, (parent.Y - child.Y) / 2
		}
		return (parent.X - child.X) / 2, 0
	}
	return 0, 0
}

func (fe *flexElement) offsetOfArrangement(spacing int, n int, idx int) (int, int) {
	if spacing <= 0 {
		return 0, 0
	}

	switch fe.Arrangement {
	case arrangement.SpaceEvenly:
		if fe.Axis == direction.Horizontal {
			return spacing / (n + 1) * (idx + 1), 0
		}
		return 0, spacing / (n + 1) * (idx + 1)
	case arrangement.SpaceAround:
		if fe.Axis == direction.Horizontal {
			return spacing / (n * 2) * (idx*2 + 1), 0
		}
		return 0, spacing / (n * 2) * (idx*2 + 1)
	case arrangement.SpaceBetween:
		if (n - 1) == 0 {
			return 0, 0
		}
		if fe.Axis == direction.Horizontal {
			return spacing / (n - 1) * idx, 0
		}
		return 0, spacing / (n - 1) * idx
	case arrangement.Start:
		if fe.Axis == direction.Horizontal {
			return 0, 0
		}
		return 0, 0
	case arrangement.End:
		if fe.Axis == direction.Horizontal {
			return spacing, 0
		}
		return 0, spacing
	case arrangement.Center:
		if fe.Axis == direction.Horizontal {
			return spacing / 2, 0
		}
		return 0, spacing / 2
	}
	return 0, 0
}

func sized(w Element) *flexChild {
	return &flexChild{widget: w}
}

func flexed(weight float32, w Element) *flexChild {
	return &flexChild{weight: weight, widget: w}
}

type flexChild struct {
	weight float32
	widget Element

	call    op.CallOp
	dims    layout.Dimensions
	painted bool
}

func (c *flexChild) flexed() bool {
	return c.weight > 0
}

func (c *flexChild) Paint(gtx layout.Context) {
	if c.painted {
		return
	}

	c.painted = true
	c.call = paint.Group(gtx.Ops, func() {
		c.dims = c.widget.Layout(gtx)
	})
}
