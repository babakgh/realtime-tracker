package main

import "sync"

// MisraGriesTracker keeps track of frequent elements in a thread-safe manner
type MisraGriesTracker struct {
	n      int           // number of heavy hitters to track
	avg    int64         // average threshold
	counts map[int]int64 // map of hiiter  to count
	mutex  sync.RWMutex
}

// NewMisraGries initializes the MisraGries sketch with k counters
func NewMisraGries(n int, avg int64) *MisraGriesTracker {
	return &MisraGriesTracker{
		n:      n,
		avg:    avg,
		counts: make(map[int]int64),
	}
}

// Process takes a stream element and updates the counters in a thread-safe manner
func (mg *MisraGriesTracker) Update(element int) {
	mg.mutex.Lock()
	defer mg.mutex.Unlock()

	// If the element is already being tracked, increment its count
	if _, exists := mg.counts[element]; exists {
		mg.counts[element]++
		return
	}

	// If there's still space for new elements, track this one
	if len(mg.counts) < mg.n {
		mg.counts[element] = 1
		return
	}

	// Otherwise, decrement all counts and remove any that reach zero
	for e := range mg.counts {
		mg.counts[e]--
		if mg.counts[e] == 0 {
			delete(mg.counts, e)
		}
	}
}

// HeavyHitters returns the current counts of tracked elements in a thread-safe manner
func (mg *MisraGriesTracker) HeavyHitters() map[int]int64 {
	mg.mutex.Lock()
	defer mg.mutex.Unlock()

	// Return a copy of the map to avoid concurrent read/write issues
	copy := make(map[int]int64)
	for key, value := range mg.counts {
		if value >= mg.avg { // add to heavy hitters
			copy[key] = value
		}
	}
	return copy
}
