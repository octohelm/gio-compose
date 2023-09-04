package modifier

import "context"

func When[T any](when bool, modifiers ...Modifier[T]) Modifier[T] {
	return &conditionModifier[T]{when: when, modifiers: modifiers}
}

type conditionModifier[T any] struct {
	when      bool
	modifiers []Modifier[T]
}

func (m *conditionModifier[T]) Modify(target T) {
	if m.when {
		Modify(target, m.modifiers...)
	}
}

func (m *conditionModifier[T]) Iter(ctx context.Context) <-chan Modifier[any] {
	if m.when {
		return Iter(ctx, m.modifiers...)
	}
	return closedModifierCh
}

var closedModifierCh = make(chan Modifier[any])

func init() {
	close(closedModifierCh)
}
