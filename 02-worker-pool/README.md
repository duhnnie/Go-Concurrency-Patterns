# Worker Pool

The worker pool is a concurrency pattern that:

- Defines a specific number of workers to run concurrently and perform a task.
- This workers consume task to perform from a shared channel among workers.
- Task source (or producer) sends tasks through the channel so workers can pull them from it.
- Once workers are done with all tasks they are shutdown gracefully.

It is useful when you have lots of tasks to perform but you don't want to spawn an unlimited number of goroutines, so you control/limit concurrency.

## Example

```go
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Task type
type Task struct {
	id int
}

// Worker logic
func worker(id int, tasks <-chan Task, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range tasks {
		// Simulate work
		time.Sleep(time.Duration(rand.Intn(300)) * time.Millisecond)
		fmt.Printf("Worker %d processed task %d\n", id, task.id)
	}
}

func main() {
	const numWorkers = 3
	const numTasks = 10

	tasks := make(chan Task, numTasks)
	wg := &sync.WaitGroup{}

	// Start workers
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, tasks, wg)
	}

	// Send tasks
	for i := 1; i <= numTasks; i++ {
		tasks <- Task{id: i}
	}
	close(tasks)

	// Wait for workers to finish
	wg.Wait()
	fmt.Println("All tasks processed")
}

```

### Key Parts

1. **Tasks channel** → Shared job queue.
    
2. **Fixed worker goroutines** → They all read from the same channel.
    
3. **WaitGroup** → Ensures all workers finish.
    
4. **close(tasks)** → Signals “no more work” to workers.
