package internal

type Hook interface {
	UpdateHook(next Hook)
}

type HookCommitter interface {
	Commit()
}

type HookAsyncCommitter interface {
	CommitAsync()
}

type HookDestroyer interface {
	Destroy()
}

type Hooks struct {
	hookUseIdx int
	usedHooks  []Hook
}

func (hs *Hooks) Reset() {
	hs.hookUseIdx = 0
}

func (hs *Hooks) putHook(hook Hook, idx int) Hook {
	if maxIdx := len(hs.usedHooks) - 1; maxIdx < hs.hookUseIdx {
		hs.usedHooks = append(hs.usedHooks, hook)
		return hs.usedHooks[idx]
	}
	return hs.usedHooks[idx]
}

func (hs *Hooks) use(hook Hook) Hook {
	usedHook := hs.putHook(hook, hs.hookUseIdx)
	usedHook.UpdateHook(hook)
	hs.hookUseIdx++
	return usedHook
}

func (hs *Hooks) commit() {
	for i := range hs.usedHooks {
		if hc, ok := hs.usedHooks[i].(HookCommitter); ok {
			hc.Commit()
		}
	}
}

func (hs *Hooks) destroy() {
	for i := range hs.usedHooks {
		if hc, ok := hs.usedHooks[i].(HookDestroyer); ok {
			hc.Destroy()
		}
	}
}

func (hs *Hooks) CommitAsync() {
	for i := range hs.usedHooks {
		if hc, ok := hs.usedHooks[i].(HookAsyncCommitter); ok {
			hc.CommitAsync()
		}
	}
}
