# Code Challenge: Multi-Stage Pipeline
### **Concept**

- A **pipeline** is a sequence of stages where each stage runs concurrently.
    
- Each stage:
    
    - Receives input from the previous stage (via a channel).
        
    - Processes the data.
        
    - Sends results to the next stage (via another channel).
        
- Inside each stage you can use Fan-Out (multiple workers) and Fan-In (collect their results).
    

---

### **Challenge: URL Status Analyzer Pipeline**

**Stages:**

1. **Producer Stage**
    
    - Send URLs into the pipeline.
        
2. **Fetcher Stage (Fan-Out)**
    
    - Workers simulate fetching URLs (like you did).
        
    - Output: `(url, status)` pairs.
        
3. **Classifier Stage (Fan-Out)**
    
    - Workers classify each status code into `"success"` (200–299), `"redirect"` (300–399), or `"error"` (400+).
        
    - Output: `(url, category)` pairs.
        
4. **Collector Stage (Fan-In)**
    
    - Gather all classified results into a `map[string]string`.
        

---

### **Data Flow**

`URLs -> [Fetcher Workers] -> (url, status) -> [Classifier Workers] -> (url, category) -> Collector -> Final Map`

---

### **Your Task**

- Implement this pipeline.
    
- Reuse your `WorkResult` struct for stage 2.
    
- Create a new `ClassifyResult` struct for stage 3.
    
- Use **channels + WaitGroups** to handle concurrency.
    
- At the end, print the map of `URL -> category`.