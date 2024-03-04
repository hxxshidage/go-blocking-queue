package tpool

import (
	"sync"
	"time"
)

type TimerPool interface {
	Get(time.Duration) *time.Timer
	Put(*time.Timer)
}

type timerPool struct {
	timerPool sync.Pool
}

func (this *timerPool) Get(timeout time.Duration) *time.Timer {
	t := this.timerPool.Get().(*time.Timer)
	t.Reset(timeout)
	return t
}

func (this *timerPool) Put(t *time.Timer) {
	t.Stop()
	this.timerPool.Put(t)
}

var TimerRecycler = NewTickerPool()

func NewTickerPool() TimerPool {
	return &timerPool{
		timerPool: sync.Pool{
			New: func() interface{} {
				t := time.NewTimer(time.Hour)
				t.Stop()
				return t
			},
		},
	}
}
