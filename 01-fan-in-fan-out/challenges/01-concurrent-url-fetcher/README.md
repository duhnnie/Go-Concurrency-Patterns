## **Challenge: Concurrent URL Fetcher**

**Goal:**

- Given a list of URLs, fetch their HTTP status codes concurrently using a worker pool (Fan-Out).
    
- Collect all the results into a single map of `URL -> Status Code` (Fan-In).
    

**Requirements:**

1. Create a list of 10-20 URLs (can be real or dummy like `https://httpbin.org/status/200`, `https://httpbin.org/status/404`).
    
2. Use **3 worker goroutines** to fetch URLs concurrently.
    
3. Each worker reads URLs from a `jobs` channel and writes results to a `results` channel.
    
4. Collect results in a map in the main function and print them.
    

**Hints:**

- Use `net/http` to fetch URLs.
    
- Use `close(jobs)` after sending all URLs.
    
- Remember to safely read from `results` channel in the main function.

