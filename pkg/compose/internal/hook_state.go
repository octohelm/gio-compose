package internal

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/octohelm/x/ptr"
)

type StateHook[T comparable] struct {
	State         T
	OnStateChange func()

	timer    atomic.Pointer[time.Timer]
	reducers atomic.Pointer[[]func(prev T) T]
	mounted  int64
}

func (s *StateHook[T]) String() string {
	return fmt.Sprintf("UseState[%v]", s.State)
}

func (s *StateHook[T]) Value() T {
	return s.State
}

func (s *StateHook[T]) CommitAsync() {
	if reducers := s.reducers.Load(); reducers != nil {
		state := s.State

		for _, r := range *reducers {
			state = r(state)
		}

		if state != s.State {
			s.State = state
			if s.OnStateChange != nil {
				s.OnStateChange()
			}
		}

		s.reducers.Store(nil)

		s.CleanTimer()
	}
}

const timeBuffer = 10 * time.Millisecond

func (s *StateHook[T]) UpdateFunc(fn func(prev T) T) {
	if atomic.LoadInt64(&s.mounted) != 1 {
		return
	}

	if reducers := s.reducers.Load(); reducers != nil {
		s.reducers.Store(ptr.Ptr(append(*reducers, fn)))
	} else {
		s.reducers.Store(ptr.Ptr([]func(prev T) T{fn}))
	}

	if t := s.timer.Load(); t == nil {
		s.timer.Store(time.AfterFunc(timeBuffer, s.CommitAsync))
	} else {
		t.Reset(timeBuffer)
	}
}

func (s *StateHook[T]) Commit() {
	s.Mount()
	s.CommitAsync()
}

func (s *StateHook[T]) Mount() {
	atomic.StoreInt64(&s.mounted, 1)
}

func (s *StateHook[T]) CleanTimer() {
	if t := s.timer.Load(); t != nil {
		t.Stop()
		s.timer.Store(nil)
	}
}

func (s *StateHook[T]) Destroy() {
	atomic.StoreInt64(&s.mounted, 0)

	s.CleanTimer()

	s.reducers.Store(nil)
}

func (s *StateHook[T]) Update(v T) {
	s.UpdateFunc(func(prev T) T {
		return v
	})
}

func (s *StateHook[T]) UpdateHook(next Hook) {
	if n, ok := next.(*StateHook[T]); ok {
		// context may change, should bind the latest callback
		s.OnStateChange = n.OnStateChange
	}
}
