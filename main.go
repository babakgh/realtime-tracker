package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/trace"
	"sync"
	"time"
)

var requests chan MockRequest

// var mg *MisraGriesTracker
var hotspotBackendTracker *AccurateTracker
var heavyHitterKeyTracker map[int]*MisraGriesTracker

func main() {
	f, _ := os.Create("trace.out")
	trace.Start(f)
	defer trace.Stop()

	p := 1 // processes
	w := 1 // workers

	runtime.GOMAXPROCS(p)

	// Start workers to process first window, 1 second
	quit := make(chan bool)
	requests = make(chan MockRequest)
	wg := &sync.WaitGroup{}
	for i := 0; i < w; i++ {
		wg.Add(1)
		go worker(i, wg, quit)
	}

	// Simulate requests
	U := 60_000 // requests per second, universe
	m := 100    // hotspot multiplier
	n := 100    // backends
	k := 1      // hotspots backends
	optimum_avg_requests_per_backend := int64(U / n)

	// tracker to find the hotsopt backend
	hotspotBackendTracker = NewAccurateTracker(n)
	heavyHitterKeyTracker := make(map[int]*MisraGriesTracker)

	// Calculate average number of requests per backend if there are k hotspots and n-k non-hotspots
	// T = m * k * hotspot_avg_requests_per_backend + (n - k) * hotspot_avg_requests_per_backend
	hotspot_avg_requests_per_backend := U / (m*k + n - k)

	// Start 1st window
	simulateWindow(1, m, n, k, hotspot_avg_requests_per_backend)

	// when we found the hotsopt backend, we are going to simulate the second window which aim to track the heavy hitters keys
	a := float64(U / n) // average requests per backend
	t := 1.0            // window in seconds
	e := 1.2            // threshold multiplier

	hotspots := hotspotBackendTracker.LoadHotspotsAndReset(a, t, e)
	fmt.Printf("Heavy Hitters using Accurate tracker %v\n", hotspots)
	if len(hotspots) > k {
		fmt.Printf("Test Failed: want %d, got %d hotspots\n", k, len(hotspots))
		return
	}

	keys := 1_000_000           //  potential keys per backend
	keys_to_track := keys / 100 // track 1% of keys

	for _, hotspotBackendID := range hotspots {
		// create a tracker for each hotspot backend
		heavyHitterKeyTracker[int(hotspotBackendID)] = NewMisraGries(keys_to_track, optimum_avg_requests_per_backend)
	}

	// Start workers to process second window, 1 second
	// so now we know the hotspot backend, we can use Misra-Gries tracker to find the heavy hitters
	// Start the second window
	simulateWindow(2, m, n, k, hotspot_avg_requests_per_backend)

	// Shutdown workers gracefully
	for j := 0; j < w; j++ {
		quit <- true
	}

	wg.Wait()

	// free resources
	close(requests)
	close(quit)
}

func worker(id int, wg *sync.WaitGroup, quit chan bool) {
	defer wg.Done()

	for {
		select {
		case <-quit:
			return
		case b := <-requests:
			hotspotBackendTracker.Update(b.BackendID)
			// mg.Update(b.BackendID)

			// time.Sleep(1 * time.Millisecond) // Debug: Do some work
		}
	}
}

// Generate requests and queue them up using channels
// channels are used to simulation and in real production we would have a better options
// this should take less than 1 second, which simulates a window of streaming data

// m: hotspot multiplier
// n: backends count
// k: hotspots backends count
// a: average requests per backend
func simulateWindow(id, m, n, k, a int) {
	now := time.Now()
	for b := k; b < n; b++ { // non-hotspots
		for j := 1; j <= a; j++ {
			requests <- *NewMockRequest("any-key", b)
		}
	}
	for b := 0; b < k; b++ { // hotspots
		hits := float32(a * m)
		for j := 1; j <= int(0.8*hits); j++ {
			requests <- *NewMockRequest("key-0", b)
		}
		for j := 1; j <= int(0.2*hits); j++ {
			key := fmt.Sprintf("key-%d", j) // create lots of keys to test the trackers with so many keys
			requests <- *NewMockRequest(key, b)
		}
	}
	duration := time.Since(now)
	fmt.Printf("It took %v\n", duration)
	if duration > 50*time.Millisecond { // Test latancy
		fmt.Printf("Benchmark Failed: Tracker is too slow\n")
	}
	// End 2nd window
}
