# Challenge: Dynamic Priority Escalation

## Goal

Extend your prioritized job dispatcher system to **dynamically adjust job priority levels** based on _waiting time_.  
In other words, jobs that wait too long in lower-priority queues should **automatically escalate** to a higher-priority queue, ensuring fairness and responsiveness under load.

This challenge builds on your previous “Prioritized Job Queues” dispatcher — but now, priority isn’t static.

---

## Background

In the real world, some low-priority tasks might become urgent if delayed for too long.  
Operating systems and load balancers often employ **priority aging** or **dynamic scheduling** strategies to prevent starvation.

Your system will simulate that:

- Jobs enter a queue with an initial priority (high, medium, or low).
    
- Over time, waiting jobs in lower-priority queues may be **promoted** to higher ones.
    

---

## Requirements

1. **Base Design:**
    
    - Keep the general structure from your previous challenge:
        
        - Three channels for high, medium, and low priority jobs.
            
        - A dispatcher that uses `select` to multiplex and assign available jobs to available workers.
            
        - A ready-worker pool mechanism (workers announce readiness).
            
2. **Dynamic Escalation Logic:**
    
    - Each job should include metadata tracking when it was created or enqueued.
        
    - If a job has been waiting longer than a configurable threshold:
        
        - Low → Medium
            
        - Medium → High
            
    - High-priority jobs remain high (no demotion).
        
    - Escalation should occur periodically or be checked before dispatching jobs.
        
3. **Configuration Parameters:**
    
    - `LOW_TO_MEDIUM_ESCALATION_MS` – how long before a low-priority job escalates to medium.
        
    - `MEDIUM_TO_HIGH_ESCALATION_MS` – how long before a medium-priority job escalates to high.
        
4. **Concurrency & Synchronization:**
    
    - Use a mutex or channels to coordinate safe movement of jobs between queues.
        
    - Escalation logic can run in a background goroutine scanning queues at intervals.
        
5. **Cancellation & Completion:**
    
    - Preserve the use of `context.Context` for graceful shutdown.
        
    - Dispatcher should terminate when all queues are empty and all workers are done.
        

---

## Constraints

- You must **still use select-based multiplexing** for worker dispatch.
    
- You **cannot** rely on external libraries (no heap-based priority queues).
    
- All priority queues must be **plain Go channels** (or slices guarded by mutexes).
    
- Your solution must **demonstrate escalation in action** (i.e., print logs when a job is promoted).
    
- Bonus (optional): support configurable escalation intervals via CLI args or constants.
    

---

## Success Criteria

- The program correctly escalates jobs over time.
    
- No job is lost or processed twice.
    
- The dispatcher never starves low-priority jobs indefinitely.
    
- Output logs clearly show promotions (e.g., `"[escalation] Job 🟢 #3 promoted to MEDIUM"`).
    
