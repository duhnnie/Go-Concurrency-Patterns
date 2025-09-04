# Code Challenge: Multi-Source log processor

## Scenario

Imagine you have multiple microservices in a system, and each one produces logs. You want to:

1. **Fan-Out:** Start multiple _producers_ (simulating microservices), each generating log entries concurrently.
    
2. **Fan-In:** Merge all log streams into a single channel.
    
3. **Fan-Out again:** Start multiple _processors_ that analyze the logs concurrently (e.g., count occurrences of “ERROR”).
    
4. **Fan-In again:** Collect the aggregated results into a single summary map:
    
    `{     "serviceA": 2 errors,     "serviceB": 5 errors,     "serviceC": 0 errors }`

## Requirements

1. **Producers**:
    
    - At least 3 producers (`serviceA`, `serviceB`, `serviceC`).
        
    - Each producer sends ~5 log lines like:
        
        - `"INFO: User logged in"`
            
        - `"ERROR: DB connection failed"`
            
        - `"WARN: Cache miss"`
            
    - Sleep randomly between sending logs to simulate real-time streaming.
        
2. **Fan-In (log merging)**:
    
    - Combine logs from all producers into a single channel.
        
3. **Processors (Fan-Out)**:
    
    - Start 3 workers that read logs and classify them (`INFO`, `WARN`, `ERROR`).
        
    - Focus only on counting `"ERROR"` logs per service.
        
4. **Final Fan-In**:
    
    - Collect results into a `map[string]int` where the key is the service name and the value is the error count.
        

---

## Expected Output Example

(since producers are concurrent, order will vary)

`serviceA: 2 errors serviceB: 3 errors serviceC: 1 errors`

---

👉 This challenge forces you to use **multiple levels of Fan-Out/Fan-In**:

- Stage 1: Multiple producers (fan-out logs).
    
- Stage 2: Merge logs (fan-in).
    
- Stage 3: Multiple processors (fan-out again).
    
- Stage 4: Aggregate results (fan-in again).