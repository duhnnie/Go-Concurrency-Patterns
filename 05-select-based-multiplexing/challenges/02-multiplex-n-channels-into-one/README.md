# Challenges: Multiplex N Channels into One

## Problem

You are given **N input channels** of the same type (`<-chan string`).  
Your task is to implement a function:

```go
func merge(channels ...<-chan string) <-chan string
```

This function should **merge messages from all input channels into a single output channel**.

The merged output channel should:

1. Forward values from all input channels as they arrive.
    
2. Close **only when all input channels are closed**.
    
3. Avoid goroutine leaks (each channel should be consumed until closed).

---
## Example usage

```go
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
```

---
## Requirements

- Your `merge` must handle **any number of channels** (`N ≥ 1`).
    
- Must **close the output channel** exactly once, when **all inputs are closed**.
    
- Must **not lose messages**.
    
- Must **not leak goroutines** after completion.
    
---

Hint: You’ll probably need a **`sync.WaitGroup`** to wait for all readers to finish before closing the output channel.