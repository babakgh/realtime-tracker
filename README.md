- This document is draft

// Arange
Given the number of keys is high, much more than the number of backend servers. Then based on the key
Then the backend will be chose base on the key - there is a kind of map somewhere.

Given we do know which backend is hotspotted - I chose to simulated
Then A tracker for each hotspot backend in next comming time window - Here is assumed for 1 second // refactor is needed

// Action

//Assert

// Tools
AccurateTracker can be useful on testing, also if we can use it without any problem, in case of low amount of buckets (here backends).
Generally we can spend a little more memory and CPU to track in smaller time windows fast enough and accurately. This is a trade-off between memory and being real time (less near realtime - if that make sense chatGPT<-).
