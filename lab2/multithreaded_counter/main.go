package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Tutorial: https://gobyexample.com/atomic-counters

func main() {
	startTime := time.Now()
	var counter uint64 = 0

	// Create a WaitGroup to track completion of all goroutines
	numThreads := 1000
	var wg sync.WaitGroup
	wg.Add(numThreads)

	// Start 1000 goroutines
	for i := 0; i < numThreads; i++ {
		go func() {
			defer wg.Done()
			for i := 0; i < 10; i++ {
				atomic.AddUint64(&counter, 1)
			}
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()

	duration := time.Since(startTime)
	fmt.Printf("Total duration: %v\n", duration)
	fmt.Printf("Final counter value: %d\n", counter)
}
