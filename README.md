This document is a draft.

Given the high number of keys, which is much greater than the number of backend servers, the backend will be chosen based on the key—there is a kind of mapping involved.

Given that we know which backend is hotspotted, I chose to simulate a tracker for each hotspot backend in the upcoming time window—this is assumed to be 1 second. // Refactoring is needed.


MisraGriesTracker has been used to track keys and detect the heavy hitter key

The AccurateTracker has been used to track hotspotted backend, it is a usfull tools for testing, especially if we can use it without any issues in scenarios with a low number of buckets (here, backends). 
Generally, we can afford to use a bit more memory and CPU to track in smaller time windows quickly and accurately. This involves a trade-off between memory usage and real-time responsiveness (less near real-time, if that makes sense).

Will work with such a load balancer architecure https://github.com/babakgh/load-balancer
