## Challenge: Select-Based Multiplexing — Load-Balanced Worker Pool

### **Goal**

Design a **load-balanced worker pool** where each worker **pulls jobs only when it’s available**, rather than competing for them.  
Use **`select` statements** for multiplexing between:
- job availability
- worker readiness
- graceful shutdown via `context.Context`
    

---

## **Conceptual Overview**

Traditional worker pools often do this:
- All workers read directly from a shared `jobs` channel.
- Whichever worker wins the race gets the job.
    

This works, but:
- It doesn’t account for **worker availability** or load differences.
- If one worker is slow, others might idle waiting for jobs to appear.
    

In this challenge, you’ll reverse the relationship:

> Workers **announce their readiness** to a dispatcher,  
> and the dispatcher assigns jobs to them **only when they’re ready**.

This is **pull-based load balancing** — coordinated through **select-based multiplexing**.

---

## **Requirements**

### **Entities**

You’ll implement three core entities:

1. **`Worker`**
    - Has a unique `id`
    - Has its own `jobChan` for receiving work
    - Runs in a goroutine
    - When idle, sends itself (or its channel) to a shared `ready` channel
        
2. **`Dispatcher`**
    
    - Listens on both `jobs` and `ready` channels
        
    - When a job and a worker are both available, it assigns the job
        
    - Uses a `select` loop to handle:
        - New jobs
        - Ready workers
        - Shutdown (`ctx.Done()`)
            
3. **`main()`**
    - Spawns a `Dispatcher` and several workers
    - Sends simulated jobs through a buffered channel
    - Cancels the context after some time (graceful shutdown)
        

---

## **Constraints & Realism Enhancements**

|Type|Constraint|
|---|---|
|**Job simulation**|Each job has a random “processing time” between 300–1000 ms|
|**Worker behavior**|Each worker prints when it starts and finishes a job|
|**Job queue**|Limit the `jobs` channel buffer size to 5 — simulating backpressure|
|**Availability delay**|Workers randomly delay before signaling readiness, simulating heterogeneous load|
|**Graceful shutdown**|Use `context.Context` cancellation to stop all goroutines|
|**Fair distribution**|Ensure every job gets processed by exactly one worker, no duplication or loss|

---

## **Expected Behavior**

- Jobs are processed _concurrently_ by multiple workers.
- Slower workers don’t block faster ones.
- The dispatcher cleanly stops when the context is cancelled.
- Output shows interleaved log lines like:
    
    ```
    worker-2 ready
	worker-2 received job 3
	worker-1 ready
	worker-1 received job 4
	worker-3 ready
	...
	worker-2 finished job 3
    ```
    

---
## 💡 **Your Task**

Implement in Go:

1. A `Worker` type that:
    - Sends itself (or its `jobChan`) to the `ready` channel when idle.
    - Processes jobs with simulated variable delay.
    - Stops gracefully on context cancellation.
        
2. A `Dispatcher` type that:
    - Uses a `select` loop to match jobs to ready workers.
    - Handles graceful shutdowns.
        
3. A `main()` function that:
    - Creates workers and a dispatcher.
    - Sends jobs periodically to the `jobs` channel.
    - Cancels the context after some time (e.g., 5 seconds).