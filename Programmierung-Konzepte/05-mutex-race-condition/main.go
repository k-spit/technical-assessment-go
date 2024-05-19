package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	balance int
	mutex   sync.Mutex
)

func deposit(value int, wg *sync.WaitGroup) {
	// mutex.Lock() // Sperrt den Zugriff für andere Goroutines
	temp := balance
	time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond) // Zufällige Verzögerung
	balance = temp + value
	// mutex.Unlock() // Gibt den Zugriff frei
	wg.Done()
}

func main() {
	var wg sync.WaitGroup
	transactions := []int{50, -20, 70, -30, 100, -50, 30, -10, 60, 10, -5, 15, -25, 40, -10}

	for _, value := range transactions {
		wg.Add(1)
		go deposit(value, &wg)
	}
	wg.Wait()
	fmt.Printf("Endgültiger Kontostand: %d\n", balance)
}
