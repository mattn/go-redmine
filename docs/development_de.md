# go-redmine entwickeln

## Erstellen

### Tools

Das Bauen erfordert diese installierten Werkzeuge. Im Allgemeinen verlässt sich dieses Projekt nicht auf topaktuelle oder experimentelle Versionen. Halbwegs aktuelle Versionen sind in Ordnung, aber wenn das Bauen fehlschlägt, sollte überprüft werden, ob stark veraltete Versionen der jeweiligen Tools vorliegen:

- Make
   - f. e. GNU Make 4.2.1
- Docker
   - f. e. Client/server version 19.03.5
- Golang compiler
   - f. e. 1.14.12

### Binary lokal bauen und andere interessante `make`-Targets

make-Target | Aktion
------------|-------
`make` / `make compile` | baut das Provider-Binary in einem Go-Container
