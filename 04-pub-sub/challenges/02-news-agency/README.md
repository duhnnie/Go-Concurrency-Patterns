# Challenge: News Agency with Multiple Reporters and Subscribers

## Scenario

You’re building a **news agency system**.

- Multiple **reporters** (goroutines) publish news to different topics (`sports`, `tech`, `politics`).
    
- Multiple **subscribers** (goroutines) listen to topics they care about.
    
- Some subscribers listen to **all news** (like an archiver).
    
- Subscribers should be able to **unsubscribe** after a while.

## Requirements

1. Implement reporters that periodically publish random news headlines.
    
2. Subscribers should print received news.
    
3. After some time, unsubscribe one of the subscribers (e.g., the sports subscriber).
    
4. Keep the system running for a few seconds to demonstrate the Pub/Sub flow.

## Expected Behavior Example (sample run)

```
Sports: Team A won the match!
Tech: New AI chip released by major company
Politics: Election results announced
All: Team A won the match!
All: New AI chip released by major company
All: Election results announced
Sports subscriber is unsubscribing...
Tech: Quantum breakthrough in labs
All: Quantum breakthrough in labs
```

