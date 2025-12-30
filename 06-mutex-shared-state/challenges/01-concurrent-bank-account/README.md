# Challenge: Concurrent Bank Account

## Goal

Implement a simple `Account` type that supports:

- concurrent deposits and withdrawals
    
- a method to check balance safely
    
## Constraints

- multiple goroutines will call deposit/withdraw simultaneously
    
- must prevent race conditions
    
- simulate random deposits/withdrawals
    
- print final balance after all goroutines complete