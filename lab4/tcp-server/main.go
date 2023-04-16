package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
)

// A TCP server that listens on port 12031 and spawns a new goroutine for each new connection
func main() {
	addr, err := net.ResolveTCPAddr("tcp", ":12031")
	if err != nil {
		log.Fatal(err)
	}

	listen, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Server started on %s", listen.Addr())
	defer listen.Close()

	numThreads := 20
	workerPool := make(chan struct{}, numThreads) // Buffered channel for limiting workers
	ctx := context.Background()
	var counter int64
	var wg sync.WaitGroup

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
		}

		workerPool <- struct{}{} // Acquire worker if available, else block
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer atomic.AddInt64(&counter, -1) // Decrement by 1 when the goroutine exits
			defer func() {
				<-workerPool // Release the worker
			}()
			handleConnection(ctx, conn, &counter)
		}()
	}
}

func handleConnection(ctx context.Context, conn net.Conn, counter *int64) {
	defer conn.Close()

	numGoroutines := atomic.AddInt64(counter, 1)
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		// Read from the TCP conn until the newline delimiter
		var message string
		message, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				return // Client disconnected
			}
			log.Printf("error reading: %s", err)
			return
		}
		fmt.Printf("%s", message)

		// Write response to the connection
		_, err = writer.WriteString(fmt.Sprintf("Number of goroutines on server: %d\n", numGoroutines))
		if err != nil {
			log.Printf("error writing: %s", err)
			return
		}
		writer.Flush()
	}
}
