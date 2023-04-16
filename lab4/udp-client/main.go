package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// A UDP client that dials into localhost:12031 and sends a message to the server
func main() {
	var hostname string
	var port int

	hostnamePtr := flag.String("h", "localhost", "eg. localhost")
	portPtr := flag.Int("p", 12031, "eg. 12031")
	flag.Parse()
	hostname = *hostnamePtr
	port = *portPtr

	numThreads := 10_000
	maxConcurrentConn := 100                               // Number of concurrent connections to the server. Prevent DDOS
	connPool := make(chan *net.UDPConn, maxConcurrentConn) // Buffered channel for connection pool
	ctx := context.Background()
	var counter int64
	var wg sync.WaitGroup

	// Initialize connection pool
	// Since UDP is connectionless, just share a small pool of connections across all client goroutines
	// Reuse 100 static UDP connections
	for i := 0; i < maxConcurrentConn; i++ {
		conn, err := connect(hostname, port)
		if err != nil {
			log.Fatalf("error creating UDP connection: %s", err)
		}
		connPool <- conn
	}

	startTime := time.Now()

	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			conn := <-connPool // Acquire a conn from the connPool
			defer func() {
				connPool <- conn // Release the conn into the connPool
			}()
			err := runJob(ctx, conn, &counter)
			if err != nil {
				log.Printf("%s", err.Error())
			}
		}()
	}
	wg.Wait()

	// Close the connPool channel and close all connections
	close(connPool)
	for conn := range connPool {
		conn.Close()
	}

	duration := time.Since(startTime)
	fmt.Printf("Time taken for %d threads: %v\n", numThreads, duration)
}

func runJob(ctx context.Context, conn *net.UDPConn, counter *int64) error {
	// Write to the connection
	clientID := atomic.AddInt64(counter, 1)
	message := fmt.Sprintf("Client ID is %d\n", clientID)
	_, err := conn.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("error writing: %s", err)
	}

	// Wait and read from the connection
	buf := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		return fmt.Errorf("error reading: %s", err)
	}
	fmt.Printf("%s", string(buf[:n]))
	return nil
}

func connect(hostname string, port int) (*net.UDPConn, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", hostname, port))
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
