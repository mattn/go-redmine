package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/mattn/go-redmine"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type config struct {
	Endpoint string `json:"endpoint"`
	Apikey   string `json:"apikey"`
	Project  int    `json:"project"`
	Editor   string `json:"editor"`
	Insecure bool   `json:"insecure"`
}

var profile *string = flag.String("p", os.Getenv("GODMINE_ENV"), "profile")
var conf config

func fatal(format string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, format, err)
	} else {
		fmt.Fprint(os.Stderr, format)
	}
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

func notesFromEditor(issue *redmine.Issue) (string, error) {
	file := ""
	newf := fmt.Sprintf("%d.txt", rand.Int())
	if runtime.GOOS == "windows" {
		file = filepath.Join(os.Getenv("APPDATA"), "godmine", newf)
	} else {
		file = filepath.Join(os.Getenv("HOME"), ".config", "godmine", newf)
	}
	defer os.Remove(file)
	editor := getEditor()

	body := "### Notes Here ###\n"
	contents := issue.GetTitle() + "\n" + body

	ioutil.WriteFile(file, []byte(contents), 0600)

	if err := run([]string{editor, file}); err != nil {
		return "", err
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	text := strings.Join(strings.SplitN(string(b), "\n", 2)[1:], "\n")

	if text == body {
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
	editor := getEditor()

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
	editor := getEditor()

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

func getEditor() string {
	editor := conf.Editor
	if editor == "" {
		editor = os.Getenv("EDITOR")
		if editor == "" {
			if runtime.GOOS == "windows" {
				editor = "notepad"
			} else {
				editor = "vim"
			}
		}
	}
	return editor
}
func getConfig() config {
	file := "settings.json"

	if *profile != "" {
		file = "settings." + *profile + ".json"
	}

	if runtime.GOOS == "windows" {
		file = filepath.Join(os.Getenv("APPDATA"), "godmine", file)
	} else {
		file = filepath.Join(os.Getenv("HOME"), ".config", "godmine", file)
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
	var issue redmine.Issue
	c := redmine.NewClient(conf.Endpoint, conf.Apikey)
	issue.ProjectId = conf.Project
	issue.Subject = subject
	issue.Description = description
	_, err := c.CreateIssue(issue)
	if err != nil {
		fatal("Failed to create issue: %s\n", err)
	}
}

func createIssue() {
	issue, err := issueFromEditor("")
	if err != nil {
		fatal("%s\n", err)
	}
	c := redmine.NewClient(conf.Endpoint, conf.Apikey)
	issue.ProjectId = conf.Project
	_, err = c.CreateIssue(*issue)
	if err != nil {
		fatal("Failed to create issue: %s\n", err)
	}
}

func updateIssue(id int) {
	c := redmine.NewClient(conf.Endpoint, conf.Apikey)
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
	issue.ProjectId = conf.Project
	err = c.UpdateIssue(*issue)
	if err != nil {
		fatal("Failed to update issue: %s\n", err)
	}
}

func deleteIssue(id int) {
	c := redmine.NewClient(conf.Endpoint, conf.Apikey)
	err := c.DeleteIssue(id)
	if err != nil {
		fatal("Failed to delete issue: %s\n", err)
	}
}

func closeIssue(id int) {
	c := redmine.NewClient(conf.Endpoint, conf.Apikey)
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
	c := redmine.NewClient(conf.Endpoint, conf.Apikey)
	issue, err := c.Issue(id)
	if err != nil {
		fatal("Failed to update issue: %s\n", err)
	}

	content, err := notesFromEditor(issue)
	if err != nil {
		fatal("%s\n", err)
	}
	issue.Notes = content
	issue.ProjectId = conf.Project
	err = c.UpdateIssue(*issue)
	if err != nil {
		fatal("Failed to update issue: %s\n", err)
	}
}

func showIssue(id int) {
	c := redmine.NewClient(conf.Endpoint, conf.Apikey)
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
	c := redmine.NewClient(conf.Endpoint, conf.Apikey)
	issues, err := c.Issues()
	if err != nil {
		fatal("Failed to list issues: %s\n", err)
	}
	for _, i := range issues {
		fmt.Printf("%4d: %s\n", i.Id, i.Subject)
	}
}

func addProject(name, identifier, description string) {
	var project redmine.Project
	c := redmine.NewClient(conf.Endpoint, conf.Apikey)
	project.Name = name
	project.Identifier = identifier
	project.Description = description
	_, err := c.CreateProject(project)
	if err != nil {
		fatal("Failed to create project: %s\n", err)
	}
}

func createProject() {
	project, err := projectFromEditor("")
	if err != nil {
		fatal("%s\n", err)
	}
	c := redmine.NewClient(conf.Endpoint, conf.Apikey)
	_, err = c.CreateProject(*project)
	if err != nil {
		fatal("Failed to create project: %s\n", err)
	}
}

func updateProject(id int) {
	c := redmine.NewClient(conf.Endpoint, conf.Apikey)
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
	c := redmine.NewClient(conf.Endpoint, conf.Apikey)
	err := c.DeleteProject(id)
	if err != nil {
		fatal("Failed to delete project: %s\n", err)
	}
}

func showProject(id int) {
	c := redmine.NewClient(conf.Endpoint, conf.Apikey)
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
	c := redmine.NewClient(conf.Endpoint, conf.Apikey)
	issues, err := c.Projects()
	if err != nil {
		fatal("Failed to list projects: %s\n", err)
	}
	for _, i := range issues {
		fmt.Printf("%4d: %s\n", i.Id, i.Name)
	}
}

func showMembership(id int) {
	c := redmine.NewClient(conf.Endpoint, conf.Apikey)
	membership, err := c.Membership(id)
	if err != nil {
		fatal("Failed to show membership: %s\n", err)
	}

	fmt.Printf(`
Id: %d
Project: %s
User: %s
Role: `[1:],
		membership.Id,
		membership.Project.Name,
		membership.User.Name)
	for i, role := range membership.Roles {
		if i != 0 {
			fmt.Print(", ")
		}
		fmt.Printf(role.Name)
	}
	fmt.Println()
}

func listMemberships(projectId int) {
	c := redmine.NewClient(conf.Endpoint, conf.Apikey)
	memberships, err := c.Memberships(projectId)
	if err != nil {
		fatal("Failed to list memberships: %s\n", err)
	}
	for _, i := range memberships {
		fmt.Printf("%4d: %s\n", i.Id, i.User.Name)
	}
}

func showUser(id int) {
	c := redmine.NewClient(conf.Endpoint, conf.Apikey)
	user, err := c.User(id)
	if err != nil {
		fatal("Failed to show user: %s\n", err)
	}

	fmt.Printf(`
Id: %d
Login: %s
Firstname: %s
Lastname: %s
Mail: %s
CreatedOn: %s
`[1:],
		user.Id,
		user.Login,
		user.Firstname,
		user.Lastname,
		user.Mail,
		user.CreatedOn)
}

func listUsers() {
	c := redmine.NewClient(conf.Endpoint, conf.Apikey)
	users, err := c.Users()
	if err != nil {
		fatal("Failed to list users: %s\n", err)
	}
	for _, i := range users {
		fmt.Printf("%4d: %s\n", i.Id, i.Login)
	}
}

func showNews(id int) {
	c := redmine.NewClient(conf.Endpoint, conf.Apikey)
	news, err := c.News(id)
	if err != nil {
		fatal("Failed to show user: %s\n", err)
	}

	found := -1
	for i, n := range news {
		if n.Id == id {
			found = i
		}
	}
	if found != -1 {
		fmt.Printf(`
Id: %d
Project: %s
Title: %s
Summary: %s
CreatedOn: %s

%s
`[1:],
			news[found].Id,
			news[found].Project.Name,
			news[found].Title,
			news[found].Summary,
			news[found].CreatedOn,
			news[found].Description)
	} else {
		fatal("Failed to show news: not found\n", nil)
	}
}

func listNews() {
	c := redmine.NewClient(conf.Endpoint, conf.Apikey)
	news, err := c.News(conf.Project)
	if err != nil {
		fatal("Failed to list users: %s\n", err)
	}
	for _, i := range news {
		fmt.Printf("%4d: %s\n", i.Id, i.Title)
	}
}

func usage() {
	fmt.Println(`godmine <command> <subcommand> [arguments]

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

Membership Commands:
  show     s show given membership.
             $ godmine m s 1

  list     l listing memberships of given project.
             $ godmine m l 1

User Commands:
  show     s show given user.
             $ godmine u s 1

  list     l listing users.
             $ godmine u l
`)
	os.Exit(1)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	flag.Parse()
	if flag.NArg() <= 1 {
		usage()
	}
	conf = getConfig()
	if conf.Insecure {
		http.DefaultClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	}

	switch flag.Arg(0) {
	case "i", "issue":
		switch flag.Arg(1) {
		case "a", "add":
			createIssue()
			break
		case "c", "create":
			if flag.NArg() == 4 {
				addIssue(flag.Arg(2), flag.Arg(3))
			} else {
				usage()
			}
			break
		case "u", "update":
			if flag.NArg() == 3 {
				id, err := strconv.Atoi(flag.Arg(2))
				if err != nil {
					fatal("Invalid issue id: %s\n", err)
				}
				updateIssue(id)
			} else {
				usage()
			}
			break
		case "d", "delete":
			if flag.NArg() == 3 {
				id, err := strconv.Atoi(flag.Arg(2))
				if err != nil {
					fatal("Invalid issue id: %s\n", err)
				}
				deleteIssue(id)
			} else {
				usage()
			}
			break
		case "n", "notes":
			if flag.NArg() == 3 {
				id, err := strconv.Atoi(flag.Arg(2))
				if err != nil {
					fatal("Invalid issue id: %s\n", err)
				}
				notesIssue(id)
			} else {
				usage()
			}
			break
		case "s", "show":
			if flag.NArg() == 3 {
				id, err := strconv.Atoi(flag.Arg(2))
				if err != nil {
					fatal("Invalid issue id: %s\n", err)
				}
				showIssue(id)
			} else {
				usage()
			}
			break
		case "x", "close":
			if flag.NArg() == 3 {
				id, err := strconv.Atoi(flag.Arg(2))
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
		switch flag.Arg(1) {
		case "a", "add":
			createProject()
			break
		case "c", "create":
			if flag.NArg() == 5 {
				addProject(flag.Arg(2), flag.Arg(3), flag.Arg(4))
			} else {
				usage()
			}
			break
		case "s", "show":
			if flag.NArg() == 3 {
				id, err := strconv.Atoi(flag.Arg(2))
				if err != nil {
					fatal("Invalid project id: %s\n", err)
				}
				showProject(id)
			} else {
				usage()
			}
			break
		case "d", "delete":
			if flag.NArg() == 3 {
				id, err := strconv.Atoi(flag.Arg(2))
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
	case "m", "membership":
		switch flag.Arg(1) {
		case "s", "show":
			if flag.NArg() == 3 {
				id, err := strconv.Atoi(flag.Arg(2))
				if err != nil {
					fatal("Invalid membership id: %s\n", err)
				}
				showMembership(id)
			} else {
				usage()
			}
			break
		case "l", "list":
			if flag.NArg() == 3 {
				projectId, err := strconv.Atoi(flag.Arg(2))
				if err != nil {
					fatal("Invalid project id: %s\n", err)
				}
				listMemberships(projectId)
			} else {
				usage()
			}
			break
		default:
			usage()
		}
	case "u", "user":
		switch flag.Arg(1) {
		case "s", "show":
			if flag.NArg() == 3 {
				id, err := strconv.Atoi(flag.Arg(2))
				if err != nil {
					fatal("Invalid user id: %s\n", err)
				}
				showUser(id)
			} else {
				usage()
			}
			break
		case "l", "list":
			listUsers()
			break
		default:
			usage()
		}
	case "n", "news":
		switch flag.Arg(1) {
		case "s", "show":
			if flag.NArg() == 3 {
				id, err := strconv.Atoi(flag.Arg(2))
				if err != nil {
					fatal("Invalid project id: %s\n", err)
				}
				showNews(id)
			} else {
				usage()
			}
			break
		case "l", "list":
			listNews()
			break
		default:
			usage()
		}
	default:
		usage()
	}
}
