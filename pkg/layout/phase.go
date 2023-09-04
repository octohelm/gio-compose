package layout

import (
	"github.com/octohelm/gio-compose/pkg/event"
	"github.com/octohelm/gio-compose/pkg/unit"
)

type Phase int

const (
	PhaseBeforeSize Phase = iota
	PhaseDidSize

	PhaseBeforePosition
	PhaseDidPosition
)

type PhaseHandler = event.Handler[Phase]

func OnBeforeSize(action func()) PhaseHandler {
	return event.NewHandler(PhaseBeforeSize, action)
}

func OnDidSize(action func()) PhaseHandler {
	return event.NewHandler(PhaseDidSize, action)
}

func OnBeforePosition(action func()) PhaseHandler {
	return event.NewHandler(PhaseBeforePosition, action)
}

func OnDidPosition(action func()) PhaseHandler {
	return event.NewHandler(PhaseDidPosition, action)
}

type PhaseHandlersSetter interface {
	SetPhaseHandlers(handlers ...PhaseHandler)
}

type PhaseRecorder struct {
	event.Events[Phase]

	boundingRect BoundingRect
}

var _ PhaseHandlersSetter = &PhaseRecorder{}

func (d *PhaseRecorder) SetPhaseHandlers(handlers ...PhaseHandler) {
	d.Add(handlers...)
}

func (d *PhaseRecorder) BoundingRect() BoundingRect {
	return d.boundingRect
}

func (d *PhaseRecorder) PositionBy(calc func() (x unit.Dp, y unit.Dp)) {
	d.Trigger(PhaseBeforePosition, nil)
	x, y := calc()
	d.RecordPosition(x, y)
	d.Trigger(PhaseDidPosition, nil)
}

func (d *PhaseRecorder) RecordSize(w unit.Dp, h unit.Dp) {
	d.boundingRect.Width, d.boundingRect.Height = w, h
}

func (d *PhaseRecorder) RecordPosition(x unit.Dp, y unit.Dp) {
	d.boundingRect.X, d.boundingRect.Y = x, y
}
