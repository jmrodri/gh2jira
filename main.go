package main

import (
	"fmt"
	"os"

	"github.com/google/go-github/v47/github"

	"github.com/jmrodri/gh2jira/internal/gh"
)

// So we will want to allow this to be able to take in a specific GH issue id or
// --all.

// Login to GH and get the Issue.
// Then login to jira and create a new issue (attach to an epic if supplied)
//

// Needs a --dryrun flag which will print out what jira issue it will create

// global flags
// --no-color
// --oneline
// gh2jira genconfig
// gh2jira list --project operator-framework/operator-sdk [--milestone=] [--assignee=]
// gh2jira clone GH# [--dry-run]
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: list or clone")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "list":
		gh.GetGithubIssues()
	case "clone":
		fmt.Println("clone command is TBD")
	default:
		fmt.Printf("Unsupported command: %s\n", command)
	}
	// CreateJiraIssue(nil)
}

func CreateJiraIssue(issue *github.Issue) {
}
