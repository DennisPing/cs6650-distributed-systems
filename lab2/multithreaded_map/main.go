package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	startTime := time.Now()
	rand.Seed(time.Now().UnixNano())
	var intMap sync.Map

	numThreads := 1000
	var wg sync.WaitGroup
	wg.Add(numThreads)

	for i := 0; i < numThreads; i++ {
		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ { // Do 1000 writes each
				randInt := rand.Intn(100) + 1 // Random number between 1 and 100
				value, ok := intMap.Load(randInt)
				if ok {
					intMap.Store(randInt, value.(int)+1)
				} else {
					intMap.Store(randInt, 1)
				}
			}
		}()
	}
	wg.Wait()
	duration := time.Since(startTime)
	fmt.Printf("Total duration: %v\n", duration)
	// There should be a total of 1 million writes to the map
}
