package textinput

import (
	"context"

	"gioui.org/io/key"
	giolayout "gioui.org/layout"
	"github.com/octohelm/gio-compose/pkg/event"
)

type EventsAccessor interface {
	InputEvents() *event.Events[Input]
}

type InputEventDetector struct {
	events event.Events[Input]
}

func (d *InputEventDetector) InputEvents() *event.Events[Input] {
	return &d.events
}

func (d *InputEventDetector) Layout(gtx giolayout.Context) {
	if d.events.Disabled() {
		return
	}

	var keySet key.Set

	for h := range d.events.IterHandler(context.Background(), KeyUp, KeyDown) {
		for _, m := range h.Metadata() {
			if keySetFromMeta, ok := m.(key.Set); ok {
				keySet += keySetFromMeta
			}
		}
	}

	if keySet != "" {
		for _, e := range gtx.Events(d) {
			switch x := e.(type) {
			case key.Event:
				switch x.State {
				case key.Press:
					for h := range d.events.IterHandler(context.Background(), KeyDown) {
						for _, m := range h.Metadata() {
							if keySetFromMeta, ok := m.(key.Set); ok {
								if keySetFromMeta.Contains(x.Name, x.Modifiers) {
									h.Handle(nil)
								}
							}
						}
					}
				case key.Release:
					for h := range d.events.IterHandler(context.Background(), KeyUp) {
						for _, m := range h.Metadata() {
							if keySetFromMeta, ok := m.(key.Set); ok {
								if keySetFromMeta.Contains(x.Name, x.Modifiers) {
									h.Handle(nil)
								}
							}
						}
					}
				}
			}
		}

		key.InputOp{Tag: d, Keys: keySet}.Add(gtx.Ops)
	}
}
