package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/xanzy/go-gitlab"
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

	token := os.Getenv("GPTSCRIPT_GITLAB_TOKEN")
	if token == "" {
		logrus.Errorf("GitLab token is required")
		os.Exit(1)
	}

	git, err := gitlab.NewClient(token)
	if err != nil {
		logrus.Fatalf("Failed to create client: %v", err)
	}

	issues, err := search(ctx, git, a)
	if err != nil {
		logrus.Errorf("error listing issues: %v", err)
		os.Exit(1)
	}

	printIssues(issues)
}

func search(ctx context.Context, git *gitlab.Client, a args) ([]*gitlab.Issue, error) {
	page, err := strconv.Atoi(a.Page)
	if err != nil {
		return nil, fmt.Errorf("error parsing page number: %v", err)
	}

	searchOpts := &gitlab.ListProjectIssuesOptions{
		ListOptions: gitlab.ListOptions{Page: page, PerPage: 20},
		Search:      &a.Query,
	}

	// Replace PROJECT_ID with your project ID
	issues, _, err := git.Issues.ListProjectIssues("PROJECT_ID", searchOpts)
	if err != nil {
		return nil, err
	}
	return issues, nil
}

func printIssues(issues []*gitlab.Issue) {
	if len(issues) == 0 {
		fmt.Println("No issues found")
		return
	}

	for _, issue := range issues {
		printIssue(issue)
	}
}

func printIssue(issue *gitlab.Issue) {
	fmt.Printf("Title: %s\n", issue.Title)
	fmt.Printf("URL: %s\n", issue.WebURL)
	fmt.Printf("Description: %s\n\n---\n\n", trunc(issue.Description))
}

func trunc(s string) string {
	if len(s) > 500 {
		return s[:500] + "..."
	}
	return s
}
