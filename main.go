package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v47/github"
	"golang.org/x/oauth2"
)

// So we will want to allow this to be able to take in a specific GH issue id or
// --all.

// Login to GH and get the Issue.
// Then login to jira and create a new issue (attach to an epic if supplied)
//

// Needs a --dryrun flag which will print out what jira issue it will create

// gh2jira list operator-framework operator-sdk
// gh2jira copy GH# [--dry-run]
func main() {
	GetGithubIssues()
	CreateJiraIssue(nil)
}

func CreateJiraIssue(issue *github.Issue) {
}

func GetGithubIssues() {
	token, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		fmt.Println("please supply your GITHUB_TOKEN")
		os.Exit(1)
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	opt := &github.IssueListByRepoOptions{
		ListOptions: github.ListOptions{PerPage: 50},
		State:       "open",
	}

	var allIssues []*github.Issue

	for {
		issues, resp, err := client.Issues.ListByRepo(context.Background(), "operator-framework", "operator-sdk", opt)
		if err != nil {
			fmt.Println(err.Error())
		}
		// fmt.Println(len(issues))

		allIssues = append(allIssues, issues...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	// fmt.Println(len(allIssues))
	for _, issue := range allIssues {
		if issue.IsPullRequest() {
			// We have a PR, skipping
			continue
		}
		// fmt.Println(*issue.ID)
		fmt.Println(issue.GetNumber())
		// fmt.Println(*issue.Title)
		fmt.Println(issue.GetState())
		if issue.GetAssignee() != nil {
			fmt.Println(issue.GetAssignee().Login)
		}

		// NOTE: This should be the jira body
		fmt.Println(issue.GetBody())

		// Look through the labels
		// important soon should be Urgent
		// important long term should be Medium
		// fmt.Println(issue.Labels)
		fmt.Println("================")
	}
}
