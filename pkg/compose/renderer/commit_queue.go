package renderer

import (
	"sync/atomic"
	"time"

	"github.com/octohelm/x/ptr"
)

type commitQueue struct {
	tasks atomic.Pointer[[]func()]
	timer *time.Timer
	done  int32
}

func (q *commitQueue) push(fn func()) {
	if tasks := q.tasks.Load(); tasks != nil {
		q.tasks.Store(ptr.Ptr(append(*tasks, fn)))
	} else {
		q.tasks.Store(ptr.Ptr([]func(){fn}))
	}
}

func (q *commitQueue) runIfExists() {
	tasks := q.tasks.Swap(nil)

	if tasks != nil {
		actionQueueBufferEach(*tasks, bufSize, func(buf []func()) {
			for i := range buf {
				buf[i]()
			}
		})
	}
}

const bufSize = 2048
const bufferTime = 10 * time.Millisecond

func actionQueueBufferEach(queue []func(), bufSize int, each func(queue []func())) {
	n := len(queue)
	partN := n / bufSize
	for i := 0; i < partN+1; i++ {
		left := i * bufSize
		if left > n-1 {
			break
		}
		right := (i + 1) * bufSize
		if right > n {
			right = n
		}
		each(queue[i*bufSize : right])
	}
}

func (q *commitQueue) ForceCommit() {
	q.runIfExists()
}

func (q *commitQueue) Close() error {
	atomic.StoreInt32(&q.done, 1)
	return nil
}

func (q *commitQueue) Start() {
	q.timer = time.NewTimer(bufferTime)

	go func() {
		for {
			<-q.timer.C

			if atomic.LoadInt32(&q.done) == 0 {
				q.timer.Reset(bufferTime)
			}

			q.runIfExists()
		}
	}()
}

func (q *commitQueue) Commit(fn func()) {
	q.push(fn)
}
