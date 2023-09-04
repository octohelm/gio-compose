package debounce

import (
	"testing"
	"time"

	testingx "github.com/octohelm/x/testing"
)

func TestDebouncing(t *testing.T) {

	d := Debouncing{BufferTime: 10 * time.Millisecond}

	ret := 0

	for i := 1; i <= 10; i++ {
		d.Do(func(i int) func() {
			return func() {
				ret = i
			}
		}(i))
	}

	time.Sleep(d.BufferTime * 2)
	testingx.Expect(t, ret, testingx.Be(10))
}
