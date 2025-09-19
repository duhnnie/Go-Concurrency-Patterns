# Pub/Sub pattern

The **Publish–Subscribe (Pub/Sub)** pattern is a messaging pattern where:

- **Publishers** send messages to a central hub (called a **topic** or **broker**).
    
- **Subscribers** express interest in certain topics.
    
- The **broker** delivers published messages to all subscribers of that topic.
    

The key idea: **Publishers don’t know who the subscribers are**, and subscribers don’t know who the publishers are. This makes the system **loosely coupled** and highly extensible.

---
## Simple Real-Life Analogy

- Think of a **YouTube channel**:
    
    - The **YouTuber (Publisher)** uploads a new video (message).
        
    - The **Subscribers** to that channel automatically get notified.
        
    - The YouTuber doesn’t know (or care) exactly who the subscribers are.

---
## Example

```go
package main

import (
	"fmt"
	"time"
)

type Broker struct {
	subscribers []chan string
}

func NewBroker() *Broker {
	return &Broker{
		subscribers: []chan string{},
	}
}

func (b *Broker) Subscribe() <-chan string {
	ch := make(chan string)
	b.subscribers = append(b.subscribers, ch)
	return ch
}

func (b *Broker) Publish(msg string) {
	for _, sub := range b.subscribers {
		go func(c chan string) {
			c <- msg
		}(sub)
	}
}

func main() {
	broker := NewBroker()

	// Two subscribers
	sub1 := broker.Subscribe()
	sub2 := broker.Subscribe()

	// Publish some messages
	go func() {
		for i := 1; i <= 3; i++ {
			broker.Publish(fmt.Sprintf("Message %d", i))
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// Subscribers consume messages
	go func() {
		for msg := range sub1 {
			fmt.Println("Sub1 received:", msg)
		}
	}()
	go func() {
		for msg := range sub2 {
			fmt.Println("Sub2 received:", msg)
		}
	}()

	time.Sleep(2 * time.Second)
}
```

