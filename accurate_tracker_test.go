package main_test

import (
	"reflect"
	"sync/atomic"
	"testing"

	. "github.com/babakgh/kata/play-hash-table"
)

func TestAccurateTracker_Update(t *testing.T) {
	count := 5
	at := NewAccurateTracker(count)

	// Update a specific bucket and check the result
	bucketToUpdate := 2
	at.Update(bucketToUpdate)

	if atomic.LoadInt64(&at.Buckets[bucketToUpdate]) != 1 {
		t.Errorf("expected bucket %d to have count 1, got %d", bucketToUpdate, at.Buckets[bucketToUpdate])
	}
}

func TestAccurateTracker_LoadHotspotsAndReset(t *testing.T) {
	count := 5
	at := NewAccurateTracker(count)

	// Update buckets to simulate data
	at.Update(1)
	at.Update(1)
	at.Update(3)
	at.Update(3)
	at.Update(3)

	// Parameters for LoadHotspotsAndReset
	avg := 1.0
	window := 1.0
	threshold := 2.0

	expectedHeavyHitters := map[int]int64{
		1: 2,
		3: 3,
	}

	heavyHitters := at.LoadHotspotsAndReset(avg, window, threshold)

	if !reflect.DeepEqual(heavyHitters, expectedHeavyHitters) {
		t.Errorf("expected heavy hitters to be %v, got %v", expectedHeavyHitters, heavyHitters)
	}

	// Check if buckets were reset
	for _, bucket := range at.Buckets {
		if bucket != 0 {
			t.Errorf("expected all buckets to be reset to 0, found %d", bucket)
		}
	}
}

func TestAccurateTracker_Concurrency(t *testing.T) {
	count := 5
	at := NewAccurateTracker(count)

	// Concurrent updates
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < count; j++ {
				at.Update(j)
			}
			done <- true
		}()
	}

	// Wait for all goroutines to finish
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify counts are as expected (10 updates per bucket)
	for i := 0; i < count; i++ {
		expectedCount := int64(10)
		if atomic.LoadInt64(&at.Buckets[i]) != expectedCount {
			t.Errorf("expected bucket %d to have count %d, got %d", i, expectedCount, at.Buckets[i])
		}
	}
}
