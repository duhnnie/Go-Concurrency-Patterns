# Once / Lazy initialization

## What problem this pattern solves

In concurrent programs, some resources must be initialized **exactly once**, even when accessed by many goroutines at the same time.

Typical examples include:

- configuration loading
- connection pools
- expensive lookup tables
- singleton services
- caches    

The key requirements are:

- **Single execution**
	Initialization logic must run once and only once.

- **Concurrency safety**
	Multiple goroutines may attempt initialization concurrently.

- **Visibility guarantees**
	Multiple goroutines may attempt initialization concurrently.

- **Blocking semantics**
	Goroutines that arrive during initialization must wait until it finishes.

---
## Why a Mutex Alone is Not Enough

A naïve approach might look like this:
```go
if resource == nil {
    resource = initResource()
}
```

Under concurrency, this fails because:

- multiple goroutines can pass the `nil` check
- initialization may run multiple times
- partially initialized state may be observed

Even with a mutex, the logic becomes error-prone and repetitive.

---
## `sync.Once`: the canonical solution

Go provides `sync.Once` to express this pattern explicitly.

#### Guarantees of `sync.Once`

- The function passed to `Do` is executed **at most once**
- All goroutines calling `Do` block until the function completes
- Memory visibility is guaranteed (no stale reads)

---
## Canonical Example: Lazy configuration loading

**Scenario**
- A configuration file is expensive to load
- Multiple goroutines may request configuration concurrently
- Configuration should be loaded only when first needed

**Implementation**

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type Config struct {
	AppName string
	Version string
}

var (
	config     *Config
	configOnce sync.Once
)

func loadConfig() {
	fmt.Println("Loading config...")
	time.Sleep(500 * time.Millisecond) // simulate expensive I/O

	config = &Config{
		AppName: "MyService",
		Version: "1.0.0",
	}
}

func GetConfig() *Config {
	configOnce.Do(loadConfig)
	return config
}

func main() {
	var wg sync.WaitGroup
	wg.Add(5)

	for i := 0; i < 5; i++ {
		go func(id int) {
			defer wg.Done()
			cfg := GetConfig()
			fmt.Printf("Goroutine %d got config: %s %s\n", id, cfg.AppName, cfg.Version)
		}(i)
	}

	wg.Wait()
}

```

**What This Demonstrates**

- `loadConfig` executes exactly once
- All goroutines block until initialization finishes
- No explicit locking is needed at call sites
- After initialization, reads are effectively lock-free

---
## Important Limitation of `sync.Once`

If the initialization function:

- **panics** → `Once` is considered “done”
- **fails silently** → no retry is possible

This makes error handling a central design concern, which we will explore in later challenges.

