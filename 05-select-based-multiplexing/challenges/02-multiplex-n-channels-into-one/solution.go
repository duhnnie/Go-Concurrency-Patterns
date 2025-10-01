package main

import (
	"fmt"
	"sync"
	"time"
)

func merge(inputs ...<-chan string) <-chan string {
	output := make(chan string, len(inputs))
	wg := &sync.WaitGroup{}

	wg.Add(len(inputs))

	for _, input := range inputs {
		go func(inputChannel <-chan string) {
			for value := range inputChannel {
				output <- value
			}

			wg.Done()
		}(input)
	}

	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)
	ch3 := make(chan string)

	go func() {
		defer close(ch1)
		for i := 0; i < 3; i++ {
			ch1 <- fmt.Sprintf("ch1: %d", i)
			time.Sleep(200 * time.Millisecond)
		}
	}()

	go func() {
		defer close(ch2)
		for i := 0; i < 3; i++ {
			ch2 <- fmt.Sprintf("ch2: %d", i)
			time.Sleep(400 * time.Millisecond)
		}
	}()

	go func() {
		defer close(ch3)
		for i := 0; i < 3; i++ {
			ch3 <- fmt.Sprintf("ch3: %d", i)
			time.Sleep(600 * time.Millisecond)
		}
	}()

	out := merge(ch1, ch2, ch3)

	for v := range out {
		fmt.Println(v)
	}

	fmt.Println("Done")
}
