package main

import (
	"io"
	"log"
	"net"
)

// telnet localhost 20080

func echo(conn net.Conn) {
	defer conn.Close()

	b := make([]byte, 512)
	for {
		size, err := conn.Read(b[0:])
		if err == io.EOF {
			log.Println("Client disconnected")
			break
		}
		if err != nil {
			log.Println("Unexpected error")
			break
		}

		// From the book example, it's printing a newline and
		// keeps printing final bytes from previous message:
		// log.Printf("Received %d bytes: %s\n", size, string(b))

		// Printing only bytes from 0 to size-1 gives the same
		// output as the book:
		log.Printf("Received %d bytes: %s\n", size, string(b[0:size-1]))

		log.Println("Writing data")
		if _, err := conn.Write(b[0:size]); err != nil {
			log.Fatalln("Unable to write data")
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":20080")
	if err != nil {
		log.Fatalln("Unable to bind to port")
	}

	log.Println("Listening on 0.0.0.0:20080")

	for {
		conn, err := listener.Accept()
		log.Println("Received connection")
		if err != nil {
			log.Fatalln("Unable to accept connection")
		}

		go echo(conn)
	}
}
