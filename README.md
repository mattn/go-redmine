# go-redmine

Intefaces to redmine.

## Module Install

```
go get github.com/mattn/go-redmine 
```

## APIs

Provide Interfaces to redmine APIs.

|API                |Implements|
|-------------------|---------:|
|Issues             |      100%|
|Projects           |      100%|
|Project Memberships|      100%|
|Users              |        0%|
|Time Entries       |      100%|
|News               |      100%|
|Issue Relations    |      100%|
|Versions           |      100%|
|Wiki Pages         |      100%|
|Queries            |        0%|
|Attachments        |        0%|
|Issue Statuses     |      100%|
|Trackers           |      100%|
|Enumerations       |      100%|
|Issue Categories   |      100%|
|Roles              |      100%|
|Groups             |        0%|

## Godmine

Provide command line tool for redmine.

## Install

```
go install github.com/mattn/go-redmine/cmd/godmine@latest
```

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

# Settings

To use this, you should create `settings.json` in:

UNIX:

    ~/.config/godmine/settings.json

WINDOWS:

    %APPDATA%\godmine/settings.json

Write following:

    {
    	"endpoint": "http://redmine.example.com",
    	"apikey": "YOUR-API-KEY",
    	"project": 1 // default project id
    }

If you want switching configuration file, you should use `GODMINE_ENV` environment variable.
If you set `GODMINE_ENV` to *mine*, godmine use `settings.mine.json` to configuration file.

# License

MIT

# Author

Yasuhiro Matsumoto (a.k.a mattn)
