package event

import (
	"context"

	"github.com/octohelm/gio-compose/pkg/iter"
)

type Events[T comparable] struct {
	disabled bool
	handlers map[T][]Handler[T]
}

func (events *Events[T]) Disable(disable bool) {
	events.disabled = disable
}

func (events *Events[T]) Disabled() bool {
	return events.disabled || len(events.handlers) == 0
}

func (events *Events[T]) Reset() {
	events.handlers = nil
}

func (events *Events[T]) Add(handlers ...Handler[T]) {
	if events.handlers == nil {
		events.handlers = map[T][]Handler[T]{}
	}

	for i := range handlers {
		h := handlers[i]
		if h == nil {
			continue
		}
		events.handlers[h.Type()] = append(events.handlers[h.Type()], h)
	}
}

func (events *Events[T]) Remove(handlers ...Handler[T]) {
	for i := range handlers {
		h := handlers[i]
		if h == nil || events.handlers == nil {
			continue
		}

		typ := h.Type()

		remainHandlers := iter.Filter(events.handlers[typ], func(e Handler[T]) bool {
			return e != h
		})

		if len(remainHandlers) != 0 {
			events.handlers[h.Type()] = remainHandlers
		} else {
			delete(events.handlers, typ)
		}
	}
}

func (events *Events[T]) Trigger(t T, data any) {
	if handlers, ok := events.handlers[t]; ok {
		for i := range handlers {
			if h := handlers[i]; h != nil {
				h.Handle(data)
			}
		}
	}
}

func (events *Events[T]) Watched(types ...T) bool {
	if handlers := events.handlers; handlers != nil {
		for _, g := range types {
			if _, ok := handlers[g]; ok {
				return true
			}
		}
	}
	return false
}

func (events *Events[T]) IterHandler(ctx context.Context, types ...T) <-chan Handler[T] {
	ch := make(chan Handler[T])

	go func() {
		defer close(ch)

		if handlers := events.handlers; handlers != nil {
			for _, t := range types {
				if handlersOfType, ok := handlers[t]; ok {
					for _, h := range handlersOfType {
						select {
						case <-ctx.Done():
							return
						case ch <- h:
						}
					}
				}
			}
		}
	}()

	return ch
}
