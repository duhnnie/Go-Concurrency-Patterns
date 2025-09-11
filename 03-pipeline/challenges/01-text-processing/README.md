# Challenge: Text Processing Pipeline

You need to implement a **concurrent pipeline** in Go that processes sentences through multiple stages:

## Stages

1. **Generator**:
    
    - Takes a slice of sentences (`[]string`) and emits them one by one into a channel.
        
2. **Tokenizer**:
    
    - Reads sentences from the input channel.
        
    - Splits each sentence into words.
        
    - Emits each word into the next channel.
        
3. **Filter**:
    
    - Reads words and filters out words shorter than 4 characters.
        
    - Emits only words with length >= 4.
        
4. **Uppercaser**:
    
    - Converts the filtered words to uppercase.
        
    - Emits the uppercase words.
        
5. **Collector**:
    
    - Collects all final words and returns them as a slice (`[]string`).

### Example input

```go
sentences := []string{
    "Go is expressive concise clean and efficient",
    "Concurrency is not parallelism",
    "Channels orchestrate communication",
}
```

### Example ouput

```go
[EXPRESSIVE CONCISE CLEAN EFFICIENT CONCURRENCY PARALLELISM CHANNELS ORCHESTRATE COMMUNICATION]
```

### Constraints

- Each stage must run as a separate goroutine.
    
- Communication between stages happens **only through channels**.
    
- No global state or shared variables allowed (other than channels).