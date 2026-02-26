# Rate Limiting Pattern
## What Is Rate Limiting?

**Rate Limiting** is a control mechanism that restricts how many requests or operations can be executed within a defined time window.

In backend systems (REST APIs, microservices, gateways), rate limiting is primarily used for:

- Protecting services from overload
    
- Preventing abuse (e.g., brute-force attacks)
    
- Enforcing fair usage between clients
    
- Implementing API quotas (e.g., 100 requests/minute)
    

From a distributed systems perspective, rate limiting is a **traffic shaping** and **back-pressure** mechanism.

---

## Common Rate Limiting Algorithms

Before coding, it’s important to understand the main algorithms:

1. **Fixed Window Counter**
    
2. **Sliding Window**
    
3. **Leaky Bucket**
    
4. **Token Bucket** (most common in Go systems)
    

In Go ecosystems, the **Token Bucket** model is dominant because it is simple, precise, and efficient.

---
## 3. Token Bucket — Concept

Think of it like this:

- A bucket holds **N tokens**
    
- Tokens are added at a constant rate (e.g., 10 tokens/sec)
    
- Each request consumes 1 token
    
- If no tokens are available → request is rejected or blocked
    

This allows:

- Short bursts
    
- Controlled sustained rate

---
## The Standard Way in Go

Go provides an official implementation in:

- golang.org/x/time (package `rate`)

This is the production-grade way to do rate limiting in Go.

---
## Basic Example: Token Bucket with `rate.Limiter`

### Scenario:

Allow:
- 2 requests per second
- Burst capacity of 5

### Code Example

```go
package main

import (
	"fmt"
	"time"

	"golang.org/x/time/rate"
)

func main() {
	// 2 events per second, burst size 5
	limiter := rate.NewLimiter(2, 5)

	for i := 0; i < 10; i++ {
		if limiter.Allow() {
			fmt.Println("Request", i, "allowed")
		} else {
			fmt.Println("Request", i, "blocked")
		}
		time.Sleep(200 * time.Millisecond)
	}
}
```

Explanation:

```go
rate.NewLimiter(2, 5)
```

- `2` -> tokens added per second
- `5` -> bucket capacity (burst size)

`limiter.Allow()`:
 if a token is available, it gets consumed and returns `true`. If no token is available it returns `false`.

## Blocking version

Instead of rejecting requests, we can wait: 
```go
for i := 0; i < 10; i++ {
	err := limiter.Wait(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println("Processed request", i)
}
```

`Wait()`:
It blocks until a token can be obtained or its associated context.Context is canceled.

## How this is used in real systems

- API Gateway rate limiting
    
- Per-user limiter (map[userID]*rate.Limiter)
    
- Per-IP limiter in HTTP middleware
    
- Background job throttling
    
- Kafka consumer throttling

Example in HTTP middleware:
```go
func rateLimitMiddleware(next http.Handler) http.Handler {  
	limiter := rate.NewLimiter(10, 20)  
  
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {  
		if !limiter.Allow() {  
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)  
			return  
		}  
		next.ServeHTTP(w, r)  
	})  
}
```

## Architectural Note

`rate.Limiter` works **in-memory**.

That means:

- Works per instance
    
- Not distributed
    
- Not cluster-wide
    

For distributed systems, you’d need:

- Redis-based rate limiter
    
- API Gateway-level limiter
    
- NGINX / Envoy rate limiting
