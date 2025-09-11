package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type FetchResult struct {
	url    string
	status int
}

var validStatuses []int = []int{
	http.StatusOK,
	http.StatusNoContent,
	http.StatusAccepted,
	http.StatusNotFound,
	http.StatusForbidden,
	http.StatusUnauthorized,
	http.StatusInternalServerError,
	http.StatusNotImplemented,
	http.StatusBadGateway,
}

func worker(id int, input <-chan string, output chan<- FetchResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for url := range input {
		statusIndex := rand.Intn(len(validStatuses))
		status := validStatuses[statusIndex]

		// Simulate time-consuming job
		time.Sleep(time.Duration(rand.Int63n(300)) * time.Millisecond)
		fmt.Printf("Worker %d fetched %s with status %d\n", id, url, status)

		output <- FetchResult{
			url,
			status,
		}
	}
}

func main() {
	jobs := []string{
		"https://whatever.com",
		"https://la-razon.com",
		"https://macanas.com",
		"https://youtube.com",
		"https://google.com",
		"https://dinosaurs.com",
		"https://nin.com",
		"https://playstation.com",
		"https://sony.com",
		"https://nose.com",
		"https://golang.com",
		"https://yeyeye.com",
		"https://mercadolibre.com",
		"https://gatos.com",
	}

	const NUMBER_OF_WORKERS = 3

	jobsChannel := make(chan string, NUMBER_OF_WORKERS)
	resultsChannel := make(chan FetchResult)
	wg := &sync.WaitGroup{}

	for i := 0; i < NUMBER_OF_WORKERS; i++ {
		wg.Add(1)
		go worker(i+1, jobsChannel, resultsChannel, wg)
	}

	go func() {
		for _, job := range jobs {
			jobsChannel <- job
		}

		close(jobsChannel)
	}()

	go func() {
		wg.Wait()
		close(resultsChannel)
	}()

	results := make(map[string]int)

	for result := range resultsChannel {
		results[result.url] = result.status
	}

	fmt.Println("Final results:")
	fmt.Print(results)
}
