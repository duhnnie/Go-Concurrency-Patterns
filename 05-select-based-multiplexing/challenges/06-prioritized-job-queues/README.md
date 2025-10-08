# Challenge: Prioritized Job Queues with Select-Based Multiplexing

## Goal

Extend your current **load-balanced worker pool** so it supports **multiple job queues** with **different priorities** — for example, “high”, “medium”, and “low”.

Workers should **always pull from the highest-priority queue available**, but still continue processing lower-priority jobs when no high-priority jobs are waiting.

All this must be achieved **without central locking** — only using **channels**, **select statements**, and **context cancellation**for graceful shutdown.

---
## Requirements

1. **Multiple Queues**
    
    - Create **3 job channels**:  
        `highPriorityJobs`, `mediumPriorityJobs`, and `lowPriorityJobs`.
        
    - Each can receive independent jobs.
        
2. **Priority-Based Dispatch**
    
    - Your dispatcher must use multiplexing through `select` to:
        
        - Prefer high-priority jobs if available.
            
        - Otherwise fall back to medium, then low.
        
3. **Worker Reuse (Pull-based model)**
    
    - Keep the same “ready worker” pattern from your previous challenge, workers announce themselves when ready.
        
    - Dispatcher still assigns one job to one ready worker at a time.
        
4. **Graceful Shutdown**
    
    - When all job queues close and all workers finish, everything must exit cleanly.
        
5. **Dynamic Load Simulation**
    
    - Randomly feed jobs into all three priority queues with variable frequency to simulate different traffic patterns.
        

---
## Constraints

|Constraint|Description|
|---|---|
|**Channels only**|No mutexes for job queue management.|
|**Context-based cancellation**|Dispatcher and workers must stop gracefully when main cancels.|
|**No busy waiting**|Avoid `for { select { default: ... }}` loops that spin without purpose.|
|**Prioritized dispatch**|If high-priority jobs exist, lower-priority jobs must wait.|
|**No lost jobs**|Every enqueued job must be processed exactly once.|
|**Random delays**|Each worker should sleep randomly (like before) to simulate variable load.|

---

## Suggested Starting Structure

```go
type Dispatcher struct {     
	highPriorityJobs   <-chan string     
	mediumPriorityJobs <-chan string     
	lowPriorityJobs    <-chan string     
	readyWorkersChan   <-chan *Worker     
	ctx                context.Context     
	cancel             context.CancelFunc 
}
```

---

## Extra Credit (Optional)

Add a **“fairness mode” toggle** where:

- When enabled, after each job from a higher-priority queue, the next dispatch cycle temporarily prefers a lower-priority queue (to prevent starvation).
    

---

## Success Criteria

- Output clearly shows workers receiving jobs of **different priorities**.
    
- Higher-priority jobs appear **processed earlier**.
    
- The program **terminates cleanly** after all queues close.
    
- No goroutines leak.
    