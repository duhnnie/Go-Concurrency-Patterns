## Advanced Pipeline Challenge: Parallel Word Counters with Merge

**Goal:**  
Build a pipeline that processes text with the following rules:

1. **Generator** → emits sentences.
    
2. **Tokenizer** → splits sentences into words.
    
3. **Normalizer** → lowercases & removes punctuation.
    
4. **Filter** → keep only words with length ≥ 5.
    
5. **Fan-out** → distribute words to **N counters running in parallel**.
    
6. **Fan-in** → merge the maps from those counters.
    
7. **Final Reducer** → combine all partial maps into one global word count.
    

---

### 🔧 New Concepts

- **Fan-out**: We’ll send each word to one of several counters (e.g., round-robin or by hash of word).
    
- **Fan-in**: Merge results from the counters (map[string]int) into one channel.
    
- **Reducer**: Merge maps into a single final map.
    

---

## Example Input

```go
[]string{
    "Go is expressive, concise, clean, and efficient.",
	"Concurrency is not parallelism.",
	"Channels orchestrate communication.",
	"Parallel pipelines can improve performance.",
	"Goroutines are lightweight threads.",
}
```
## Example Output

If you run it on the sentences above, you might get something like:

```go
map[channels:1 communication:1 concurrency:1 efficient:1 goroutines:1 lightweight:1 orchestrate:1 parallel:1 parallelism:1 pipelines:1]
```