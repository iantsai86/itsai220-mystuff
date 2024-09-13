package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	url := "http://localhost:8081/payload"
	numRequests := 100

	start := time.Now()

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println("Request failed:", err)
				return
			}
			resp.Body.Close()
		}()
	}

	wg.Wait()

	duration := time.Since(start)
	fmt.Printf("Sent %d requests in %v\n", numRequests, duration)
}
