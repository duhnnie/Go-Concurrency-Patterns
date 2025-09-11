# Pipeline

The Pipeline concurrency pattern consists in a sequence of stages that perform some work or transformation on the consumed values from an input channel and then sends that result through an output channel. 

The stages are connected through channels in a way in which the output channel for one stage is the input channel for the next one.

This allows you to process data in a _streaming fashion_ → as soon as one stage finishes processing an item, the next stage can start, even if previous items are still being worked on.

## Example

```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Stage 1: generate numbers
func generator(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

// Stage 2: square numbers
func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond) // simulate work
			out <- n * n
		}
		close(out)
	}()
	return out
}

// Stage 3: double numbers
func doubler(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * 2
		}
		close(out)
	}()
	return out
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Connect the stages
	numbers := generator(1, 2, 3, 4, 5)
	squares := square(numbers)
	doubles := doubler(squares)

	// Collect results
	for result := range doubles {
		fmt.Println(result)
	}
}

```

### Key Insights

1. Each stage is **independent** and only knows about its input/output channels.
    
2. Multiple goroutines can run per stage if you want parallelism.
    
3. Closing channels at the right time is crucial.