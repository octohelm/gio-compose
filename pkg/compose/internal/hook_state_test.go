package internal

import (
	"testing"

	testingx "github.com/octohelm/x/testing"
)

func TestStateHook(t *testing.T) {
	s := &StateHook[int]{}

	s.Mount()

	for i := 0; i < 3; i++ {
		s.UpdateFunc(func(prev int) int {
			return prev + 1
		})
	}

	s.CommitAsync()

	testingx.Expect(t, s.State, testingx.Be(3))
}
