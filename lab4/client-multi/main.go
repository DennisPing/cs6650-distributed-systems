package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// A TCP client that dials into localhost:12031 and sends a message to the server
func main() {
	var hostname string
	var port int

	hostnamePtr := flag.String("h", "localhost", "eg. localhost")
	portPtr := flag.Int("p", 12031, "eg. 12031")
	flag.Parse()
	hostname = *hostnamePtr
	port = *portPtr

	numThreads := 1000
	maxConcurrentConn := 100                                // Number of concurrent connections to the server. Prevent DDOS
	connThrottler := make(chan struct{}, maxConcurrentConn) // Buffered channel for throttling outgoing connections
	ctx := context.Background()
	var counter int64
	var wg sync.WaitGroup

	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		connThrottler <- struct{}{} // Acquire a slot if there is available space, else block
		go func() {
			defer wg.Done()
			defer func() {
				<-connThrottler // Release the slot so that other connections can acquire slot
			}()
			err := runJob(ctx, hostname, port, &counter)
			if err != nil {
				log.Printf("%s", err.Error())
			}
		}()
	}
	wg.Wait()
}

func runJob(ctx context.Context, hostname string, port int, counter *int64) error {

	var conn *net.TCPConn
	var err error

	maxRetries := 5
	backOffDuration := 100 * time.Millisecond

	for i := 0; i < maxRetries; i++ {
		conn, err = connect(hostname, port)
		if err == nil {
			break // Connect successful
		}

		// Exponential backoff
		sleepDuration := time.Duration(rand.Int63n(int64(backOffDuration)))
		log.Printf("connection failed. Retrying in %s...", sleepDuration)
		time.Sleep(sleepDuration)
		backOffDuration *= 2
	}

	if err != nil {
		return fmt.Errorf("failed to connect after %d retries: %v", maxRetries, err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// Write to the connection
	clientID := atomic.AddInt64(counter, 1)
	_, err = writer.WriteString(fmt.Sprintf("Client ID is %d\n", clientID))
	if err != nil {
		return fmt.Errorf("error writing: %s", err)
	}
	writer.Flush()

	// Wait and read from the connection
	var resp string
	resp, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error reading: %s", err)
	}
	fmt.Printf("%s", resp)

	return nil
}

func connect(hostname string, port int) (*net.TCPConn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", hostname, port))
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
