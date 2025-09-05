# Code Challenge: Multi-level fan-in

In your previous solution:

- Services → generate logs (fan-out).
    
- Processors → classify logs (fan-out).
    
- Collector → aggregate into a single map (fan-in).
    

Now let’s extend it:

### 📌 Problem Statement

You need to **aggregate results in two stages**:

1. **Local Collectors**: Instead of a single collector, have _one collector per processor group_. Each local collector should produce a partial result map.
    
2. **Global Collector (Final Fan-In)**: All local results should then be merged into a final aggregated result.
    

### 🛠️ Requirements

- Keep **fan-out processors** as before.
    
- Each processor should send results to a **dedicated collector**.
    
- Collectors should send their partial results into a **global channel**.
    
- A final goroutine should perform a **global fan-in** of those partial results into a single final map.

### Expected Outcome

When the program ends, you should have one final `map[string]int` that shows the **total error logs per service**, just like before — but now produced via a **two-step fan-in** process.
