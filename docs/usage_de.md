## Wie man `go-redmine` benutzt

## Verwendung von `go-redmine` als API-Client

`go-redmine` in ein Go-Projekt installieren:

```bash
go get github.com/cloudogu/go-redmine
```

So schnell geht es, einen go-redmine-Client zu konfigurieren und zu instanziieren:

```go
package main

import "github.com/cloudogu/go-redmine"
...
func main() {
  client, err := redmine.NewClientBuilder().
    Endpoint("http://example.com:3000").
  	AuthBasicAuth("admin", "password").
  	Build()

	project := redmine.Project{
		Name:           "Heart of Gold",
		Identifier:     "heartofgold",
		Description:    "A Redmine example project",
	}
  
  updatedProject, err := client.CreateProject(project)
}
```

## `go-redmine` als Konsolenanwendung `godmine`

Kommandozeilenwerkzeug für Redmine-Zugriffe.

### Verwendung

    godmine <Befehl> <Unterbefehl> [Argumente]
    
    Projekt-Befehle:
      Projekt mit Texteditor anlegen.
                 $ godmine p a
    
      create c Projekt aus gegebenen Argumenten erstellen.
                 $ godmine p c Name Bezeichner Beschreibung
    
      update u Update des gegebenen Projekts.
                 $ godmine p u 1
    
      show s zeigt gegebenes Projekt an.
                 $ godmine p s 1
    
      delete d angegebenes Projekt löschen.
                 $ godmine p d 1
    
      list l Projekte auflisten.
                 $ godmine p l
    
    Issue-Befehle:
      add a issue mit Texteditor erstellen.
                 $ godmine i a
    
      create c Issue aus gegebenen Argumenten erstellen.
                 $ godmine i c Thema Beschreibung
    
      update u Aktualisiere gegebenes Issue.
                 $ godmine i u 1
    
      show s zeigt den angegebenen Eintrag an.
                 $ godmine i s 1
    
      delete d löscht den angegebenen Eintrag.
                 $ godmine i d 1
    
      close x schließt den angegebenen Eintrag.
                 $ godmine i x 1
    
      notes n Hinzufügen von Notizen zu einer bestimmten Ausgabe.
                 $ godmine i n 1
    
      list l Ausgaben auflisten.
                 $ godmine i l

### Einstellungen

Um das CLI-Tool zu verwenden, sollten eine `settings.json` verwendet werden:

UNIX:

    ~/.config/godmine/settings.json

WINDOWS:

    %APPDATA%\godmine\settings.json

Folgendes schreiben:

    {
    	"endpoint": "http://redmine.example.com",
    	"apikey": "IHR-API-KEY",
    	"project": 1 // Standard-Projekt-ID
    }

Wenn Sie die Konfigurationsdatei wechseln wollen, sollten Sie die Umgebungsvariable `GODMINE_ENV` verwenden.
Wenn Sie `GODMINE_ENV` auf *mine* setzen, verwendet Godmine `settings.mine.json` als Konfigurationsdatei.

## Ein Hinweis zu Redmine-Versionen

`go-redmine` sollte gegen neuere Versionen von Redmine verwendet werden. Projekte und Issues funktionieren definitiv mit Redmine v4.1.x

Übersetzt mit www.DeepL.com/Translator (kostenlose Version)