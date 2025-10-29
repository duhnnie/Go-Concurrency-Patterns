# Challenge: Concurrent Metrics Collector

## Goal

Implement a **thread-safe metrics collector** that multiple goroutines can safely update concurrently.  
Each goroutine will simulate a worker producing metrics — e.g., “requests processed,” “errors,” “latency” — and the collector will maintain global statistics using **mutex-protected shared state**.

At the end, the program should print the **aggregated results** without data races or inconsistent values.

---
## Requirements

1. **Define a `Metrics` struct** that holds:
    - `Requests int`
    - `Errors int`
    - `TotalLatency time.Duration`
    - `mu sync.Mutex` (protects all fields)
        
2. **Expose safe methods** on `Metrics`:
    - `RecordRequest(latency time.Duration)`  
        increments `Requests` and adds to `TotalLatency`
    - `RecordError()`  
        increments `Errors`
    - `AverageLatency() time.Duration`  
        computes the average latency safely (returns 0 if no requests)
        
3. **Simulate N workers** (e.g., 10 goroutines):  
    Each worker:
    - Processes 10–30 requests
    - Randomly sleeps between 100–500ms per request
    - Has a small chance (e.g., 15%) to record an error instead of success
    - Records its metrics using the shared `Metrics` instance.
        
4. **At the end**, after all workers finish:
    - Print total requests, errors, average latency, and success rate.
        
5. **No global variables.**  
    The shared metrics must be passed explicitly to workers.

---
## Constraints

- Must use **only `sync.Mutex`** (not channels or atomic primitives).
    
- Must be **race-free** (verify using `go run -race`).
    
- Use a **custom random source** (`rand.New(rand.NewSource(time.Now().UnixNano()))`).
    
- Must handle the case when **no requests** are recorded (avoid division by zero).
    

---
## Optional Enhancements (after main version works)

- Add a **Reset()** method to clear metrics atomically.
- Add per-worker IDs in the logs for traceability.
- Measure **total runtime** of the simulation.