package main

import (
	"fmt"
	"strings"
	"sync"
	"unicode"
)

const COUNTER_NUM = 3

func generator(sentences []string) <-chan string {
	output := make(chan string)

	go func() {
		for _, sentence := range sentences {
			output <- sentence
		}

		close(output)
	}()

	return output
}

func tokenizer(sentencesChan <-chan string) <-chan string {
	output := make(chan string)

	go func() {
		for sentence := range sentencesChan {
			words := strings.Fields(sentence)

			for _, word := range words {
				output <- word
			}
		}

		close(output)
	}()

	return output
}

func normalizer(wordsChan <-chan string) <-chan string {
	output := make(chan string)

	go func() {
		for word := range wordsChan {
			wordWithoutPunctuation := []rune{}

			for _, r := range word {
				if !unicode.IsPunct(r) {
					wordWithoutPunctuation = append(wordWithoutPunctuation, r)
				}
			}

			output <- strings.ToLower(string(wordWithoutPunctuation))
		}

		close(output)
	}()

	return output
}

func filter(wordsChan <-chan string) <-chan string {
	output := make(chan string)

	go func() {
		for word := range wordsChan {
			if len([]rune(word)) >= 5 {
				output <- word
			}
		}

		close(output)
	}()

	return output
}

func counter(wordsChan <-chan string) <-chan map[string]int {
	output := make(chan map[string]int)

	go func() {
		resultsMap := make(map[string]int)

		for word := range wordsChan {
			resultsMap[word]++
		}

		output <- resultsMap
		close(output)
	}()

	return output
}

func counterDistribuitor(wordsChan <-chan string, counterNum int) []<-chan map[string]int {
	output := []<-chan map[string]int{}

	for i := 0; i < counterNum; i++ {
		ch := counter(wordsChan)
		output = append(output, ch)
	}

	return output
}

func counterMerger(resultsMapChanList ...<-chan map[string]int) <-chan map[string]int {
	output := make(chan map[string]int)
	wg := &sync.WaitGroup{}
	wg.Add(3)

	for _, resultsMapChan := range resultsMapChanList {

		go func(inputChan <-chan map[string]int) {
			for resultsMap := range inputChan {
				output <- resultsMap
			}

			wg.Done()
		}(resultsMapChan)
	}

	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}

func reducer(resultsMapChan <-chan map[string]int) <-chan map[string]int {
	output := make(chan map[string]int)

	go func() {
		finalResults := make(map[string]int)

		for resultsMap := range resultsMapChan {
			for k, v := range resultsMap {
				finalResults[k] += v
			}
		}

		output <- finalResults
		close(output)
	}()

	return output
}

func main() {
	sentences := []string{
		"Go is expressive, concise, clean, and efficient.",
		"Concurrency is not parallelism.",
		"Channels orchestrate communication.",
		"Parallel pipelines can improve performance.",
		"Goroutines are lightweight threads.",
	}

	sentencesChan := generator(sentences)
	tokeninzerChan := tokenizer(sentencesChan)
	normalizerChan := normalizer(tokeninzerChan)
	filter := filter(normalizerChan)
	counterChanList := counterDistribuitor(filter, COUNTER_NUM)
	counterChan := counterMerger(counterChanList...)

	fmt.Println(<-reducer(counterChan))
}
