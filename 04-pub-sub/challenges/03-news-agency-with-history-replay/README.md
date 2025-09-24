# Pub/Sub Challenge: Fan-out with Message History Replay

You are tasked with enhancing your [**concurrent Pub/Sub system** (`NewsAgency`)](../02-news-agency/README.md) with the following features:

1. **Fan-out delivery (already present):**
    
    - When a publisher sends a message on a topic, the broker should deliver it to **all subscribers** of that topic and to the global subscribers (`SubscribeAll`).
        
    - Each subscriber must get their own copy of the message (fan-out).
        
2. **Message history replay (new feature):**
    
    - When a new subscriber subscribes to a topic, they should immediately receive the **last N messages** published for that topic.
        
    - After receiving the history, they continue to receive new incoming messages in real time.
        
    - Existing subscribers should _not_ replay old messages, only new ones.
        
    - The broker must support this for multiple topics and multiple subscribers per topic.
        
3. **Concurrency safety:**
    
    - Publishers and subscribers must work concurrently without race conditions.
        
    - Message history should be stored in a **bounded buffer (ring buffer)** to avoid unbounded memory growth.
        

---

👉 Example scenario:

- `Publisher` sends 10 news items on `SPORTS`.
    
- `Subscriber A` subscribes early → gets messages starting from #1 onward.
    
- `Subscriber B` subscribes later → should immediately replay the last 5 messages (#6–#10) and then continue with new ones.