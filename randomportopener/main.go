// randomportopener opens a number of random ports on your system for testing.
// It will accept connections and close immediately.
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
)

var (
	listenQt   = flag.Int("qt", 10, "how many random ports to listen")
	listenFrom = flag.Int("from", 1, "from port")
	listenTo   = flag.Int("to", 1024, "to port")
	listenAddr = flag.String("addr", "127.0.0.1", "")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Opens random ports on the specified address.\n\n")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	fmt.Println("Use -h for help.")
	fmt.Println("If no ports could be opened, run this as root.")
	fmt.Println()

	if *listenTo < 1 || *listenTo > 65535 {
		fmt.Println("-to out of the range 1-65535")
		os.Exit(1)
	}
	if *listenFrom < 1 || *listenFrom > 65535 {
		fmt.Println("-from out of the range 1-65535")
		os.Exit(1)
	}
	if *listenFrom > *listenTo {
		fmt.Println("-from should be less than or equal to -to")
		os.Exit(1)
	}
	if *listenQt < 1 {
		fmt.Println("-qt should be greater than zero")
		os.Exit(1)
	}
	// Clamp qt if greater than the range of ports
	if *listenQt > *listenTo-*listenFrom {
		*listenQt = *listenTo - *listenFrom + 1
	}

	rand.Seed(time.Now().UnixNano())

	usedPorts := make([]int, 0, *listenQt)

mainLoop:
	for len(usedPorts) < cap(usedPorts) {
		port := *listenFrom + rand.Intn(*listenTo-*listenFrom+1)

		for _, p := range usedPorts {
			if p == port {
				continue mainLoop
			}
		}

		usedPorts = append(usedPorts, port)
		go openPort(port)
	}

	select {}
}

func openPort(port int) {
	address := fmt.Sprintf("%s:%d", *listenAddr, port)
	fmt.Println("Opening", address)

	ln, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		fmt.Printf("Accepted client on: %s, remote addr: %s. Closing connection.\n", address, conn.RemoteAddr())
		conn.Close()
	}
}
