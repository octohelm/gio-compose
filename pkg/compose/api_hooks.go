package compose

import (
	"github.com/octohelm/gio-compose/pkg/compose/internal"
)

type State[T comparable] interface {
	Value() T
	Update(value T)
	UpdateFunc(func(prev T) T)
}

func UseState[T comparable](ctx BuildContext, defaultState T) State[T] {
	bc := ctx.(internal.BuildContextAccessor)
	vn := bc.VNode()

	return internal.UseHook(vn, &internal.StateHook[T]{
		State: defaultState,
		OnStateChange: func() {
			// use raw context to avoid context stack overflow
			vn.Update(bc.RawContext())
		},
	})
}

func UseEffect(ctx BuildContext, setup func() func(), deps []any) {
	vn := ctx.(internal.BuildContextAccessor).VNode()

	internal.UseHook(vn, &internal.EffectHook{
		Setup: setup,
		Deps:  deps,
	})
}

func UseMemo[T any](ctx BuildContext, setup func() T, deps []any) T {
	vn := ctx.(internal.BuildContextAccessor).VNode()

	h := internal.UseHook(vn, &internal.MemoHook[T]{
		Setup: setup,
		Deps:  deps,
	})

	return h.Memorised()
}

type Ref[T any] struct {
	Current T
}

func UseRef[T any](ctx BuildContext, initialValue T) *Ref[T] {
	vn := ctx.(internal.BuildContextAccessor).VNode()

	h := internal.UseHook(vn, &internal.RefHook[*Ref[T]]{
		Current: &Ref[T]{Current: initialValue},
	})

	return h.Current
}
