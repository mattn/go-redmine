# How to use `go-redmine`

## Using `go-redmine` as API client

Install `go-redmine` to your Go project:

```bash
go get github.com/cloudogu/go-redmine
```

Configure and instantiate a go-redmine client, and you are ready to interact with Redmine:

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

## Using `go-redmine` as CLI interface `godmine`

Provide command line tool for redmine.

### Usage

    godmine <command> <subcommand> [arguments]
    
    Project Commands:
      add      a create project with text editor.
                 $ godmine p a
    
      create   c create project from given arguments.
                 $ godmine p c name identifier description
    
      update   u update given project.
                 $ godmine p u 1
    
      show     s show given project.
                 $ godmine p s 1
    
      delete   d delete given project.
                 $ godmine p d 1
    
      list     l listing projects.
                 $ godmine p l
    
    Issue Commands:
      add      a create issue with text editor.
                 $ godmine i a
    
      create   c create issue from given arguments.
                 $ godmine i c subject description
    
      update   u update given issue.
                 $ godmine i u 1
    
      show     s show given issue.
                 $ godmine i s 1
    
      delete   d delete given issue.
                 $ godmine i d 1
    
      close    x close given issue.
                 $ godmine i x 1
    
      notes    n add notes to given issue.
                 $ godmine i n 1
    
      list     l listing issues.
                 $ godmine i l

### Settings

To use this, you should create `settings.json` in:

UNIX:

    ~/.config/godmine/settings.json

WINDOWS:

    %APPDATA%\godmine\settings.json

Write following:

    {
    	"endpoint": "http://redmine.example.com",
    	"apikey": "YOUR-API-KEY",
    	"project": 1 // default project id
    }

If you want switching configuration file, you should use `GODMINE_ENV` environment variable.
If you set `GODMINE_ENV` to *mine*, godmine use `settings.mine.json` to configuration file.

## A Note on Redmine versions

`go-redmine` is supposed to be used against more recent versions of Redmine. Projects and issues definitely work with Redmine v4.1.x