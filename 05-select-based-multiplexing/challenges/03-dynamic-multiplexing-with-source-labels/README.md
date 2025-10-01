# Challenge: Dynamic Multiplexing with Source Labels

## Problem Statement

You need to build a **dynamic fan-in multiplexer** in Go that:

1. Accepts **any number of input channels**.
    
2. Tags each message with a **source label** (string) instead of just an index.
    
3. Supports **dynamic addition and removal of input channels** at runtime.
    
4. Ensures that **slow channels do not block faster channels**.
    
5. Closes the output channel gracefully once all active input channels have been closed.
    

---
## **Example Behavior**

- You have channels `sportsChan`, `techChan`, `newsChan`.
    
- You merge them into one output channel.
    
- Each message is a struct:
    

```go
type TaggedMessage struct {     Source string     Value  string }
```

- Messages from `sportsChan` might look like:
    

```go
{Source: "sports", Value: "Team A won"}
```

- Messages from `techChan` might look like:
    

```go
{Source: "tech", Value: "New AI chip released"}
```

- The multiplexer should allow adding a new channel dynamically:
    

```go
gamingChan := make(chan string) AddChannel("gaming", gamingChan)
```

- And support removing channels dynamically:
    

```go
RemoveChannel("tech")
```

---

### **Requirements**

1. Output channel should **never be blocked by a slow sender**.
    
2. You may use **goroutines, sync.WaitGroup, and channels**.
    
3. Provide a simple **main function** with 3–4 channels sending messages at different rates.
    
4. Demonstrate **adding a new channel** after some messages, and **removing a channel** dynamically.
    
5. Print all messages from the merged output channel in the format:
    

```go
[Source] Message
```

---
This challenge combines:

- **N-to-1 multiplexing**
    
- **Tagged messages with meaningful source labels**
    
- **Dynamic channel management** (add/remove at runtime)
    
- **Non-blocking behavior**
