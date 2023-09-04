package gesture

import (
	"gioui.org/io/key"
	"github.com/octohelm/gio-compose/pkg/event"
	"github.com/octohelm/gio-compose/pkg/node"
	"github.com/octohelm/gio-compose/pkg/unit"
)

type Gesture int

const (
	Press Gesture = iota + 10
	Tap
	Release

	DoubleTap
	LongPress
)

const (
	MouseEnter Gesture = iota + 20
	Hover
	MouseLeave
)

const (
	Focus Gesture = iota + 30
	Blur
)

type Handler = event.Handler[Gesture]

func WatchGestureEvent(n node.Node, handlers ...Handler) func() {
	if d, ok := n.(EventsAccessor); ok {
		events := d.GestureEvents()

		events.Add(handlers...)
		return func() {
			events.Remove(handlers...)
		}
	}

	return func() {}
}

func OnPress(action func()) Handler {
	return event.NewHandler(Press, action)
}

func OnPressWithEvent(action func(event *PointerData)) Handler {
	return event.NewHandlerWithEventData(Press, action)
}

func OnRelease(action func()) Handler {
	return event.NewHandler(Release, action)
}

func OnReleaseWithEvent(action func(p *PointerData)) Handler {
	return event.NewHandlerWithEventData(Release, action)
}

func OnTap(action func()) Handler {
	return event.NewHandler(Tap, action)
}

func OnTapWithEvent(action func(event *PointerData)) Handler {
	return event.NewHandlerWithEventData(Tap, action)
}

func OnDoubleTap(action func()) Handler {
	return event.NewHandler(DoubleTap, action)
}

func OnLongPress(action func()) Handler {
	return event.NewHandler(LongPress, action)
}

func OnMouseEnter(action func()) Handler {
	return event.NewHandler(MouseEnter, action)
}

func OnMouseLeave(action func()) Handler {
	return event.NewHandler(MouseLeave, action)
}

func OnHover(action func()) Handler {
	return event.NewHandler(Hover, action)
}

func OnFocus(action func()) Handler {
	return event.NewHandler(Focus, action)
}

func OnBlur(action func()) Handler {
	return event.NewHandler(Blur, action)
}

type PointerData struct {
	// Position related of element
	Position  unit.Point
	Modifiers key.Modifiers
	Target    node.Node
}
