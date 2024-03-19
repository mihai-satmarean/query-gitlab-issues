package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/acorn-io/z"
	"github.com/google/go-github/v60/github"
	"github.com/sirupsen/logrus"
)

type args struct {
	Query string `json:"query"`
	Page  string `json:"page"`
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if len(os.Args) != 2 {
		logrus.Errorf("Usage: %s <JSON parameters>", os.Args[0])
		os.Exit(1)
	}

	var a args
	if err := json.Unmarshal([]byte(os.Args[1]), &a); err != nil {
		logrus.Errorf("error parsing JSON parameters: %v", err)
		os.Exit(1)
	}

	if a.Page == "" {
		a.Page = "1"
	}

	gh := github.NewClient(nil)
	if os.Getenv("GPTSCRIPT_GITHUB_TOKEN") != "" {
		gh = gh.WithAuthToken(os.Getenv("GPTSCRIPT_GITHUB_TOKEN"))
	}

	issues, err := search(ctx, gh, a)
	if err != nil {
		logrus.Errorf("error listing issues: %v", err)
		os.Exit(1)
	}

	printIssues(issues)
}

func search(ctx context.Context, gh *github.Client, a args) ([]*github.Issue, error) {
	page, err := strconv.Atoi(a.Page)
	if err != nil {
		return nil, fmt.Errorf("error parsing page number: %v", err)
	}

	searchOpts := &github.SearchOptions{
		ListOptions: github.ListOptions{Page: page},
	}

	issues, _, err := gh.Search.Issues(ctx, a.Query, searchOpts)
	if issues != nil {
		return issues.Issues, err
	}
	return nil, err
}

func printIssues(issues []*github.Issue) {
	if len(issues) == 0 {
		fmt.Println("No issues found")
		return
	}

	for _, issue := range issues {
		printIssue(issue)
	}
}

func printIssue(issue *github.Issue) {
	if issue.Title == nil {
		issue.Title = z.Pointer("(no title)")
	}
	fmt.Printf("Title: %s\n", *issue.Title)

	if issue.HTMLURL != nil {
		fmt.Printf("URL: %s\n", *issue.HTMLURL)
	}

	fmt.Printf("Description: %s\n\n---\n\n", trunc(z.Dereference(issue.Body)))
}

func trunc(s string) string {
	if len(s) > 500 {
		return s[:500] + "..."
	}
	return s
}
