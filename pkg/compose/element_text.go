package compose

import (
	"context"
	"image"
	"image/color"

	"gioui.org/op"
	gtext "gioui.org/text"
	"gioui.org/widget"
	"github.com/octohelm/x/ptr"

	"github.com/octohelm/gio-compose/pkg/modifier"
	"github.com/octohelm/gio-compose/pkg/paint"
	"github.com/octohelm/gio-compose/pkg/paint/size"
	"github.com/octohelm/gio-compose/pkg/unit"

	"github.com/octohelm/gio-compose/pkg/cmp"
	"github.com/octohelm/gio-compose/pkg/compose/internal"
	"github.com/octohelm/gio-compose/pkg/layout"
	"github.com/octohelm/gio-compose/pkg/node"
	"github.com/octohelm/gio-compose/pkg/text"
	"github.com/octohelm/gio-compose/pkg/util/contextutil"
)

func Text(text string, modifies ...modifier.Modifier[*text.Style]) VNode {
	return H(&textWidget{
		Element: node.Element{
			Name: "Text",
		},
	}, modifier.ModifiersOf(modifies...), &textModifier{textContents: text})
}

type textModifier struct {
	textContents string
}

func (m *textModifier) Modify(w any) {
	if setter, ok := w.(TextSetter); ok {
		setter.SetText(m.textContents)
	}
}

type TextSetter interface {
	SetText(txt string)
}

var TextStyleContext = contextutil.New[text.Style](contextutil.Defaulter(func() text.Style {
	return text.Style{
		FontFamily: "",
		FontSize:   unit.Sp(16),
		LineHeight: unit.Sp(18),
		Color:      color.NRGBA{A: 0xff},
		TextAlign:  ptr.Ptr(text.Start),
	}
}))

func ProvideTextStyle(modifiers ...modifier.Modifier[*text.Style]) VNode {
	return Provider(func(ctx context.Context) context.Context {
		s := ptr.Ptr(TextStyleContext.Extract(ctx))
		modifier.Modify(s, modifiers...)
		return TextStyleContext.Inject(ctx, *s)
	})
}

type textWidget struct {
	internal.ElementComponent
	node.Element

	textWidgetAttrs
	textWidgetState

	layout.PhaseRecorder
}

type textWidgetState struct {
	shaper *gtext.Shaper
}

func (w textWidget) New(ctx context.Context) internal.Element {
	return &textWidget{
		Element: node.Element{
			Name: w.Name,
		},
	}
}

func (w *textWidget) Update(ctx context.Context, modifiers ...modifier.Modifier[any]) bool {
	w.shaper = text.ShaperContext.Extract(ctx)

	attrs := &textWidgetAttrs{
		Style: TextStyleContext.Extract(ctx),
	}

	modifier.Modify[any](attrs, modifiers...)

	return cmp.UpdateWhen(
		cmp.Not(w.textWidgetAttrs.Eq(attrs)),
		&w.textWidgetAttrs, attrs,
	)
}

type textWidgetAttrs struct {
	text.Style
	textContents string
}

func (attrs *textWidgetAttrs) SetText(c string) {
	attrs.textContents = c
}

var _ Element = &textWidget{}

func (attrs *textWidgetAttrs) Eq(next *textWidgetAttrs) cmp.Result {
	return cmp.All(
		func() bool {
			return attrs.textContents == next.textContents
		},
		attrs.Style.Eq(&next.Style),
	)
}

var _ paint.SizedChecker = &textWidget{}

func (w *textWidget) Sized(typ size.Type) size.SizingType {
	return size.Exactly
}

var _ WidgetWithPositionBy = &textWidget{}

func (w *textWidget) Layout(gtx layout.Context) (dims layout.Dimensions) {
	w.PhaseRecorder.Trigger(layout.PhaseBeforeSize, nil)
	defer func() {
		w.PhaseRecorder.RecordSize(gtx.Metric.PxToDp(dims.Size.X), gtx.Metric.PxToDp(dims.Size.Y))
		w.PhaseRecorder.Trigger(layout.PhaseDidSize, nil)
	}()

	texture := paint.ColorTexture(gtx.Ops, w.Color)

	label := widget.Label{
		LineHeight:      w.LineHeight,
		LineHeightScale: 1,
	}

	if textAlign := w.TextAlign; textAlign != nil {
		label.Alignment = textAlign.TextAlignment()
	}

	gtx.Constraints.Min.Y = 0

	txt := paint.Group(gtx.Ops, func() {
		dims = label.Layout(gtx, w.shaper, w.Style.ToFont(), w.FontSize, w.textContents, texture)
	})

	fontSize, lineHeight := gtx.Sp(w.FontSize), gtx.Sp(w.LineHeight)

	///    xxxxxx
	///    x
	///    xxxxx
	///    x             -- baseline
	///    xxxxxx
	///
	///    xxxxxx
	///    x
	///    xxxxx
	///    x             -- baseline
	///    xxxxxx
	///
	/// line height here is height between baselines
	/// font size should contains accent + decent
	/// then the final height should be
	///		fontSize + lineHeight * lineCount
	///
	accent := dims.Size.Y - dims.Baseline
	lineCount := int(float32(dims.Size.Y-accent)/float32(lineHeight)) + 1
	fixedY := fontSize + lineHeight*(lineCount-1)

	defer op.Offset(image.Pt(0, dims.Size.Y-fixedY)).Push(gtx.Ops).Pop()
	txt.Add(gtx.Ops)

	dims.Size.Y = fixedY

	return dims
}
