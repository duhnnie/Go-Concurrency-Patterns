# Mutex / Shared State

## Concept Overview

In contrast to channel-based designs (which emphasize message passing), **mutex-based designs** focus on **protecting shared state** — ensuring that only one goroutine can access or modify it at a time.

At its core:

> “Don’t communicate by sharing memory; instead, share memory by communicating.”  
> — Go proverb

However, sometimes _shared memory_ is the practical choice — and Go’s `sync.Mutex` provides a simple, safe way to do it.

---

## What is a Mutex?

**Mutex (mutual exclusion)** is a synchronization primitive that ensures **exclusive access** to a resource.  
In Go, it’s provided by `sync.Mutex`.

It has two basic operations:

- `Lock()` — acquire ownership of the lock. If another goroutine holds it, wait.
    
- `Unlock()` — release ownership so others may proceed.
    

---

## Real-World Analogy

Think of a **bathroom key** at a café:

- Only one person can use it at a time.
    
- Others must wait until the person returns the key.
    

The shared resource (bathroom) is like your variable or map, and the key is the mutex.

---
## Example

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	var counter int
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 3; j++ {
				mu.Lock()
				counter++
				fmt.Printf("[goroutine %d] incremented counter → %d\n", id, counter)
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("Final counter:", counter)
}
```

---
### What’s Happening

|Step|Action|Explanation|
|---|---|---|
|1|`mu.Lock()`|One goroutine enters the _critical section_. Others block.|
|2|`counter++`|Safe — no concurrent write.|
|3|`mu.Unlock()`|Releases the lock; next waiting goroutine proceeds.|
|4|`wg.Wait()`|Wait for all goroutines to finish.|

If you remove the mutex, the program will have a **race condition**, producing inconsistent results.

---

## Common Pitfalls

- Forgetting to `Unlock()` (use `defer mu.Unlock()` inside critical sections).
    
- Locking too much (large critical sections = contention).
    
- Copying structs containing mutexes (never do that — they’re not copy-safe).