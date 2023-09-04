package textinput

import (
	"gioui.org/io/key"

	"github.com/octohelm/gio-compose/pkg/event"
	"github.com/octohelm/gio-compose/pkg/node"
)

type Input int

const (
	KeyDown Input = iota + 10
	KeyUp
)

const (
	Change Input = iota + 20
)

type Handler = event.Handler[Input]

func WatchInputEvent(n node.Node, handlers ...Handler) func() {
	if d, ok := n.(EventsAccessor); ok {
		inputEvents := d.InputEvents()

		inputEvents.Add(handlers...)
		return func() {
			inputEvents.Remove(handlers...)
		}
	}

	return func() {}
}

func OnChange(action func(event *ChangeData)) Handler {
	return event.NewHandlerWithEventData(Change, action)
}

type ChangeData struct {
	Value string
}

type KeyData struct {
	Name      string
	Modifiers key.Modifiers
}

func OnKeyDown(keySet key.Set, action func()) Handler {
	return event.NewHandler(KeyDown, action, keySet)
}

func OnKeyUp(keySet key.Set, action func()) Handler {
	return event.NewHandler(KeyDown, action, keySet)
}
