package main

import (
	"fmt"
	"time"
)

func main() {
	startTime := time.Now()

	numElements := 100_000

	intSlice := make([]int, numElements)
	for i := 0; i < numElements; i++ {
		intSlice[i] = i
	}
	duration := time.Since(startTime)
	fmt.Printf("Total duration: %v\n", duration)
	fmt.Printf("Number of elements: %d\n", len(intSlice))
}
