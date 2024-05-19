package main

import (
	"fmt"
	"sync"
)

type BankAccount struct {
	balance int
	mutex   sync.Mutex
}

func (account *BankAccount) deposit(amount int) {
	account.mutex.Lock()
	fmt.Printf("Einzahlen: %d\n", amount)
	account.balance += amount
	account.mutex.Unlock()
}

func (account *BankAccount) withdraw(amount int) {
	account.mutex.Lock()
	fmt.Printf("Abheben: %d\n", amount)
	if account.balance >= amount {
		account.balance -= amount
	} else {
		fmt.Println("Nicht genügend Guthaben!")
	}
	account.mutex.Unlock()
}

func main() {
	var wg sync.WaitGroup
	account := &BankAccount{balance: 1000}

	transactions := []int{200, -100, 150, -200, 300, -50, 400, -500, 100, -400}

	for _, amount := range transactions {
		wg.Add(1)
		go func(amount int) {
			if amount > 0 {
				account.deposit(amount)
			} else {
				account.withdraw(-amount)
			}
			wg.Done()
		}(amount)
	}
	wg.Wait()
	fmt.Printf("Endgültiger Kontostand: %d\n", account.balance)
}
