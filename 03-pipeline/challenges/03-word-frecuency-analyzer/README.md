# Pipeline Challenge 2: Word Frequency Analyzer

## Problem Statement

You are tasked with building a **pipeline** in Go that analyzes a collection of sentences and produces a **map of word frequencies**, but with the following requirements:

1. **Stage 1 – Generator**  
    Emit all sentences into a channel.
    
2. **Stage 2 – Tokenizer**  
    Split each sentence into words and emit them individually.
    
3. **Stage 3 – Normalizer**  
    Normalize words by:
    
    - Lowercasing them
        
    - Removing punctuation (e.g., `, . ! ?`)
        
4. **Stage 4 – Filter**  
    Only pass words with length ≥ 4 characters.
    
5. **Stage 5 – Counter**  
    Maintain a frequency map of the words and output the final result.
    
## Example Input

```go
[]string{
    "Go is expressive, concise, clean, and efficient.",
    "Concurrency is not parallelism.",
    "Channels orchestrate communication.",
}
```

## Expected Output

```go
map[string]int{
    "expressive": 1,
    "concise":    1,
    "clean":      1,
    "efficient":  1,
    "concurrency":1,
    "parallelism":1,
    "channels":   1,
    "orchestrate":1,
    "communication":1,
}
```