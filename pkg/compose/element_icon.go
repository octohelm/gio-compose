package compose

import (
	"bytes"
	"context"
	"image"
	"image/draw"

	"github.com/octohelm/gio-compose/pkg/modifier"

	"gioui.org/op/clip"
	"github.com/octohelm/gio-compose/pkg/paint/canvas"
	"github.com/octohelm/gio-compose/pkg/paint/size"
	"github.com/octohelm/gio-compose/pkg/unit"

	giopaint "gioui.org/op/paint"
	"github.com/octohelm/gio-compose/pkg/paint/f32color"
	"github.com/octohelm/x/ptr"
	"golang.org/x/exp/shiny/iconvg"

	"github.com/octohelm/gio-compose/pkg/cmp"
	"github.com/octohelm/gio-compose/pkg/compose/internal"
	"github.com/octohelm/gio-compose/pkg/layout"
	"github.com/octohelm/gio-compose/pkg/node"
	"github.com/octohelm/gio-compose/pkg/paint"
	"github.com/octohelm/gio-compose/pkg/text"
)

// Icon
// only support iconvg golang.org/x/exp/shiny/iconvg
func Icon(src []byte, modifiers ...modifier.Modifier[any]) VNode {
	return H(&iconWidget{
		Element: node.Element{
			Name: "Icon",
		},
	}, modifier.ModifiersOf(modifiers...), &iconModifier{src: src})
}

type iconModifier struct {
	src []byte
}

func (m *iconModifier) Modify(w any) {
	if setter, ok := w.(IconvgSetter); ok {
		setter.SetIconvg(m.src)
	}
}

type IconvgSetter interface {
	SetIconvg(src []byte)
}

var _ Element = &iconWidget{}

type iconWidget struct {
	internal.ElementComponent
	node.Element

	iconWidgetAttrs
}

func (w *iconWidget) New(ctx context.Context) internal.Element {
	return &iconWidget{
		Element: node.Element{
			Name: w.Name,
		},
	}
}

var _ IconvgSetter = &iconWidgetAttrs{}

type iconWidgetAttrs struct {
	*container
	imagevg []byte
	paint.ContentScale
	textStyle text.Style
	imageOp   *giopaint.ImageOp
	iconSize  image.Point
}

func (attrs *iconWidgetAttrs) SetIconvg(src []byte) {
	attrs.imagevg = src
}

func (attrs *iconWidgetAttrs) Eq(next *iconWidgetAttrs) cmp.Result {
	return cmp.All(
		attrs.container.Eq(next.container),
		attrs.textStyle.Eq(&next.textStyle),
		func() bool {
			return bytes.Equal(attrs.imagevg, next.imagevg)
		},
	)
}

func (w *iconWidget) Update(ctx context.Context, modifiers ...modifier.Modifier[any]) bool {
	attrs := &iconWidgetAttrs{
		container: newContainer(),
		textStyle: TextStyleContext.Extract(ctx),
	}

	// use font size as the default icon size
	attrs.SetSize(unit.Dp(attrs.textStyle.FontSize), size.Width, size.Height)

	modifier.Modify[any](attrs, modifiers...)

	return cmp.UpdateWhen(
		cmp.Not(w.iconWidgetAttrs.Eq(attrs)),
		&w.iconWidgetAttrs, attrs,
	)
}

func (w *iconWidget) Layout(gtx layout.Context) layout.Dimensions {
	if w.imagevg == nil {
		return layout.Dimensions{}
	}

	return w.container.Layout(gtx, w, ElementPainterFunc(func(gtx layout.Context) layout.Dimensions {
		// clipped
		defer clip.Outline{Path: canvas.PathSpec(gtx.Ops, w.container.shape.Path(gtx))}.Op().Push(gtx.Ops).Pop()

		if w.imageOp == nil {
			sz := max(gtx.Constraints.Max.X, gtx.Constraints.Max.Y)

			m, _ := iconvg.DecodeMetadata(w.imagevg)
			dx, dy := m.ViewBox.AspectRatio()
			img := image.NewRGBA(image.Rectangle{Max: image.Point{X: sz, Y: int(float32(sz) * dy / dx)}})
			ico := &iconvg.Rasterizer{}
			ico.SetDstImage(img, img.Bounds(), draw.Src)

			m.Palette[0] = f32color.NRGBAToLinearRGBA(w.textStyle.Color)
			_ = iconvg.Decode(ico, w.imagevg, &iconvg.DecodeOptions{
				Palette: &m.Palette,
			})

			w.imageOp = ptr.Ptr(giopaint.NewImageOp(img))
			w.iconSize = img.Bounds().Max
		}

		w.imageOp.Add(gtx.Ops)
		w.ContentScale.Fit(gtx.Ops, w.iconSize, gtx.Constraints.Max)
		giopaint.PaintOp{}.Add(gtx.Ops)
		return layout.Dimensions{Size: gtx.Constraints.Max}
	}))
}
