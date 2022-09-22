package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jmrodri/gh2jira/internal/gh"
	"github.com/jmrodri/gh2jira/internal/jira"
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
		// Pass in options to the List command
		gh.ListIssues()
	case "clone":
		// Need github issue number
		if len(os.Args) < 3 {
			fmt.Println("Usage: clone ISSUE_TO_CLONE")
			os.Exit(1)
		}
		issueNum, _ := strconv.Atoi(os.Args[2])
		issue := gh.GetIssue(issueNum)

		var dryRun bool
		if len(os.Args) > 3 {
			dryRun = os.Args[3] == "--dryrun"
		}
		jira.CloneIssueToJira(issue, dryRun)
	default:
		fmt.Printf("Unsupported command: %s\n", command)
	}
}
