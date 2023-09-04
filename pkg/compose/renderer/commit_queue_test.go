package renderer

import (
	"testing"

	testingx "github.com/octohelm/x/testing"
)

func Test_actionQueueBufferEach(t *testing.T) {
	q := []func(){
		func() {},
		func() {},
		func() {},
		func() {},
		func() {},
	}

	bufLens := make([]int, 0)

	actionQueueBufferEach(q, 2, func(queue []func()) {
		bufLens = append(bufLens, len(queue))
	})

	testingx.Expect(t, bufLens, testingx.Equal([]int{
		2, 2, 1,
	}))
}
