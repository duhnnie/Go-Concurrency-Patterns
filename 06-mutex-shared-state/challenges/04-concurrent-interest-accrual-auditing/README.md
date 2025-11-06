# Challenge: Concurrent Interest Accrual and Auditing

## Goal

Simulate a bank system where:

- Multiple **accounts** accrue periodic interest concurrently.
    
- Several **auditors** simultaneously inspect (read) account balances to verify total consistency.
    
- Access to shared state (balances) must remain **thread-safe** using mutexes.
    
- The **total interest accrued** and the **total balance consistency** must be correctly tracked.
    

---

## Description

Each account:

- Has an initial balance and accrues interest every fixed interval (e.g. every 200–400 ms).
    
- Interest is a small random percentage (e.g. 0.1% – 1%) applied to its balance.
    
- Can be safely read or written at any time by auditors or other goroutines.
    

Auditors:

- Periodically compute the total sum of **all balances** and **log it**, ensuring no read/write conflict occurs.
    
- Run concurrently with interest accrual goroutines.
    

---

## Constraints

1. You must use **sync.RWMutex** to protect shared state:
    
    - `RLock/RUnlock` for auditors (readers).
        
    - `Lock/Unlock` for accrual updates (writers).
        
2. Support **N accounts** (e.g. 10–20) and **M auditors** (e.g. 3–5).
    
3. The program should run for a **limited simulation time** (e.g. 5 seconds), then stop gracefully.
    
4. All goroutines must exit cleanly (use `sync.WaitGroup` and `context.Context` for cancellation).
    
5. At the end, print:
    
    - Total number of interest updates performed.
        
    - The final total balance across all accounts.
        

---

## Bonus (Optional)

If you want to challenge yourself further:

- Track and print **the average interest rate applied** during the simulation.
    
- Detect and log **inconsistencies** (e.g. total balance unexpectedly decreases).
    
- Simulate **a random account freeze**: one account periodically becomes read-only.

---
## Skeleton Starter

```go
package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Account struct {
	id      int
	balance float64
	mu      *sync.RWMutex
}

func (a *Account) AccrueInterest(rate float64) {
	// TODO: safely apply interest
}

func (a *Account) Balance() float64 {
	// TODO: safely read balance
	return 0
}

func accrueWorker(ctx context.Context, acc *Account, wg *sync.WaitGroup, updates *int64) {
	defer wg.Done()
	// TODO: periodically apply random interest
}

func auditorWorker(ctx context.Context, accounts []*Account, wg *sync.WaitGroup) {
	defer wg.Done()
	// TODO: periodically read all balances safely
}

func main() {
	// TODO: initialize accounts, start accrual and audit workers, and stop after timeout
}

```