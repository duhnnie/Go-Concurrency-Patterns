# Challenge: Merge two channels

### Problem

Implement a function:

```go
func merge(ch1, ch2 <-chan string) <-chan string
```

that:

- Starts goroutines to read from **both input channels**.
    
- Uses **`select`** to multiplex between them.
    
- Forwards messages into a single output channel.
    
- Stops and closes the output channel **only when both inputs are closed**.

## Example usage

```go
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
```

## Expected output

```go
ch1: 0
ch2: 0
ch1: 1
ch1: 2
ch2: 1
...
Done
```

## Expected behavior

- Output should **interleave values** from `ch1` and `ch2` depending on timing.
    
- Program should **exit cleanly** when both input channels are closed.

---
The challenge focuses only on **merging channels using `select`**, which is the simplest and most practical application of multiplexing.
