package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Account struct {
	id      int
	balance uint
	*sync.RWMutex
}

func NewAccount(id int, initialBalance uint) *Account {
	return &Account{
		id:      id,
		balance: initialBalance,
		RWMutex: &sync.RWMutex{},
	}
}

func (a *Account) Withdraw(amount uint) error {
	if a.balance >= amount {
		a.balance -= amount
		return nil
	}

	return errors.New("not enough balance")
}

func (a *Account) DepositTo(amount uint, destination *Account) {
	first, second := a, destination

	if first.id > second.id {
		first, second = destination, a
	}

	first.Lock()
	second.Lock()
	defer second.Unlock()
	defer first.Unlock()

	err := a.Withdraw(amount)

	if err == nil {
		destination.Deposit(amount)
		time.Sleep(time.Duration(rand.Intn(int(SLEEP_MAX)-int(SLEEP_MIN)) + int(SLEEP_MIN)))
		fmt.Printf("Transfer %d from #%d → #%d (OK)\n", amount, a.id, destination.id)
		return
	}

	fmt.Printf("Transfer %d from #%d → #%d (Insufficient funds)\n", amount, a.id, destination.id)
}

func (a *Account) GetBalance() uint {
	a.RLock()
	defer a.RUnlock()
	return a.balance
}

func (a *Account) Deposit(amount uint) {
	a.balance += amount
}

const (
	ACCOUNTS_NUM                   = 10
	TOTAL_BALANCE                  = 100_000
	INITIAL_BALANCE_PER_ACCOUNT    = TOTAL_BALANCE / ACCOUNTS_NUM
	CONCURRENT_TRANSFER_OPERATIONS = 500
	SLEEP_MIN                      = 100 * time.Millisecond
	SLEEP_MAX                      = 500 * time.Millisecond
)

func main() {
	accounts := []*Account{}
	wg := &sync.WaitGroup{}
	wg.Add(CONCURRENT_TRANSFER_OPERATIONS)

	for i := 0; i < ACCOUNTS_NUM; i++ {
		accounts = append(accounts, NewAccount(i+1, INITIAL_BALANCE_PER_ACCOUNT))
	}

	for i := 0; i < CONCURRENT_TRANSFER_OPERATIONS; i++ {
		go func() {
			defer wg.Done()
			originIndex := 0
			destinationIndex := 0

			for originIndex == destinationIndex {
				originIndex = rand.Intn(ACCOUNTS_NUM)
				destinationIndex = rand.Intn(ACCOUNTS_NUM)
			}

			amount := uint(rand.Intn(INITIAL_BALANCE_PER_ACCOUNT) + 1)
			accounts[originIndex].DepositTo(amount, accounts[destinationIndex])
		}()
	}

	wg.Wait()
	var total_balance uint = 0

	for _, account := range accounts {
		total_balance += account.GetBalance()
	}

	fmt.Printf("Final total balance: %d (initial total: %d)", total_balance, TOTAL_BALANCE)

}
