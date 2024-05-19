OSINT
===

# Übung: Bedrohungsmodellierung

Bedrohungsmodellierung ist ein Prozess zur Identifizierung, Bewertung und Minderung von Risiken in einem System. Wir werden ein Bedrohungsmodell für ein einfaches Webanwendungssystem erstellen.
Szenario: Webanwendung mit Benutzerauthentifizierung

Komponenten:

* Webserver (NGINX)
* API-Server (Go, JWT-basierte Authentifizierung)
* Datenbank (PostgreSQL)
* Cache (Redis)
* Client (Webbrowser)

Bedrohungsmodell:

* Identifikation von Assets:
  * Benutzerkonten
  * Datenbank mit Benutzerdaten
  * API-Schlüssel und JWT-Tokens
  * Cache-Daten

  Identifikation potenzieller Bedrohungen:
  * SQL-Injection: Angreifer könnten versuchen, über SQL-Injection Zugriff auf die Datenbank zu erhalten.
  * Cross-Site Scripting (XSS): Angreifer könnten versuchen, schädlichen Code auf der Webseite zu injizieren.
  * Man-in-the-Middle (MITM): Angreifer könnten versuchen, den Datenverkehr zwischen Client und Server abzufangen.
  * Brute Force: Angreifer könnten versuchen, Benutzerpasswörter durch wiederholte Versuche zu erraten.
  * Denial of Service (DoS): Angreifer könnten versuchen, den Dienst durch Überlastung außer Betrieb zu setzen.

* Minderungsmaßnahmen:
  * SQL-Injection: Verwenden von vorbereiteten Anweisungen und ORM (Object-Relational Mapping).
  * XSS: Validierung und Escape von Benutzereingaben.
  * MITM: Einsatz von HTTPS für die gesamte Kommunikation.
  * Brute Force: Implementierung von Rate Limiting und Captchas.
  * DoS: Implementierung von Rate Limiting und Lastverteilung.

