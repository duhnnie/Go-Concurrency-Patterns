package main

import (
	"fmt"
	"strings"
	"unicode"
)

// Generator: Emit sentences into a channel
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

// Tokenizer: Split each sentence into words and emit them individually
func tokenizer(input <-chan string) <-chan string {
	output := make(chan string)

	go func() {
		for sentence := range input {
			words := strings.Fields(sentence)

			for _, word := range words {
				output <- word
			}
		}

		close(output)
	}()

	return output
}

// Normalizer: lowercase each word, remove punctuation
func normalizer(input <-chan string) <-chan string {
	output := make(chan string)

	go func() {
		for token := range input {
			wordRunes := []rune{}

			for _, r := range token {
				if !unicode.IsPunct(r) {
					wordRunes = append(wordRunes, r)
				}
			}

			output <- strings.ToLower(string(wordRunes))
		}

		close(output)
	}()

	return output
}

// Filter: only pass words with length greater or equal to 4
func filter(input <-chan string) <-chan string {
	output := make(chan string)

	go func() {
		for word := range input {
			if len([]rune(word)) >= 4 {
				output <- word
			}
		}

		close(output)
	}()

	return output
}

// Counter: count word occurrences
func counter(input <-chan string) <-chan map[string]int {
	output := make(chan map[string]int)

	go func() {
		resultMap := make(map[string]int)

		for word := range input {
			resultMap[word]++
		}

		output <- resultMap
		close(output)
	}()

	return output
}

func main() {
	sentences := []string{
		"Go is expressive, concise, clean, and efficient.",
		"Concurrency is not parallelism.",
		"Channels orchestrate communication.",
	}

	generatorChan := generator(sentences)
	tokenizerChan := tokenizer(generatorChan)
	normalizerChan := normalizer(tokenizerChan)
	filterChan := filter(normalizerChan)
	counterChan := counter(filterChan)

	fmt.Println(<-counterChan)
}
