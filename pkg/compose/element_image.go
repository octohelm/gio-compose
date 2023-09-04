package compose

import (
	"context"
	"image"

	"github.com/octohelm/gio-compose/pkg/modifier"

	"gioui.org/op/clip"
	"github.com/octohelm/gio-compose/pkg/paint/canvas"

	giopaint "gioui.org/op/paint"

	"github.com/octohelm/gio-compose/pkg/cmp"
	"github.com/octohelm/gio-compose/pkg/compose/internal"
	"github.com/octohelm/gio-compose/pkg/layout"
	"github.com/octohelm/gio-compose/pkg/node"
	"github.com/octohelm/gio-compose/pkg/paint"
)

func Image(image image.Image, modifiers ...modifier.Modifier[any]) VNode {
	return H(&imageWidget{
		Element: node.Element{
			Name: "Image",
		},
	}, modifier.ModifiersOf(modifiers...), &imageModifier{image: image})
}

type imageModifier struct {
	image image.Image
}

func (m *imageModifier) Modify(w any) {
	if setter, ok := w.(ImageSetter); ok {
		setter.SetImage(m.image)
	}
}

type ImageSetter interface {
	SetImage(img image.Image)
}

var _ Element = &imageWidget{}

type imageWidget struct {
	internal.ElementComponent
	node.Element

	imageWidgetAttrs
}

func (w *imageWidget) New(ctx context.Context) internal.Element {
	return &imageWidget{
		Element: node.Element{
			Name: w.Name,
		},
	}
}

var _ ImageSetter = &imageWidgetAttrs{}

type imageWidgetAttrs struct {
	image image.Image
	paint.ContentScale
	*container
}

func (attrs *imageWidgetAttrs) SetImage(img image.Image) {
	attrs.image = img
}

func (attrs *imageWidgetAttrs) Eq(next *imageWidgetAttrs) cmp.Result {
	return cmp.All(
		attrs.container.Eq(next.container),
		func() bool {
			return attrs.image == next.image
		},
	)
}

func (w *imageWidget) Update(ctx context.Context, modifiers ...modifier.Modifier[any]) bool {
	attrs := &imageWidgetAttrs{container: newContainer()}

	modifier.Modify[any](attrs, modifiers...)

	return cmp.UpdateWhen(
		cmp.Not(w.imageWidgetAttrs.Eq(attrs)),
		&w.imageWidgetAttrs, attrs,
	)
}

func (w *imageWidget) Layout(gtx layout.Context) layout.Dimensions {
	if w.image == nil {
		return layout.Dimensions{}
	}

	return w.container.Layout(gtx, w, ElementPainterFunc(func(gtx layout.Context) layout.Dimensions {
		// clipped
		defer clip.Outline{Path: canvas.PathSpec(gtx.Ops, w.container.shape.Path(gtx))}.Op().Push(gtx.Ops).Pop()

		imageOp := giopaint.NewImageOp(w.image)
		imageOp.Add(gtx.Ops)

		w.ContentScale.Fit(gtx.Ops, w.image.Bounds().Size(), gtx.Constraints.Max)
		giopaint.PaintOp{}.Add(gtx.Ops)
		return layout.Dimensions{Size: gtx.Constraints.Max}
	}))
}
