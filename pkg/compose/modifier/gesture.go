package modifier

import (
	"github.com/octohelm/gio-compose/pkg/event/gesture"
	"github.com/octohelm/gio-compose/pkg/modifier"
)

func DetectGesture(handlers ...gesture.Handler) modifier.Modifier[any] {
	return &gestureHandlerModifiers{
		handlers: handlers,
	}
}

type gestureHandlerModifiers struct {
	handlers []gesture.Handler
}

func (m *gestureHandlerModifiers) String() string {
	return "[Gesture]"
}

func (g *gestureHandlerModifiers) Modify(w any) {
	if setter, ok := w.(gesture.EventsAccessor); ok {
		setter.GestureEvents().Add(g.handlers...)
	}
}

func GestureDisabled() modifier.Modifier[any] {
	return &gestureDisabledModifiers{}
}

type gestureDisabledModifiers struct{}

func (g gestureDisabledModifiers) Modify(w any) {
	if setter, ok := w.(gesture.EventsAccessor); ok {
		setter.GestureEvents().Disable(true)
	}
}
