package main

import (
	"fmt"
	"time"
)

func merge(c1, c2 <-chan string) <-chan string {
	out := make(chan string)

	go func() {
		defer close(out)

		ch1Open := true
		ch2Open := true

	loop:
		for {
			select {
			case v, ok := <-c1:
				if ok {
					out <- v
				} else {
					ch1Open = false

					if !ch2Open {
						break loop
					}
				}
			case v, ok := <-c2:
				if ok {
					out <- v
				} else {
					ch2Open = false

					if !ch1Open {
						break loop
					}
				}
			}
		}
	}()

	return out
}

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		defer close(ch1)
		for i := 0; i < 5; i++ {
			ch1 <- fmt.Sprintf("ch1: %d", i)
			time.Sleep(300 * time.Millisecond)
		}
	}()

	go func() {
		defer close(ch2)
		for i := 0; i < 5; i++ {
			ch2 <- fmt.Sprintf("ch2: %d", i)
			time.Sleep(500 * time.Millisecond)
		}
	}()

	out := merge(ch1, ch2)

	for v := range out {
		fmt.Println(v)
	}

	fmt.Println("Done")
}
