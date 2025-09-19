# Challenge: Topic-based Pub/Sub system

**Challenge:**  
Implement a **topic-based Pub/Sub system** in Go:

1. Each subscriber can subscribe to **specific topics** (e.g., `"sports"`, `"news"`, `"tech"`).
    
2. Publishers send messages to a **specific topic**.
    
3. Only subscribers to that topic receive the message.
    
4. Add a **collector subscriber** that prints all messages from all topics (like a logger).

## Example

```go
broker := NewBroker()
sportsSub := broker.Subscribe("sports")
newsSub := broker.Subscribe("news")
logger := broker.Subscribe("all")

broker.Publish("sports", "Team A won the match!")
broker.Publish("news", "Breaking: Go 2.0 announced!")

// sportsSub → receives only sports messages
// newsSub   → receives only news messages
// logger    → receives everything
```