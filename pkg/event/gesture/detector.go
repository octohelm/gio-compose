package gesture

import (
	"image"

	"github.com/octohelm/gio-compose/pkg/node"
	"github.com/octohelm/gio-compose/pkg/paint"
	"github.com/octohelm/gio-compose/pkg/unit"

	"gioui.org/gesture"
	"gioui.org/io/pointer"
	giolayout "gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/widget"
	"github.com/octohelm/gio-compose/pkg/event"
	"github.com/octohelm/gio-compose/pkg/layout"
)

type EventsAccessor interface {
	GestureEvents() *event.Events[Gesture]
}

type FocusedChecker interface {
	Focused() bool
}

var _ EventsAccessor = &PointerEventDetector{}

type PointerEventDetector struct {
	// TODO impl Focus when tabindex exists
	// https://developer.mozilla.org/en-US/docs/Web/HTML/Global_attributes/tabindex
	// key.FocusOp{Tag: &b.keyTag}.Add(gtx.Ops)
	tabIndex int
	widget.Clickable

	click gesture.Click

	lastMouseEventData *PointerData

	events event.Events[Gesture]

	pressed State[bool]
	hovered State[bool]
	focused State[bool]

	focusedChecker FocusedChecker
}

func (d *PointerEventDetector) BindFocusedChecker(fc FocusedChecker) {
	d.focusedChecker = fc
}

func (d *PointerEventDetector) GestureEvents() *event.Events[Gesture] {
	return &d.events
}

func (d *PointerEventDetector) LayoutChild(gtx giolayout.Context, n node.Node, child giolayout.Widget) (dims layout.Dimensions) {
	if d.events.Disabled() {
		return child(gtx)
	}

	clicked := false
	for _, e := range d.click.Events(gtx) {
		switch e.Type {
		case gesture.TypePress:
			d.lastMouseEventData = &PointerData{
				Position:  unit.PtFromPx(gtx.Metric, e.Position.X, e.Position.Y),
				Modifiers: e.Modifiers,
				Target:    n,
			}
			d.pressed.Set(true)
		case gesture.TypeClick:
			clicked = true
			d.pressed.Set(false)
		case gesture.TypeCancel:
			d.pressed.Set(false)
		}
	}

	g := paint.Group(gtx.Ops, func() {
		dims = child(gtx)
	})
	defer g.Add(gtx.Ops)

	defer clip.Rect(image.Rectangle{Max: dims.Size}).Push(gtx.Ops).Pop()
	d.click.Add(gtx.Ops)

	// trigger events

	d.hovered.Set(d.click.Hovered())

	if focusedChecker := d.focusedChecker; focusedChecker != nil {
		d.focused.Set(focusedChecker.Focused())
	} else {
		//
	}

	if d.pressed.Changed {
		if d.pressed.Value {
			d.events.Trigger(Press, d.lastMouseEventData)
		} else {
			if clicked {
				d.events.Trigger(Tap, d.lastMouseEventData)
			}
			d.events.Trigger(Release, d.lastMouseEventData)
		}
	}

	if d.click.Hovered() {
		if d.events.Watched(Press, Tap, DoubleTap, LongPress) {
			layout.PostLayout(gtx.Ops, pointer.CursorPointer.Add)
		}
	} else {
		layout.Layout(gtx.Ops, pointer.CursorDefault.Add)
	}

	if d.hovered.Changed {
		if d.hovered.Value {
			d.events.Trigger(MouseEnter, nil)
		} else {
			d.events.Trigger(MouseLeave, nil)
		}
	}

	if d.focused.Changed {
		if d.focused.Value {
			d.events.Trigger(Focus, nil)
		} else {
			d.events.Trigger(Blur, nil)
		}
	}

	return
}

type State[T comparable] struct {
	Value   T
	Changed bool
}

func (once *State[T]) Set(v T) {
	if once.Value != v {
		once.Value = v
		once.Changed = true
	}
}
