// Various implementations of simple port scanners.
package main

import (
	"fmt"
	"net"
	"sort"
	"sync"
)

var host = "127.0.0.1"

func main() {
	fmt.Println("You can use randomportopener as root to open ports on your system:")
	fmt.Println("  go run randomportopener/main.go")
	fmt.Println()

	fmt.Println("Scan on port 80:")
	{
		_, err := net.Dial("tcp", host+":80")
		if err == nil {
			fmt.Println("Port 80 is open!")
		} else {
			fmt.Println("Port 80 is NOT open!")
		}
	}
	fmt.Println()

	fmt.Println("Scan on ports 1 to 1024:")
	{
		for i := 1; i <= 1024; i++ {
			address := fmt.Sprintf("%s:%d", host, i)

			conn, err := net.Dial("tcp", address)
			if err != nil {
				// port is closed or filtered
				continue
			}
			conn.Close()
			fmt.Printf("%d open\n", i)
		}
	}
	fmt.Println()

	fmt.Println(`Concurrent but "too fast":`)
	{
		var wg sync.WaitGroup
		for i := 1; i <= 1024; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				address := fmt.Sprintf("%s:%d", host, i)

				conn, err := net.Dial("tcp", address)
				if err != nil {
					// port is closed or filtered
					return
				}
				conn.Close()
				fmt.Printf("%d open\n", i)
			}(i)
		}
		wg.Wait()
	}
	fmt.Println()

	fmt.Println("100 workers:")
	{
		ports := make(chan int, 100)
		var wg sync.WaitGroup
		for i := 0; i < cap(ports); i++ {
			go worker(ports, &wg)
		}

		for i := 1; i <= 1024; i++ {
			wg.Add(1)
			ports <- i
		}

		wg.Wait()
		close(ports)
	}
	fmt.Println()

	fmt.Println("Multichannel communication:")
	{
		ports := make(chan int, 100)
		results := make(chan int)
		var openports []int

		for i := 0; i < cap(ports); i++ {
			go worker2(ports, results)
		}

		go func() {
			for i := 1; i <= 1024; i++ {
				ports <- i
			}
		}()

		for i := 0; i < 1024; i++ {
			port := <-results
			if port != 0 {
				openports = append(openports, port)
			}
		}

		close(ports)
		close(results)
		sort.Ints(openports)
		for _, port := range openports {
			fmt.Printf("%d open\n", port)
		}
	}
	fmt.Println()
}

func worker(ports chan int, wg *sync.WaitGroup) {
	for p := range ports {
		address := fmt.Sprintf("%s:%d", host, p)

		conn, err := net.Dial("tcp", address)
		if err != nil {
			// port is closed or filtered
			wg.Done()
			continue
		}
		conn.Close()
		fmt.Printf("%d open\n", p)
		wg.Done()
	}
}

func worker2(ports chan int, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("%s:%d", host, p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}
