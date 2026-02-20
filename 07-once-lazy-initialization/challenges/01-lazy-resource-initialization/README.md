## Lazy Resource Initialization

### Goal

Implement a lazily initialized resource that is safe under concurrency.

### Scenario

You have an expensive resource that:

- takes time to initialize
    
- should only be initialized once
    
- is accessed by many goroutines
    

### Requirements

1. Use `sync.Once`
    
2. Initialization must print exactly once
    
3. All goroutines must receive the same instance
    
4. No global mutexes allowed outside `sync.Once`

### Skeleton
```go
type Resource struct {
	ID int
}

func GetResource() *Resource {
	// TODO
}
```

### Test Behavior

- Spawn 10 goroutines
    
- Each goroutine calls `GetResource`
    
- Output should show:
    
    - exactly one initialization message
        
    - consistent resource ID across all goroutines