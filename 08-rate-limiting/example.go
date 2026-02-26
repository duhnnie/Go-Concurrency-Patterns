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
