package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type Account struct {
	id      int
	balance float64
	mu      *sync.RWMutex
}

func (a *Account) AccrueInterest(rate float64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.balance += a.balance * (rate / 100.0)
}

func (a *Account) Balance() float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.balance
}

func accrueWorker(ctx context.Context, acc *Account, wg *sync.WaitGroup, updates *int64) {
	defer wg.Done()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			acc.AccrueInterest(rate)
			atomic.AddInt64(updates, 1)
		case <-ctx.Done():
			return
		}
	}
}

func auditorWorker(ctx context.Context, accounts []*Account, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(AUDIT_INTERVAL)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			total_balance := 0.0

			for _, account := range accounts {
				total_balance += account.Balance()
			}

			fmt.Printf("[audit] current total balance: %.2f\n", total_balance)
		case <-ctx.Done():
			return
		}
	}
}

const INTEREST_INTERVAL_MIN_MS = 200 * time.Millisecond
const INTEREST_INTERVAL_MAX_MS = 400 * time.Millisecond

const INTEREST_RATE_MIN = 0.1
const INTEREST_RATE_MAX = 1.0

const INITIAL_BALANCE_MIN = 200.0
const INIITAL_BALANCE_MAX = 600.0

const ACCOUNTS_MIN = 10
const ACCOUNTS_MAX = 20

const AUDIT_MIN = 3
const AUDIT_MAX = 5

const RUN_TIME = 5 * time.Second

const AUDIT_INTERVAL = RUN_TIME / 3

var interval = time.Duration(rand.Int63n(int64(INTEREST_INTERVAL_MAX_MS-INTEREST_INTERVAL_MIN_MS))) + INTEREST_INTERVAL_MIN_MS
var rate = (rand.Float64() * (1.0 - 0.1)) + 0.1

func main() {
	accounts_num := rand.Intn(ACCOUNTS_MAX-ACCOUNTS_MIN) + ACCOUNTS_MIN
	audit_num := rand.Intn(AUDIT_MAX-AUDIT_MIN) + AUDIT_MIN

	fmt.Printf("Starting with:\nInterest accrue interval: %.2f%% each %dms\nAccounts: %d\nAuditors: %d\n", rate, interval.Milliseconds(), accounts_num, audit_num)

	ctx, _ := context.WithTimeout(context.Background(), RUN_TIME)
	wg := &sync.WaitGroup{}
	wg.Add(accounts_num + audit_num)
	var updates int64 = 0
	accounts := []*Account{}

	for id := range accounts_num {
		initial_balance := (rand.Float64() * (INIITAL_BALANCE_MAX - INITIAL_BALANCE_MIN)) + INITIAL_BALANCE_MIN

		account := &Account{
			id:      id,
			balance: initial_balance,
			mu:      &sync.RWMutex{},
		}

		accounts = append(accounts, account)
		go accrueWorker(ctx, account, wg, &updates)
	}

	for range audit_num {
		go auditorWorker(ctx, accounts, wg)
	}

	wg.Wait()

	total_balance := 0.0

	for _, account := range accounts {
		total_balance += account.Balance()
	}

	fmt.Printf("Interest updates: %d\n", updates)
	fmt.Printf("Final total balance: %.2f\n", total_balance)
}
