package debounce

import (
	"time"

	"sync/atomic"
)

type Debouncing struct {
	BufferTime time.Duration

	timer atomic.Pointer[time.Timer]
}

func (d *Debouncing) Do(fn func()) {
	if t := d.timer.Load(); t != nil {
		t.Stop()
	}

	d.timer.Store(time.AfterFunc(d.BufferTime, func() {
		fn()
		d.timer.Store(nil)
	}))
}
