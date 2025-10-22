package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const NUM_OF_WITHDRAWS = 10
const NUM_OF_DEPOSITS = 8
const AMOUNT_MIN = 1
const AMOUNT_MAX = 15

type Account struct {
	balance float64
	mu      *sync.RWMutex
}

func NewAccount(initialBalance float64) *Account {
	return &Account{
		balance: initialBalance,
		mu:      &sync.RWMutex{},
	}
}

func (a *Account) GetBalance() float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.balance
}

func (a *Account) Deposit(amount float64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.balance += amount
}

func (a *Account) Withdraw(amount float64) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.balance >= amount {
		a.balance -= amount
		return nil
	}

	return errors.New("Not enough balance")
}

func GetRandomAmount() float64 {
	diff := AMOUNT_MAX - AMOUNT_MIN
	return float64(rand.Intn(diff) + AMOUNT_MIN)
}

func GetRandomSleepTime() time.Duration {
	return time.Duration(rand.Intn(500) * int(time.Millisecond))
}

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(NUM_OF_DEPOSITS + NUM_OF_WITHDRAWS)
	a := NewAccount(20)

	go func() {
		for range NUM_OF_DEPOSITS {
			amount := GetRandomAmount()
			a.Deposit(amount)
			fmt.Printf("Deposit: %.2f\n", amount)
			time.Sleep(GetRandomSleepTime())
			wg.Done()
		}
	}()

	go func() {
		for range NUM_OF_WITHDRAWS {
			amount := GetRandomAmount()
			err := a.Withdraw(amount)

			if err != nil {
				fmt.Printf("Can't withdraw %.2f: %s\n", amount, err)
			} else {
				fmt.Printf("Withdraw: $%.2f\n", amount)
			}

			time.Sleep(GetRandomSleepTime())
			wg.Done()
		}
	}()

	wg.Wait()
	fmt.Printf("Final Balance: %.2f\n", a.GetBalance())
}
