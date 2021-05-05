package main

import (
	"io"
	"log"
	"net"
)

// Add Host header so the remote server knows what site you want to load:
// curl -i -X GET -H 'Host: www.google.com' http://localhost

func handle(src net.Conn) {
	dst, err := net.Dial("tcp", "www.google.com:80")
	if err != nil {
		log.Fatalln("Unable to connect to our unreachable host")
	}
	defer dst.Close()

	go func() {
		if _, err := io.Copy(dst, src); err != nil {
			log.Fatalln(err)
		}
	}()

	if _, err := io.Copy(src, dst); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	listener, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Fatalln("Unable to bind to port")
	}

	log.Println("Listening on 0.0.0.0:80")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("Unable to accept connection")
		}

		go handle(conn)
	}
}
