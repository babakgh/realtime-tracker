package main_test

import (
	"reflect"
	"testing"

	. "github.com/babakgh/realtime-tracker"
)

func TestMisraGries_UpdateAndHeavyHitters(t *testing.T) {
	k := 3
	avg := int64(2)
	mg := NewMisraGries(k, avg)

	// Update with some elements
	mg.Update(1)
	mg.Update(1)
	mg.Update(2)
	mg.Update(3)
	mg.Update(3)
	mg.Update(3)
	mg.Update(4)

	expectedCounts := map[int]int64{
		3: 2,
	}

	heavyHitters := mg.HeavyHitters()

	if !reflect.DeepEqual(heavyHitters, expectedCounts) {
		t.Errorf("expected heavy Hitters to be %v, got %v", expectedCounts, heavyHitters)
	}

	// Update with more elements to test eviction logic
	mg.Update(5)
	mg.Update(5)
	mg.Update(5)
	mg.Update(5)

	expectedCountsAfterEviction := map[int]int64{
		3: 2,
		5: 4,
	}

	heavyHittersAfterEviction := mg.HeavyHitters()

	if !reflect.DeepEqual(heavyHittersAfterEviction, expectedCountsAfterEviction) {
		t.Errorf("expected heavy Hitters after eviction to be %v, got %v", expectedCountsAfterEviction, heavyHittersAfterEviction)
	}
}

func TestMisraGries_Concurrency(t *testing.T) {
	k := 3
	avg := int64(2)
	mg := NewMisraGries(k, avg)

	// Concurrent updates
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(element int) {
			for j := 0; j < 10; j++ {
				mg.Update(element)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to finish
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify counts
	heavyHitters := mg.HeavyHitters()
	for key, value := range heavyHitters {
		if value < avg {
			t.Errorf("expected element %d to have count >= %d, got %d", key, avg, value)
		}
	}
}
