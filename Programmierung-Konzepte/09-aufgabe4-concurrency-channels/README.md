Aufgabenstellung: Erstellen eines Nachrichtenverteilungssystems

Ziel: Entwickle ein System, das Nachrichten von einem zentralen Eingangspunkt sammelt und sie gleichzeitig an mehrere Verbraucher (Consumer) verteilt. Jeder Consumer soll in der Lage sein, Nachrichten unabhängig von den anderen zu verarbeiten.
Anforderungen:

 * Nachrichtenquelle:
    * Simuliere eine Nachrichtenquelle, die zufällig Nachrichten generiert. Jede Nachricht kann einfach als eine Zeichenkette oder eine strukturierte Datenform (z.B. eine Map oder ein benutzerdefinierter Typ) dargestellt werden.

 * Verbraucher (Consumers):
    * Implementiere mehrere Goroutines als Consumers, die Nachrichten von der Nachrichtenquelle erhalten und verarbeiten.
    * Jeder Consumer sollte in der Lage sein, Nachrichten parallel zu den anderen zu verarbeiten.   

 * Nachrichtenkanal:
    * Verwende Go Channels, um Nachrichten von der Quelle zu den Consumers zu übertragen. Überlege, ob ein unbuffered Channel oder ein buffered Channel für deine Anwendung sinnvoller ist.

 * Kontrollmechanismen:
    * Implementiere Kontrollmechanismen, die es ermöglichen, die Nachrichtenverteilung zu starten, zu überwachen und sicher zu beenden.

 * Fehlerbehandlung:
    * Stelle sicher, dass das System robust gegenüber möglichen Fehlern ist, wie z.B. der Behandlung von blockierten oder geschlossenen Channels.

Details zur Implementierung:

 * Nachrichtenquelle:
 * Erstelle eine Funktion generateMessages, die regelmäßig Nachrichten erzeugt und sie an einen zentralen Channel sendet. Die Frequenz der Nachrichtengenerierung kann variabel sein.

 * Consumer:
 * Erstelle mehrere Consumer-Funktionen, die Nachrichten aus dem Channel lesen und sie verarbeiten. Jede Verarbeitung könnte einfach darin bestehen, die Nachricht auszugeben, eine Transformation durchzuführen oder sie in eine Log-Datei zu schreiben.

 * Start und Beendigung:
 * Das Hauptprogramm sollte in der Lage sein, die Nachrichtengenerierung und die Consumer zu starten und das System ordnungsgemäß zu beenden, z.B. durch Signale von der Tastatur oder nach einer bestimmten Laufzeit.

Beispielkonzept:

```go
package main

import (
    "fmt"
    "time"
)

func generateMessages(msgChan chan<- string) {
    for i := 0; ; i++ {
        msgChan <- fmt.Sprintf("Nachricht %d", i)
        time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond) // Zufälliges Intervall
    }
}

func consumer(id int, msgChan <-chan string) {
    for msg := range msgChan {
        fmt.Printf("Consumer %d erhielt: %s\n", id, msg)
        // Füge hier weitere Verarbeitungslogik hinzu
    }
}

func main() {
    msgChan := make(chan string)
    go generateMessages(msgChan)
    for i := 0; i < 5; i++ {
        go consumer(i, msgChan)
    }

    // Hier könnten Steuerungs- und Beendigungslogiken implementiert werden
}
```