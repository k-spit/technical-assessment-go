package main // oder `package main_test` wenn du ausschließlich exportierte Funktionen testen möchtest

import (
	"testing"
	"time"
)

func TestBearbeiteBestellung(t *testing.T) {
	r := &Restaurant{
		startTime:  time.Now(),
		dishCounts: make(map[string]int),
	}
	b := &Bestellung{
		ID:       1,
		Gerichte: []Gericht{{Name: "Pizza", Preis: 12}},
	}

	err := r.bearbeiteBestellung(b)
	if err != nil {
		t.Errorf("Fehler bei der Verarbeitung der Bestellung: %s", err)
	}
	if r.totalRevenue != 12 {
		t.Errorf("totalRevenue erwartet 12, erhalten %d", r.totalRevenue)
	}
	if count, ok := r.dishCounts["Pizza"]; !ok || count != 1 {
		t.Errorf("Erwartete 1 Pizza, erhalten %d", count)
	}
}

func TestGerichteToString(t *testing.T) {
	gerichte := []Gericht{
		{Name: "Pizza", Preis: 12},
		{Name: "Salat", Preis: 9},
	}
	expected := "[Pizza Salat]"
	result := gerichteToString(gerichte)
	if result != expected {
		t.Errorf("Erwartet %s, erhalten %s", expected, result)
	}
}

func TestGeneriereBestellung(t *testing.T) {
	b := generiereBestellung(1)
	if b.ID != 1 {
		t.Errorf("Erwartete Bestell-ID 1, erhalten %d", b.ID)
	}
	if len(b.Gerichte) == 0 {
		t.Errorf("Keine Gerichte generiert")
	}
}
