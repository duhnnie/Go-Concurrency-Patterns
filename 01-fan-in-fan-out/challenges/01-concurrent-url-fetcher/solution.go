package main

import (
	"fmt"
	"net/http"
	"sync"
)

type WorkResult struct {
	url    string
	status int
}

func doWork(jobs <-chan string, results chan<- WorkResult, wg *sync.WaitGroup) {
	for url := range jobs {
		response, err := http.Get(url)
		if err != nil {
			panic(err)
		}

		results <- WorkResult{url, response.StatusCode}
		wg.Done()
	}
}

func main() {
	urls := []string{
		"https://httpbin.org/status/200",
		"https://httpbin.org/status/404",
		"https://httpbin.org/status/500",
		"https://httpbin.org/status/302",
		"https://httpbin.org/status/403",
		"https://httpbin.org/status/418",
		"https://httpbin.org/status/503",
		"https://httpbin.org/status/201",
	}

	jobs := make(chan string, len(urls))
	results := make(chan WorkResult, len(urls))
	wg := sync.WaitGroup{}
	resultMap := make(map[string]int)

	// Fan-Out: start fixed number of workers
	numWorkers := 8
	for w := 0; w < numWorkers; w++ {
		go doWork(jobs, results, &wg)
	}

	// Send jobs
	for _, url := range urls {
		wg.Add(1)
		jobs <- url
	}
	close(jobs)

	// Fan-In: wait for results
	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		resultMap[result.url] = result.status
	}

	fmt.Println(resultMap)
}
