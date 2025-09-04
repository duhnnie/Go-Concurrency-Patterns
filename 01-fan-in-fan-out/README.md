# Fan-in / Fan-out

## Fan-out

**Concept:**

- You have a set of tasks to perform (e.g., processing jobs).
    
- Instead of doing them one by one, you start multiple **goroutines** to handle tasks concurrently.
    
- This is “fanning out” the work across multiple workers.
    

**Example:**  
Imagine you have 100 numbers, and you want to square each number concurrently. You could launch multiple goroutines, each handling a subset of numbers.

## Fan-in

**Concept:**

- After multiple goroutines produce results, you need to **combine** them into a single channel so you can process the output in one place.
    
- This is “fanning in” the results.
    

**Example:**

- Each goroutine sends its squared numbers to a channel.
    
- A single channel collects all these results, so your main function can process them together.

## Example

```
Numbers -> [Worker1, Worker2, Worker3] -> Results Channel -> Main
(Fan-Out)                       (Fan-In)
```

```go
package main

import (
    "fmt"
)

func worker(id int, jobs <-chan int, results chan<- int) {
    for j := range jobs {
        fmt.Printf("Worker %d processing job %d\n", id, j)
        results <- j * j // square the number
    }
}

func main() {
    jobs := make(chan int, 5)
    results := make(chan int, 5)

    // Fan-Out: start 3 workers
    for w := 1; w <= 3; w++ {
        go worker(w, jobs, results)
    }

    // Send 5 jobs
    for j := 1; j <= 5; j++ {
        jobs <- j
    }
    close(jobs)

    // Fan-In: collect results
    for a := 1; a <= 5; a++ {
        fmt.Println(<-results)
    }
}
```

**What happens here:**

1. `jobs` channel holds tasks.
    
2. Three workers read from `jobs` concurrently (**fan-out**).
    
3. Workers send results to `results` channel.
    
4. Main function collects all results (**fan-in**).

**Key Takeaways**

- Fan-Out increases concurrency, speeding up processing.
    
- Fan-In combines results safely without race conditions.
    
- Channels are the glue connecting workers and main routine.

