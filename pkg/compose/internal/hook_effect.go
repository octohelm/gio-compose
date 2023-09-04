package internal

import (
	"fmt"

	"github.com/octohelm/gio-compose/pkg/cmp"
)

type EffectHook struct {
	Setup   func() func()
	Deps    []any
	cleanup func()
	commit  func()
}

func (h *EffectHook) String() string {
	return fmt.Sprintf("UseEffect: %v", h.Deps)
}

var _ HookDestroyer = &EffectHook{}

func (h *EffectHook) Destroy() {
	if h.cleanup != nil {
		h.cleanup()
		h.cleanup = nil
	}
	h.commit = nil
}

func (h *EffectHook) Commit() {
	if h.commit != nil {
		h.commit()
	}
}

func (h *EffectHook) UpdateHook(next Hook) {
	nextHook := next.(*EffectHook)

	if nextHook == h || nextHook.Deps == nil || !cmp.ShallowEqual(nextHook.Deps, h.Deps) {
		h.commit = func() {
			if h.cleanup != nil {
				h.cleanup()
			}
			h.cleanup = nextHook.Setup()
			h.Deps = nextHook.Deps
		}
	} else {
		// avoid call if not changes
		h.commit = nil
	}
}
