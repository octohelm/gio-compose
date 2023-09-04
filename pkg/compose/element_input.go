package compose

import (
	"context"
	"image"
	"image/color"

	"github.com/octohelm/gio-compose/pkg/modifier"

	"gioui.org/op"
	gtext "gioui.org/text"
	"gioui.org/widget"

	"github.com/octohelm/gio-compose/pkg/cmp"
	"github.com/octohelm/gio-compose/pkg/compose/internal"
	"github.com/octohelm/gio-compose/pkg/event/textinput"
	"github.com/octohelm/gio-compose/pkg/layout"
	"github.com/octohelm/gio-compose/pkg/node"
	"github.com/octohelm/gio-compose/pkg/paint"
	"github.com/octohelm/gio-compose/pkg/paint/f32color"
	"github.com/octohelm/gio-compose/pkg/text"
)

func Input(modifies ...modifier.Modifier[any]) VNode {
	return H(&inputElement{
		Element: node.Element{
			Name: "Input",
		},
	}, modifies...)
}

type inputElement struct {
	internal.ElementComponent

	node.Element
	inputWidgetAttr
	inputWidgetState
}

type inputWidgetState struct {
	shaper              *gtext.Shaper
	editor              widget.Editor
	submitEventProvider SubmitEventProvider
	text                string
	textEmitted         string
}

var _ TextSetter = &inputWidgetState{}

func (s *inputWidgetState) SetText(txt string) {
	s.text = txt
}

type inputWidgetAttr struct {
	text.Style
	*container
}

func (attrs *inputWidgetAttr) Eq(next *inputWidgetAttr) cmp.Result {
	return cmp.All(
		attrs.container.Eq(next.container),
		attrs.Style.Eq(&next.Style),
	)
}

func (w inputElement) New(ctx context.Context) internal.Element {
	return &inputElement{
		Element: node.Element{
			Name: w.Name,
		},
	}
}

func (w *inputElement) Update(ctx context.Context, modifiers ...modifier.Modifier[any]) bool {
	w.shaper = text.ShaperContext.Extract(ctx)
	w.submitEventProvider = SubmitEventProviderContext.Extract(ctx)

	modifier.Modify[any](&w.inputWidgetState, modifiers...)

	attrs := &inputWidgetAttr{
		container: newContainer(),
		Style:     TextStyleContext.Extract(ctx),
	}

	modifier.Modify[any](attrs, modifiers...)

	changed := cmp.UpdateWhen(
		cmp.Not(w.inputWidgetAttr.Eq(attrs)),
		&w.inputWidgetAttr, attrs,
	)

	if changed {
		if w.text != w.textEmitted {
			w.editor.SetText(w.text)
		}
	}

	return changed
}

func (w *inputElement) Layout(gtx layout.Context) layout.Dimensions {
	w.container.BindFocusedChecker(&w.editor)

	return w.container.Layout(gtx, w, ElementPainterFunc(func(gtx layout.Context) layout.Dimensions {
		for _, event := range w.editor.Events() {
			if _, ok := event.(widget.ChangeEvent); ok {
				w.textEmitted = w.editor.Text()
				w.InputEvents().Trigger(textinput.Change, &textinput.ChangeData{
					Value: w.textEmitted,
				})
			}

			if _, ok := event.(widget.SubmitEvent); ok {
				if w.submitEventProvider != nil {
					w.submitEventProvider.TriggerSubmit()
				}
			}
		}

		textColor := paint.ColorTexture(gtx.Ops, w.Color)
		hintColor := paint.ColorTexture(gtx.Ops, f32color.MulAlpha(w.Color, 0xbb))
		selectionColor := paint.ColorTexture(gtx.Ops, blendDisabledColor(gtx.Queue == nil, f32color.MulAlpha(w.Color, 0x33)))

		maxLines := 0
		if w.editor.SingleLine {
			maxLines = 1
		}

		w.editor.Submit = true

		w.editor.ReadOnly = w.GestureEvents().Disabled()

		// FIXME support multiline
		w.editor.SingleLine = true
		w.editor.LineHeight = w.LineHeight
		w.editor.LineHeightScale = 1

		if textAlign := w.TextAlign; textAlign != nil {
			w.editor.Alignment = textAlign.TextAlignment()
		}

		tl := widget.Label{
			Alignment:       w.editor.Alignment,
			LineHeight:      w.editor.LineHeight,
			LineHeightScale: w.editor.LineHeightScale,
			MaxLines:        maxLines,
		}

		var dims layout.Dimensions

		placeholder := paint.Group(gtx.Ops, func() {
			dims = tl.Layout(gtx, w.shaper, w.Style.ToFont(), w.Style.FontSize, "", hintColor)
		})

		editor := paint.Group(gtx.Ops, func() {
			dims = w.editor.Layout(gtx, w.shaper, w.Style.ToFont(), w.Style.FontSize, textColor, selectionColor)
		})

		if w.editor.SingleLine {
			halfLeading := (dims.Size.Y - gtx.Sp(w.FontSize)) / 2
			// middle align
			defer op.Offset(image.Pt(0, halfLeading)).Push(gtx.Ops).Pop()
		}

		editor.Add(gtx.Ops)
		if w.editor.Len() == 0 {
			placeholder.Add(gtx.Ops)
		}

		return dims
	}))
}

func blendDisabledColor(disabled bool, c color.NRGBA) color.NRGBA {
	if disabled {
		return f32color.Disabled(c)
	}
	return c
}
