package main

import (
	"github.com/google/go-github/v47/github"

	"github.com/jmrodri/gh2jira/internal/gh"
)

// So we will want to allow this to be able to take in a specific GH issue id or
// --all.

// Login to GH and get the Issue.
// Then login to jira and create a new issue (attach to an epic if supplied)
//

// Needs a --dryrun flag which will print out what jira issue it will create

// gh2jira genconfig
// gh2jira list --project operator-framework/operator-sdk [--milestone=] [--assignee=]
// gh2jira copy GH# [--dry-run]
func main() {
	gh.GetGithubIssues()
	// CreateJiraIssue(nil)
}

func CreateJiraIssue(issue *github.Issue) {
}
