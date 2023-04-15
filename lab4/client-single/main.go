package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
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

	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", hostname, port))
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// Write to the connection
	_, err = writer.WriteString(fmt.Sprintf("Client ID is %d\n", 1))
	if err != nil {
		log.Fatalf("error writing: %s", err)
	}
	writer.Flush()

	// Wait and read from the connection
	var resp string
	resp, err = reader.ReadString('\n')
	if err != nil {
		log.Fatalf("error reading: %s", err)
	}
	fmt.Printf("%s", resp)
}
