package visible

import "github.com/octohelm/gio-compose/pkg/event"

type EventsAccessor interface {
	VisibleEvents() *event.Events[EventType]
}

type EventData struct {
	Visible bool
}

type EventType int

const (
	Change EventType = iota + 10
)
