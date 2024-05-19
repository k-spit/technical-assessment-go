Aufgabenstellung: Ein Bankkonto mit mehreren Zugriffen

Ziel: Implementiere ein einfaches simuliertes Bankkontosystem, das mehrere gleichzeitige Ein- und Auszahlungen auf ein Konto erlaubt, und nutze dabei Mutexes, um sicherzustellen, dass die Transaktionen korrekt verarbeitet werden.
Details:

* Kontostand: Starte mit einem Anfangskontostand von 1000 Euro.
* Transaktionen: Es gibt zwei Arten von Transaktionen: Einzahlungen und Abhebungen. Für diese Aufgabe kannst du eine Liste von Transaktionen simulieren, die zufällig generiert oder vordefiniert sein können.
* Concurrency: Jede Transaktion soll in einer eigenen Goroutine ausgeführt werden.
* Synchronisation: Verwende einen sync.Mutex um sicherzustellen, dass jede Transaktion korrekt auf den Kontostand zugreift, ohne dass Race Conditions entstehen.

Aufgaben:

* Implementiere die Transaktionen: Nutze den bereitgestellten Code-Rahmen, um das Bankkontosystem zu implementieren. Stelle sicher, dass die Transaktionen korrekt ausgeführt werden und keine Race Conditions auftreten.
* Teste verschiedene Szenarien: Ändere die Liste der Transaktionen oder füge künstliche Verzögerungen hinzu, um zu sehen, wie dein System unter verschiedenen Bedingungen funktioniert.
* Beobachte das Ergebnis: Überprüfe, ob der endgültige Kontostand korrekt ist und keine Fehler bei den Transaktionen aufgetreten sind.