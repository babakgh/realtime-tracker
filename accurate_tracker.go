package main

import (
	"sync/atomic"
)

type AccurateTracker struct {
	count   int // number of buckets
	Buckets []int64
}

func NewAccurateTracker(count int) *AccurateTracker {
	a := &AccurateTracker{
		count: count,
	}
	a.Buckets = make([]int64, count)
	for i := range a.Buckets {
		atomic.StoreInt64(&a.Buckets[i], 0)
	}

	return a
}

func (at *AccurateTracker) Update(bucket int) {
	atomic.AddInt64(&at.Buckets[bucket], 1)
}

// The idea is to reset counters after each window (per second) and store the hotsopts
func (at *AccurateTracker) LoadHotspotsAndReset(avg, window, threshold float64) map[int]int64 {
	hh := map[int]int64{} // creating a copy
	for i := range at.Buckets {
		count := atomic.SwapInt64(&at.Buckets[i], 0) // Get data and reset
		if float64(count) >= threshold*avg*window {
			hh[i] = count
		}
	}

	return hh
}
