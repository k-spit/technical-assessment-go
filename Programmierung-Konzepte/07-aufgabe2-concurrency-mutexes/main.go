package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

type Restaurant struct {
	totalRevenue           int
	mutex                  sync.Mutex
	startTime              time.Time
	currentIntervalOrders  int
	currentIntervalRevenue int
	dishCounts             map[string]int
}

type Gericht struct {
	Name  string
	Preis int
}

type Bestellung struct {
	ID       int
	Gerichte []Gericht
}

func (r *Restaurant) bearbeiteBestellung(b *Bestellung) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.currentIntervalOrders++
	for _, gericht := range b.Gerichte {
		r.totalRevenue += gericht.Preis
		r.currentIntervalRevenue += gericht.Preis
		if r.dishCounts == nil {
			r.dishCounts = make(map[string]int)
		}
		r.dishCounts[gericht.Name]++
	}
	log.Printf("Live: Bestellung %d bearbeitet. Gericht: %s, Einnahmen bisher: %d\n", b.ID, gerichteToString(b.Gerichte), r.totalRevenue)
	return nil
}

func gerichteToString(gerichte []Gericht) string {
	names := []string{}
	for _, g := range gerichte {
		names = append(names, g.Name)
	}
	return fmt.Sprintf("%v", names)
}

func simuliereRestaurant() {
	r := &Restaurant{
		startTime:  time.Now(),
		dishCounts: make(map[string]int),
	}
	ticker := time.NewTicker(60 * time.Second)
	go func() {
		for range ticker.C {
			r.mutex.Lock()
			log.Printf("\nZwischenbericht nach %v:\n", time.Since(r.startTime))
			log.Printf("Anzahl der Bestellungen: %d\n", r.currentIntervalOrders)
			log.Printf("Einnahmen in diesem Intervall: %d\n", r.currentIntervalRevenue)
			log.Printf("Beliebteste Gerichte in diesem Intervall:\n")
			for dish, count := range r.dishCounts {
				log.Printf("- %s: %d\n", dish, count)
			}
			r.currentIntervalOrders = 0
			r.currentIntervalRevenue = 0
			r.dishCounts = make(map[string]int)
			r.mutex.Unlock()
		}
	}()

	orderID := 1
	for {
		time.Sleep(time.Duration(rand.ExpFloat64() * float64(time.Second)))
		bestellung := generiereBestellung(orderID)
		r.bearbeiteBestellung(bestellung)
		orderID++
	}
}

func generiereBestellung(id int) *Bestellung {
	gerichte := []Gericht{
		{Name: "Pizza", Preis: 12},
		{Name: "Salat", Preis: 9},
		{Name: "Pasta", Preis: 15},
		{Name: "Suppe", Preis: 7},
	}
	anzahlGerichte := rand.Intn(3) + 1
	bestellteGerichte := make([]Gericht, anzahlGerichte)
	for i := 0; i < anzahlGerichte; i++ {
		bestellteGerichte[i] = gerichte[rand.Intn(len(gerichte))]
	}
	return &Bestellung{
		ID:       id,
		Gerichte: bestellteGerichte,
	}
}

func main() {
	go simuliereRestaurant()
	log.Println("Restaurant Simulation läuft. Drücken Sie CTRL+C zum Beenden.")
	select {} // Blockiert das Hauptprogramm, ohne CPU-Last zu verursachen
}
