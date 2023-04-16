package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
)

// A TCP server that listens on port 8080 and spawns a new goroutine for each new connection
func main() {
	addr, err := net.ResolveUDPAddr("udp", ":12031")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Server started on %s", conn.LocalAddr())
	defer conn.Close()

	numThreads := 20
	workerPool := make(chan struct{}, numThreads) // Buffered channel for limiting workers
	ctx := context.Background()
	var counter int64
	var wg sync.WaitGroup

	for {
		workerPool <- struct{}{} // Acquire worker if available, else block
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer atomic.AddInt64(&counter, -1) // Decrement by 1 when the goroutine exits
			defer func() {
				<-workerPool // Release the worker
			}()
			handlePacket(ctx, conn, &counter)
		}()
	}
}

func handlePacket(ctx context.Context, conn *net.UDPConn, counter *int64) {
	numGoroutines := atomic.AddInt64(counter, 1)

	buf := make([]byte, 1024)

	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("error reading: %s", err)
			return
		}
		message := string(buf[:n])
		fmt.Printf("%s", message)

		// Write the number of goroutines that are currently running
		resp := fmt.Sprintf("Number of goroutines on server: %d\n", numGoroutines)
		_, err = conn.WriteToUDP([]byte(resp), addr)
		if err != nil {
			log.Printf("error writing: %s", err)
			return
		}
	}
}
