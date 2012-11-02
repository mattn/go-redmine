package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mattn/go-iconv"
	"github.com/mattn/go-redmine"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func fatal(format string, err error) {
	fmt.Fprintf(os.Stderr, format, err)
	os.Exit(1)
}

func run(argv []string) error {
	cmd, err := exec.LookPath(argv[0])
	if err != nil {
		return err
	}
	var stdin *os.File
	if runtime.GOOS == "windows" {
		stdin, _ = os.Open("CONIN$")
	} else {
		stdin = os.Stdin
	}
	p, err := os.StartProcess(cmd, argv, &os.ProcAttr{Files: []*os.File{stdin, os.Stdout, os.Stderr}})
	if err != nil {
		return err
	}
	defer p.Release()
	w, err := p.Wait()
	if err != nil {
		return err
	}
	if !w.Exited() || !w.Success() {
		return errors.New("Failed to execute text editor")
	}
	return nil
}

func notesFromEditor() (string, error) {
	file := ""
	newf := fmt.Sprintf("%d.txt", rand.Int())
	if runtime.GOOS == "windows" {
		file = filepath.Join(os.Getenv("APPDATA"), "godmine", newf)
	} else {
		file = filepath.Join(os.Getenv("HOME"), ".config", "godmine", newf)
	}
	defer os.Remove(file)
	editor := os.Getenv("EDITOR")
	if len(editor) == 0 {
		if runtime.GOOS == "windows" {
			editor = "notepad"
		} else {
			editor = "vim"
		}
	}

	contents := "### Notes Here ###\n"

	ioutil.WriteFile(file, []byte(contents), 0600)

	if err := run([]string{editor, file}); err != nil {
		return "", err
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	text := string(b)

	if text == contents {
		return "", errors.New("Canceled")
	}
	return text, nil
}

func issueFromEditor(contents string) (*redmine.Issue, error) {
	file := ""
	newf := fmt.Sprintf("%d.txt", rand.Int())
	if runtime.GOOS == "windows" {
		file = filepath.Join(os.Getenv("APPDATA"), "godmine", newf)
	} else {
		file = filepath.Join(os.Getenv("HOME"), ".config", "godmine", newf)
	}
	defer os.Remove(file)
	editor := os.Getenv("EDITOR")
	if len(editor) == 0 {
		if runtime.GOOS == "windows" {
			editor = "notepad"
		} else {
			editor = "vim"
		}
	}

	if contents == "" {
		contents = "### Subject Here ###\n### Description Here ###\n"
	}

	ioutil.WriteFile(file, []byte(contents), 0600)

	if err := run([]string{editor, file}); err != nil {
		return nil, err
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	text := string(b)

	if text == contents {
		return nil, errors.New("Canceled")
	}
	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		return nil, errors.New("Canceled")
	}
	var subject, description string
	if len(lines) == 1 {
		subject = lines[0]
	} else {
		subject, description = lines[0], strings.Join(lines[1:], "\n")
	}
	var issue redmine.Issue
	issue.Subject = subject
	issue.Description = description
	return &issue, nil
}

func projectFromEditor(contents string) (*redmine.Project, error) {
	file := ""
	newf := fmt.Sprintf("%d.txt", rand.Int())
	if runtime.GOOS == "windows" {
		file = filepath.Join(os.Getenv("APPDATA"), "godmine", newf)
	} else {
		file = filepath.Join(os.Getenv("HOME"), ".config", "godmine", newf)
	}
	defer os.Remove(file)
	editor := os.Getenv("EDITOR")
	if len(editor) == 0 {
		if runtime.GOOS == "windows" {
			editor = "notepad"
		} else {
			editor = "vim"
		}
	}

	if contents == "" {
		contents = "### Name Here ###\n### Identifier Here ###\n### Description Here ###\n"
	}

	ioutil.WriteFile(file, []byte(contents), 0600)

	if err := run([]string{editor, file}); err != nil {
		return nil, err
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	text := string(b)

	if text == contents {
		return nil, errors.New("Canceled")
	}
	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		return nil, errors.New("Canceled")
	}
	var name, identifier, description string
	if len(lines) == 1 {
		return nil, errors.New("Invalid Format")
	} else if len(lines) == 2 {
		name = lines[0]
		identifier = lines[1]
	} else {
		name = lines[0]
		identifier = lines[1]
		description = lines[2]
	}
	var project redmine.Project
	project.Name = name
	project.Identifier = identifier
	project.Description = description
	return &project, nil
}

func getConfig() config {
	file := ""
	if runtime.GOOS == "windows" {
		file = filepath.Join(os.Getenv("APPDATA"), "godmine", "settings.json")
	} else {
		file = filepath.Join(os.Getenv("HOME"), ".config", "godmine", "settings.json")
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		fatal("Failed to read config file: %s\n", err)
	}
	var c config
	err = json.Unmarshal(b, &c)
	if err != nil {
		fatal("Failed to unmarshal file: %s\n", err)
	}
	return c
}

func addIssue(subject, description string) {
	config := getConfig()
	var issue redmine.Issue
	c := redmine.NewClient(config.Endpoint, config.Apikey)
	issue.ProjectId = config.Project
	issue.Subject = subject
	issue.Description = description
	_, err := c.CreateIssue(issue)
	if err != nil {
		fatal("Failed to create issue: %s\n", err)
	}
}

func createIssue() {
	config := getConfig()
	issue, err := issueFromEditor("")
	if err != nil {
		fatal("%s\n", err)
	}
	c := redmine.NewClient(config.Endpoint, config.Apikey)
	issue.ProjectId = config.Project
	_, err = c.CreateIssue(*issue)
	if err != nil {
		fatal("Failed to create issue: %s\n", err)
	}
}

func updateIssue(id int) {
	config := getConfig()

	c := redmine.NewClient(config.Endpoint, config.Apikey)
	issue, err := c.Issue(id)
	if err != nil {
		fatal("Failed to update issue: %s\n", err)
	}
	issueNew, err := issueFromEditor(fmt.Sprintf("%s\n%s\n", issue.Subject, issue.Description))
	if err != nil {
		fatal("%s\n", err)
	}
	issue.Subject = issueNew.Subject
	issue.Description = issueNew.Description
	issue.ProjectId = config.Project
	err = c.UpdateIssue(*issue)
	if err != nil {
		fatal("Failed to update issue: %s\n", err)
	}
}

func deleteIssue(id int) {
	config := getConfig()

	c := redmine.NewClient(config.Endpoint, config.Apikey)
	err := c.DeleteIssue(id)
	if err != nil {
		fatal("Failed to delete issue: %s\n", err)
	}
}

func closeIssue(id int) {
	config := getConfig()

	c := redmine.NewClient(config.Endpoint, config.Apikey)
	issue, err := c.Issue(id)
	if err != nil {
		fatal("Failed to update issue: %s\n", err)
	}
	is, err := c.IssueStatuses()
	if err != nil {
		fatal("Failed to get issue statuses: %s\n", err)
	}
	for _, s := range is {
		if s.IsClosed {
			issue.StatusId = s.Id
			err = c.UpdateIssue(*issue)
			if err != nil {
				fatal("Failed to update issue: %s\n", err)
			}
			break
		}
	}
}

func notesIssue(id int) {
	config := getConfig()

	content, err := notesFromEditor()
	if err != nil {
		fatal("%s\n", err)
	}
	c := redmine.NewClient(config.Endpoint, config.Apikey)
	issue, err := c.Issue(id)
	if err != nil {
		fatal("Failed to update issue: %s\n", err)
	}
	issue.Notes = content
	issue.ProjectId = config.Project
	err = c.UpdateIssue(*issue)
	if err != nil {
		fatal("Failed to update issue: %s\n", err)
	}
}

func showIssue(id int) {
	config := getConfig()
	c := redmine.NewClient(config.Endpoint, config.Apikey)
	issue, err := c.Issue(id)
	if err != nil {
		fatal("Failed to show issue: %s\n", err)
	}
	assigned := ""
	if issue.Assigned != nil {
		assigned = issue.Assigned.Name
	}

	fmt.Printf(`
Id: %d
Subject: %s
Project: %s
Tracker: %s
Status: %s
Priority: %s
Author: %s
Assigned: %s
CreatedOn: %s
UpdatedOn: %s

%s
`[1:],
		issue.Id,
		issue.Subject,
		issue.Project.Name,
		issue.Tracker.Name,
		issue.Status.Name,
		issue.Priority.Name,
		issue.Author.Name,
		assigned,
		issue.CreatedOn,
		issue.UpdatedOn,
		issue.Description)
}

func listIssues() {
	config := getConfig()
	c := redmine.NewClient(config.Endpoint, config.Apikey)
	issues, err := c.Issues()
	if err != nil {
		fatal("Failed to list issues: %s\n", err)
	}
	for _, i := range issues {
		fmt.Printf("%4d: %s\n", i.Id, i.Subject)
	}
}

func addProject(name, identifier, description string) {
	config := getConfig()
	var project redmine.Project
	c := redmine.NewClient(config.Endpoint, config.Apikey)
	project.Name = name
	project.Identifier = identifier
	project.Description = description
	_, err := c.CreateProject(project)
	if err != nil {
		fatal("Failed to create project: %s\n", err)
	}
}

func createProject() {
	config := getConfig()
	project, err := projectFromEditor("")
	if err != nil {
		fatal("%s\n", err)
	}
	c := redmine.NewClient(config.Endpoint, config.Apikey)
	_, err = c.CreateProject(*project)
	if err != nil {
		fatal("Failed to create project: %s\n", err)
	}
}

func updateProject(id int) {
	config := getConfig()

	c := redmine.NewClient(config.Endpoint, config.Apikey)
	project, err := c.Project(id)
	if err != nil {
		fatal("Failed to update project: %s\n", err)
	}
	projectNew, err := projectFromEditor(fmt.Sprintf("%s\n%s\n%s\n", project.Name, project.Identifier, project.Description))
	if err != nil {
		fatal("%s\n", err)
	}
	project.Name = projectNew.Name
	project.Identifier = projectNew.Identifier
	project.Description = projectNew.Description
	err = c.UpdateProject(*project)
	if err != nil {
		fatal("Failed to update project: %s\n", err)
	}
}

func deleteProject(id int) {
	config := getConfig()

	c := redmine.NewClient(config.Endpoint, config.Apikey)
	err := c.DeleteProject(id)
	if err != nil {
		fatal("Failed to delete project: %s\n", err)
	}
}

func showProject(id int) {
	config := getConfig()
	c := redmine.NewClient(config.Endpoint, config.Apikey)
	project, err := c.Project(id)
	if err != nil {
		fatal("Failed to show project: %s\n", err)
	}

	fmt.Printf(`
Id: %d
Name: %s
Identifier: %s
CreatedOn: %s
UpdatedOn: %s

%s
`[1:],
		project.Id,
		project.Name,
		project.Identifier,
		project.CreatedOn,
		project.UpdatedOn,
		project.Description)
}

func listProjects() {
	config := getConfig()
	c := redmine.NewClient(config.Endpoint, config.Apikey)
	issues, err := c.Projects()
	if err != nil {
		fatal("Failed to list projects: %s\n", err)
	}
	for _, i := range issues {
		fmt.Printf("%4d: %s\n", i.Id, i.Name)
	}
}

func usage() {
	fmt.Println(`gotmine <command> <subcommand> [arguments]

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
`)
	os.Exit(1)
}

type config struct {
	Endpoint string `json:"endpoint"`
	Apikey   string `json:"apikey"`
	Project  int    `json:"project"`
}

func toUtf8(s string) string {
	ic, err := iconv.Open("char", "UTF-8")
	if err != nil {
		return s
	}
	defer ic.Close()
	ret, _ := ic.Conv(s)
	return ret
}

func main() {
	if len(os.Args) <= 2 {
		usage()
	}

	switch os.Args[1] {
	case "i", "issue":
		switch os.Args[2] {
		case "a", "add":
			createIssue()
			break
		case "c", "create":
			if len(os.Args) == 5 {
				addIssue(os.Args[3], os.Args[4])
			} else {
				usage()
			}
			break
		case "u", "update":
			if len(os.Args) == 4 {
				id, err := strconv.Atoi(os.Args[3])
				if err != nil {
					fatal("Invalid issue id: %s\n", err)
				}
				updateIssue(id)
			} else {
				usage()
			}
			break
		case "d", "delete":
			if len(os.Args) == 4 {
				id, err := strconv.Atoi(os.Args[3])
				if err != nil {
					fatal("Invalid issue id: %s\n", err)
				}
				deleteIssue(id)
			} else {
				usage()
			}
			break
		case "n", "notes":
			if len(os.Args) == 4 {
				id, err := strconv.Atoi(os.Args[3])
				if err != nil {
					fatal("Invalid issue id: %s\n", err)
				}
				notesIssue(id)
			} else {
				usage()
			}
			break
		case "s", "show":
			if len(os.Args) == 4 {
				id, err := strconv.Atoi(os.Args[3])
				if err != nil {
					fatal("Invalid issue id: %s\n", err)
				}
				showIssue(id)
			} else {
				usage()
			}
			break
		case "x", "close":
			if len(os.Args) == 4 {
				id, err := strconv.Atoi(os.Args[3])
				if err != nil {
					fatal("Invalid issue id: %s\n", err)
				}
				closeIssue(id)
			} else {
				usage()
			}
			break
		case "l", "list":
			listIssues()
			break
		default:
			usage()
		}
	case "p", "project":
		switch os.Args[2] {
		case "a", "add":
			createProject()
			break
		case "c", "create":
			if len(os.Args) == 6 {
				addProject(os.Args[3], os.Args[4], os.Args[5])
			} else {
				usage()
			}
			break
		case "s", "show":
			if len(os.Args) == 4 {
				id, err := strconv.Atoi(os.Args[3])
				if err != nil {
					fatal("Invalid project id: %s\n", err)
				}
				showProject(id)
			} else {
				usage()
			}
			break
		case "d", "delete":
			if len(os.Args) == 4 {
				id, err := strconv.Atoi(os.Args[3])
				if err != nil {
					fatal("Invalid project id: %s\n", err)
				}
				deleteProject(id)
			} else {
				usage()
			}
			break
		case "l", "list":
			listProjects()
			break
		default:
			usage()
		}
	default:
		usage()
	}
}
