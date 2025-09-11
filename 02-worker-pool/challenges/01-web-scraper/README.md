# Challenge: Web Scraper with Worker Pool

## Problem Statement

You need to build a simple **concurrent web scraper simulator** using the **Worker Pool** pattern.

- You’ll be given a list of URLs.
    
- A fixed number of workers should process these URLs.
    
- Each worker will **“fetch”** the URL (simulate with random sleep).
    
- Each worker should then return a `FetchResult` struct with:
    
    - `url string`
        
    - `status int` (simulated HTTP status, random between `200`, `404`, `500`).
        

Finally, collect all results and print them.

## Requirements

1. Implement a **worker pool** with a configurable number of workers.
    
2. Each worker should pick URLs from the shared channel and process them.
    
3. The main function should:
    
    - Create the workers.
        
    - Feed the URLs into the job channel.
        
    - Wait until all workers finish.
        
    - Print a map of `url -> status`.

```
Worker 1 fetched https://example.com/page1 with status 200
Worker 2 fetched https://example.com/page2 with status 404
Worker 3 fetched https://example.com/page3 with status 500
...
Final results:
map[https://example.com/page1:200 https://example.com/page2:404 ...]
```
