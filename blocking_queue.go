package bqueue

import (
	tpool "github.com/hxxshidage/go-blocking-queue/internal"
	"sync/atomic"
	"time"
)

type BlockQueue interface {
	/* Insert elements into the queue and block when the channel capacity is full. */
	Put(ele interface{}) bool

	/* Insert elements into the queue. When the channel capacity is full, it will block until timeout. */
	Offer(ele interface{}, timeout time.Duration) bool

	/* Get and remove the head element, when the channel is empty, it will block */
	Take() interface{}

	/* Get and remove the head element, when the channel is empty, it will block until timeout */
	Poll(timeout time.Duration) (interface{}, bool)

	/* queue len */
	Len() uint32

	/* queue capacity */
	Cap() uint32
}

type blockingQueue struct {
	buffer    chan interface{}
	capacity  uint32
	len       uint32
	timerPool tpool.TimerPool // timer to implement timeout function
}

/* Create and specify the capacity, which is actually the capacity of the channel */
func NewBlockingQueue(capacity uint32) BlockQueue {
	if capacity == 0 {
		panic("capacity illegal")
	}
	return &blockingQueue{
		capacity:  capacity,
		buffer:    make(chan interface{}, capacity),
		timerPool: tpool.TimerRecycler, // a pool to reuse timer
		len:       0,
	}
}

func (this *blockingQueue) Offer(ele interface{}, timeout time.Duration) bool {
	t := this.timerPool.Get(timeout)
	select {
	case <-t.C:
		this.timerPool.Put(t)
		return false
	case this.buffer <- ele:
		atomic.AddUint32(&this.len, 1)
		this.timerPool.Put(t)
		return true
	}
}

func (this *blockingQueue) Put(ele interface{}) bool {
	this.buffer <- ele
	atomic.AddUint32(&this.len, 1)
	return true
}

func (this *blockingQueue) Take() interface{} {
	val := <-this.buffer
	atomic.AddUint32(&this.len, ^uint32(1-1))
	return val
}

func (this *blockingQueue) Poll(timeout time.Duration) (interface{}, bool) {
	t := this.timerPool.Get(timeout)
	select {
	case <-t.C:
		return nil, false // timeout, return false
	case val := <-this.buffer:
		atomic.AddUint32(&this.len, ^uint32(1-1))
		this.timerPool.Put(t)
		return val, true
	}
}

func (this *blockingQueue) Len() uint32 {
	return this.len
}

func (this *blockingQueue) Cap() uint32 {
	return this.capacity
}
