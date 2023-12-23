package tunnel

import (
	"sync"
	"sync/atomic"
)

type waitgroup struct {
	inner sync.WaitGroup
	n     int32
}

func (w *waitgroup) Add(n int) {
	w.inner.Add(n)
	atomic.AddInt32(&w.n, int32(n))
}
func (w *waitgroup) Done() {
	if n := atomic.LoadInt32(&w.n); n > 0 && atomic.CompareAndSwapInt32(&w.n, n, n-1) {
		w.inner.Done()
	}
}

func (w *waitgroup) DoneAll() {
	for atomic.LoadInt32(&w.n) > 0 {
		w.Done()
	}
}
func (w *waitgroup) Wait() {
	w.inner.Wait()
}
