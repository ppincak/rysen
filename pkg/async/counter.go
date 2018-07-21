package async

import "sync/atomic"

type Counter struct {
	value int64
}

func NewCounter() *Counter {
	return &Counter{
		value: 0,
	}
}

func (counter *Counter) Value() int64 {
	return atomic.LoadInt64(&counter.value)
}

func (counter *Counter) Add(value int64) {
	atomic.AddInt64(&counter.value, value)
}

func (counter *Counter) Inc() {
	counter.Add(1)
}

func (counter *Counter) Reset() {
	atomic.StoreInt64(&counter.value, 0)
}
