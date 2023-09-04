package modifier

import (
	"github.com/octohelm/gio-compose/pkg/event"
	"github.com/octohelm/gio-compose/pkg/layout/visible"
	"github.com/octohelm/gio-compose/pkg/modifier"
	"github.com/octohelm/gio-compose/pkg/paint"
)

func Visible(v bool) modifier.Modifier[any] {
	return &visibleModifier{visible: v}
}

type visibleModifier struct {
	visible bool
}

func (v *visibleModifier) Modify(target any) {
	if s, ok := target.(paint.VisibleSetter); ok {
		s.SetVisible(v.visible)
	}
}

func OnVisibleChange(fn func(v bool)) modifier.Modifier[any] {
	return &visibleChangeWatcher{
		action: event.NewHandlerWithEventData(visible.Change, func(data *visible.EventData) {
			fn(data.Visible)
		}),
	}
}

type visibleChangeWatcher struct {
	action event.Handler[visible.EventType]
}

func (w *visibleChangeWatcher) Modify(target any) {
	if p, ok := target.(visible.EventsAccessor); ok {
		p.VisibleEvents().Add(w.action)
	}
}
