package internal

import (
	"fmt"

	"github.com/octohelm/gio-compose/pkg/cmp"
)

type MemoHook[T any] struct {
	Setup     func() T
	Deps      []any
	memorised any
}

func (h *MemoHook[T]) Memorised() T {
	return h.memorised.(T)
}

func (h *MemoHook[T]) String() string {
	return fmt.Sprintf("UseMemo: %v", h.Deps)
}

func (h *MemoHook[T]) UpdateHook(next Hook) {
	if nextHook, ok := next.(*MemoHook[T]); ok {
		if nextHook == h || nextHook.Deps == nil || !cmp.ShallowEqual(nextHook.Deps, h.Deps) {
			h.memorised = nextHook.Setup()
		}
	}
}
