# Pipeline Challenge: Prime Number Filtering

## Problem Statement

You need to build a pipeline that processes a stream of integers and filters only **prime numbers**. The pipeline should have the following stages:

1. **Generator**: Emits integers from `2` up to `N`.
    
2. **FilterEven**: Removes even numbers greater than 2 (optimization step).
    
3. **PrimeChecker**: Checks if a number is prime. Only passes primes to the next stage.
    
4. **Collector**: Collects all primes into a slice and outputs the final list.
    
## Requirements

- Each stage should be a **separate goroutine** communicating through channels.
    
- All channels must be **closed properly** to avoid leaks.
    
- The final output should be a slice of primes.
    
- Test with `N = 50` (so your result should include `[2, 3, 5, 7, 11, ... 47]`).

## Example Expected Output

`Primes up to 50: [2 3 5 7 11 13 17 19 23 29 31 37 41 43 47]`