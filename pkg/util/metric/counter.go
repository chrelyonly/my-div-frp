package metric

import (
	"sync/atomic"
)

type Counter interface {
	Count() int32
	Inc(int32)
	Dec(int32)
	Snapshot() Counter
	Clear()
}

func NewCounter() Counter {
	return &StandardCounter{
		count: 0,
	}
}

type StandardCounter struct {
	count int32
}

func (c *StandardCounter) Count() int32 {
	return atomic.LoadInt32(&c.count)
}

func (c *StandardCounter) Inc(count int32) {
	atomic.AddInt32(&c.count, count)
}

func (c *StandardCounter) Dec(count int32) {
	atomic.AddInt32(&c.count, -count)
}

func (c *StandardCounter) Snapshot() Counter {
	tmp := &StandardCounter{
		count: atomic.LoadInt32(&c.count),
	}
	return tmp
}

func (c *StandardCounter) Clear() {
	atomic.StoreInt32(&c.count, 0)
}
