package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type URLStatus struct {
	url    string
	status int
}

type ClassifyResult struct {
	url            string
	classification string
}

func produce(output chan<- string) {
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

	for _, url := range urls {
		output <- url
	}
	close(output)
}

func fetch(input <-chan string, output chan<- URLStatus, wg *sync.WaitGroup) {
	defer wg.Done()
	for url := range input {
		response, err := http.Get(url)
		if err != nil {
			panic(err)
		}

		output <- URLStatus{url, response.StatusCode}
		response.Body.Close()
	}
}

func classify(input <-chan URLStatus, output chan<- ClassifyResult, wg *sync.WaitGroup) {
	defer wg.Done()
	for urlStatus := range input {
		var classification string
		switch {
		case urlStatus.status >= 200 && urlStatus.status < 300:
			classification = "success"
		case urlStatus.status >= 300 && urlStatus.status < 400:
			classification = "redirect"
		case urlStatus.status >= 400:
			classification = "error"
		default:
			classification = "unknown"
		}

		output <- ClassifyResult{urlStatus.url, classification}
	}
}

func collect(input <-chan ClassifyResult, output chan<- map[string]string) {
	resultMap := make(map[string]string)
	for classification := range input {
		resultMap[classification.url] = classification.classification
	}
	output <- resultMap
	close(output)
}

const FETCHERS_NUM = 8
const CLASSIFIERS_NUM = 4

func main() {
	rand.Seed(time.Now().UnixNano())

	produceChan := make(chan string)
	fetchChan := make(chan URLStatus)
	classifyChan := make(chan ClassifyResult)
	collectChan := make(chan map[string]string)

	// Producer
	go produce(produceChan)

	// Fan-Out: fetchers
	var fetchWg sync.WaitGroup
	for i := 0; i < FETCHERS_NUM; i++ {
		fetchWg.Add(1)
		go fetch(produceChan, fetchChan, &fetchWg)
	}

	// Close fetchChan once all fetchers are done
	go func() {
		fetchWg.Wait()
		close(fetchChan)
	}()

	// Fan-Out: classifiers
	var classifyWg sync.WaitGroup
	for i := 0; i < CLASSIFIERS_NUM; i++ {
		classifyWg.Add(1)
		go classify(fetchChan, classifyChan, &classifyWg)
	}

	// Close classifyChan once all classifiers are done
	go func() {
		classifyWg.Wait()
		close(classifyChan)
	}()

	// Collector
	go collect(classifyChan, collectChan)

	// Fan-In: final result
	for url, classification := range <-collectChan {
		fmt.Printf("%s: %s\n", url, classification)
	}
}
