package main

import (
	"fmt"
	"math/rand"
	"time"
)

func worker(name string, out chan<- string) {
	for {
		// Simulate variable work
		time.Sleep(time.Duration(rand.Intn(800)+500) * time.Millisecond)
		out <- fmt.Sprintf("%s finished work", name)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	worker1 := make(chan string)
	worker2 := make(chan string)

	go worker("Worker1", worker1)
	go worker("Worker2", worker2)

	// Listen to both workers using select
	for i := 0; i < 10; i++ {
		select {
		case msg1 := <-worker1:
			fmt.Println("Received from worker1:", msg1)
		case msg2 := <-worker2:
			fmt.Println("Received from worker2:", msg2)
		case <-time.After(1 * time.Second):
			fmt.Println("Timeout! No workers responded.")
		}
	}

	fmt.Println("Done listening.")
}
