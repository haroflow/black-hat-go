package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

type FooReader struct{}

func (f *FooReader) Read(p []byte) (n int, err error) {
	fmt.Print("in> ")
	return os.Stdin.Read(p)
}

type FooWriter struct{}

func (f *FooWriter) Write(p []byte) (n int, err error) {
	fmt.Print("out> ")
	return os.Stdout.Write(p)
}

func main() {
	var (
		reader FooReader
		writer FooWriter
	)

	// Copying manually
	// input := make([]byte, 4096)
	// s, err := reader.Read(input)
	// if err != nil {
	// 	log.Fatalln("Unable to read data")
	// }
	// fmt.Printf("Read %d bytes from stdin\n", s)

	// s, err = writer.Write(input)
	// if err != nil {
	// 	log.Fatalln("Unable to write data")
	// }
	// fmt.Printf("Wrote %d bytes to stdout\n", s)

	// Convenience function
	if _, err := io.Copy(&writer, &reader); err != nil {
		log.Fatalln("Unable to read/write data")
	}
}
