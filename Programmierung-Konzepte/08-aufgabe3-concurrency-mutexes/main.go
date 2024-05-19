package main

import (
	"log"
	"sync"
)

type Konto struct {
	Kontostand int
	sync.Mutex
}

func (k *Konto) einzahlen(amount int) {
	k.Lock()
	defer k.Unlock()

	k.Kontostand += amount
}

func (k *Konto) abheben(amount int) {
	k.Lock()
	defer k.Unlock()

	if k.Kontostand-amount >= 0 {
		k.Kontostand -= amount
	} else {
		log.Printf("kontostand(%d) zu gering um geld abzuheben\n", k.Kontostand)
	}
}

func main() {
	var wg sync.WaitGroup
	k := &Konto{Kontostand: 1000}

	transaktionen := []int{100, 200, 900, -500, 100, 800, -600}

	for _, t := range transaktionen {
		wg.Add(1)
		go func(t int) {
			if t > 0 {
				k.einzahlen(t)
			} else {
				k.abheben(-t)
			}
			wg.Done()
		}(t)
	}

	wg.Wait()
	log.Printf("kontostand: %d\n", k.Kontostand)
}
