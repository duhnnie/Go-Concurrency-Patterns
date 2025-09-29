# Challenge: Pub/Sub with Backpressure Handling

## Introduction

When a publisher is sending messages faster than subscribers can consume, we need strategies to avoid overwhelming the system. Instead of just dropping messages, we can apply techniques like:

1. **Block (backpressure)**
    
    - Broker **blocks the publisher** until the subscriber has room in its buffer (i.e., the send will block).
        
    - Guarantees no messages are lost, but a slow subscriber can slow down publishers.
        
2. **DropNewest (drop incoming)**
    
    - If the buffer is full, the broker **drops the incoming message** for that subscriber (does not block).
        
    - Fast publishers are never blocked, but messages can be lost for slow consumers.
        
3. **EvictOldest (drop oldest / keep newest)**
    
    - If the buffer is full, remove the **oldest queued** message from that subscriber’s buffer and insert the new one.
        
    - Keeps the most recent messages; useful when recent data is more valuable than old data.
        

(Extra/optional variants you might want later: **Block-with-timeout**, **Reject-with-error** to the publisher, or **persist-to-disk** fallback.)

## Description

Extend your Pub/Sub broker to handle **backpressure** gracefully.

- Add support for **per-subscriber strategies** for handling a full buffer.
    
- `Subscribe(topic, bufferSize int, strategy string)` — `strategy` ∈ `{ "block", "drop", "evict" }`.
    
- `Publish(topic, message string)` — for each subscriber of `topic`, deliver according to that subscriber’s strategy:
    
    - `"block"` → block until there’s room (synchronous backpressure).
        
    - `"drop"` → drop message for that subscriber (non-blocking).
        
    - `"evict"` → evict oldest queued message to make room for the new one.
        
- Demonstrate with 3 subscribers to the same topic, each using one of the strategies, and publish a burst of messages faster than consumers can read.

---
## Small notes about implementing each strategy

- **Block:** simply perform a blocking send into the subscriber’s buffered channel (or use a condition variable if you implement buffers yourself).
    
- **Drop:** use `select { case ch <- msg: default: /* drop */ }`.
    
- **EvictOldest:** needs a custom bounded queue (slice or ring) you manage per subscriber so you can pop the oldest and push the new one atomically (mutex-protected), or implement an internal goroutine per subscriber that accepts control requests to evict then enqueue.
---

## Example Output (conceptual)

```
[BlockingSub] received: msg-1 
[BufferSub] received: msg-1 
[EvictSub] received: msg-5   (older ones evicted)
```

---

## 🎯 Goal

This will simulate **real-world Pub/Sub systems** (like Kafka, NATS, or Redis Streams) where handling slow consumers is a critical part of design.

---