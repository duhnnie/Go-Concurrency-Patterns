# Challenge: _Fan-out + Multiplexing — Concurrent Consumers on a Merged Stream_

## Concept 

In a high-throughput system (e.g., a message broker, logging pipeline, or Kafka-like consumer group), multiple workers may need to process messages coming from a **single merged stream** concurrently — _without blocking each other_.

You’ll extend your **MessageMerger** so that multiple consumer goroutines can safely process messages from the merged output concurrently, while ensuring:

- Each message is processed **exactly once**.
- Slow consumers don’t block faster ones.
- The system shuts down gracefully when all work is done.

---
## Challenge description

#### Goal

Implement a **Fan-out Worker Pool** that consumes from the merged channel produced by your `MessageMerger`.

Each worker:
- Reads from the shared merged output.
- Processes messages independently (simulated by random sleep).
- Logs its worker ID and the message it processed.

The system:

- Should start _N_ worker goroutines.
- Should process all messages until the merger closes.
- Should use **context.Context** for graceful cancellation and shutdown.
- Should wait for all workers to finish before exiting.

---
## Requirements

1. **MessageMerger** — Use your refactored version (from previous challenge).
    
2.  **Worker Pool:**
    - A configurable number of workers (e.g., 3–5).
    - Each reads from the same channel (`merger.Output`).
    - Each processes independently (no blocking others).
        
3. **Processing Simulation:**
    - Random `time.Sleep()` to simulate variable workloads.
        
4. **Graceful Shutdown:**
    - When the merger closes its output channel, all workers should exit cleanly.
    - Main goroutine waits for all workers to finish (`sync.WaitGroup`).
        
5. **Logging:**
    - Each processed message should print something like:
    ```
    [Worker #2] processed: [tech] tech message #3
    ```
---

### Optional Stretch Goals

If you want to make it more interesting later:
- Add retry logic if a worker “fails” (simulate with random errors).
- Add metrics (e.g., messages processed per worker).
- Make the worker pool dynamic (add/remove workers at runtime).

---
## Hint diagram 

```
     +-------------+         +-------------------+
     |  sportsChan |--\      |                   |
     |   techChan  |---\     |                   |
     |   newsChan  |---->--> |   MessageMerger   | ---> merged output channel
     |   gameChan  |---/     |                   |
     +-------------+         +-------------------+
                                     |
                                     |
                              +--------------+
                              |  Worker #1   |  <-- reads from same merged channel
                              +--------------+
                              +--------------+
                              |  Worker #2   |
                              +--------------+
                              +--------------+
                              |  Worker #3   |
                              +--------------+

Each worker reads from the same merged output channel concurrently.

```