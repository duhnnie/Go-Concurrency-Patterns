# Challenge: Concurrent Transfers Simulation

## Goal

Simulate multiple concurrent **money transfers** between several bank accounts while keeping all balances consistent.

Your program should:

- Create **N accounts**, each with an **initial balance**.
    
- Run **M concurrent transfer operations**, where each transfer:
    
    - randomly selects **two different accounts** (source and target),
        
    - randomly selects an **amount**,
        
    - transfers it **only if the source has enough balance**.
        

After all transfers complete, the **total balance across all accounts must remain the same** as at the start.

---

## Constraints

- Use a struct `Account` with fields like `id`, `balance`, and a `sync.Mutex`.
    
- A global slice of accounts represents the shared state.
    
- Protect all balance reads/writes with proper locking.
    
- Prevent **deadlocks** when locking two accounts for a transfer.
    
    - (Hint: always lock in a consistent order, e.g. by account ID.)
        
- Randomize transfer amounts and delays to simulate load.
    
- Print all transfers and a final balance summary.
    

---

## Example Output

```
Transfer 12.00 from #2 → #5 (OK)
Transfer 8.00 from #4 → #1 (Insufficient funds)
Transfer 6.00 from #3 → #2 (OK)
...
Final total balance: 1000.00 (initial total: 1000.00)
```

---

## Goal Validation

✅ All account operations are thread-safe.  
✅ No deadlocks occur.  
✅ Total balance stays constant.  
✅ Output shows mixed success and failure due to insufficient funds.