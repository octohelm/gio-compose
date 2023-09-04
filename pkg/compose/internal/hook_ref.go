package internal

import (
	"fmt"
)

type RefHook[T any] struct {
	Current T
}

func (RefHook[T]) UpdateHook(next Hook) {
	// do nothing here
}

func (s *RefHook[T]) String() string {
	return fmt.Sprintf("UseRef: %v", s.Current)
}
