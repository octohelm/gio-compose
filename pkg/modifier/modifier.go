package modifier

import "context"

type Modifier[T any] interface {
	Modify(target T)
}

func Modify[T any](target T, modifiers ...Modifier[T]) {
	for i := range modifiers {
		m := modifiers[i]
		if m != nil {
			m.Modify(target)
		}
	}
}

type Iterator interface {
	Iter(ctx context.Context) <-chan Modifier[any]
}

func Iter[T any](ctx context.Context, modifiers ...Modifier[T]) <-chan Modifier[any] {
	ch := make(chan Modifier[any])

	go func() {
		defer close(ch)

		for _, m := range modifiers {
			switch x := m.(type) {
			case Iterator:
				for fc := range x.Iter(ctx) {
					select {
					case <-ctx.Done():
						return
					case ch <- fc:
					}
				}
				continue
			}

			select {
			case <-ctx.Done():
				return
			case ch <- &wrapper[T]{m: m}:
			}
		}
	}()

	return ch
}

func ModifiersOf[T any](modifiers ...Modifier[T]) Modifiers {
	ret := make(Modifiers, len(modifiers))
	for i := range modifiers {
		ret[i] = &wrapper[T]{m: modifiers[i]}
	}
	return ret
}

type Modifiers []Modifier[any]

func (modifiers Modifiers) Iter(ctx context.Context) <-chan Modifier[any] {
	return Iter[any](ctx, modifiers...)
}

func (modifiers Modifiers) Modify(target any) {
	Modify(target, modifiers...)
}

type Discord struct {
}

func (Discord) Modify(target any) {
}

func Unwrap(a any) any {
	if x, ok := a.(interface{ Unwrap() any }); ok {
		return x.Unwrap()
	}
	return a
}

type wrapper[T any] struct {
	m Modifier[T]
}

func (w *wrapper[T]) Unwrap() any {
	return w.m
}

func (w *wrapper[T]) Modify(target any) {
	if w.m != nil {
		if x, ok := target.(T); ok {
			w.m.Modify(x)
		}
	}
}
