package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

type CountLimit struct {
	counter      int64
	limit        int64
	intervalNano int64
	unixNano     int64
}

func NewCountLimit(interval time.Duration, limit int64) *CountLimit {
	return &CountLimit{
		counter:      0,
		limit:        limit,
		intervalNano: int64(interval),
		unixNano:     time.Now().UnixNano(),
	}
}

func (c *CountLimit) Allow() bool {
	now := time.Now().UnixNano()
	if now-c.unixNano > c.intervalNano {
		atomic.StoreInt64(&c.counter, 0)
		atomic.StoreInt64(&c.unixNano, now)
		return true
	}

	atomic.AddInt64(&c.counter, 1)
	fmt.Println(c.counter)
	return c.counter < c.limit
}

func main() {
	limit := NewCountLimit(time.Nanosecond*100000, 100)
	m := make(map[int]bool)
	for i := 0; i < 100; i++ {
		allow := limit.Allow()
		if allow {
			//fmt.Printf("i=%d is allow\n", i)
			m[i] = true
		} else {
			//fmt.Printf("i=%d is not allow\n", i)
			m[i] = false
		}
	}

}
