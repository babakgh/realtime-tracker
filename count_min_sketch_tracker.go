package main

import (
	"hash/fnv"
	"math"
	"sync"
)

type CountMinSketch struct {
	mu    sync.Mutex
	width int
	depth int
	count [][]int
	hash  []func(string) int
}

func NewCountMinSketch(width, depth int) *CountMinSketch {
	hashFunctions := make([]func(string) int, depth)
	for i := 0; i < depth; i++ {
		seed := uint32(i)
		hashFunctions[i] = func(data string) int {
			h := fnv.New32a()
			h.Write([]byte(data))
			return int(h.Sum32() ^ seed)
		}
	}

	count := make([][]int, depth)
	for i := range count {
		count[i] = make([]int, width)
	}

	return &CountMinSketch{
		width: width,
		depth: depth,
		count: count,
		hash:  hashFunctions,
	}
}

func (cms *CountMinSketch) Update(key string, value int) {
	cms.mu.Lock()
	defer cms.mu.Unlock()

	for i, hash := range cms.hash {
		cms.count[i][hash(key)%cms.width] += value
	}
}

func (cms *CountMinSketch) Estimate(key string) int {
	cms.mu.Lock()
	defer cms.mu.Unlock()

	min := math.MaxInt32
	for i, hash := range cms.hash {
		count := cms.count[i][hash(key)%cms.width]
		if count < min {
			min = count
		}
	}
	return min
}
