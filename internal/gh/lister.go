// Copyright © 2022 jesus m. rodriguez jmrodri@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gh

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v47/github"
	"golang.org/x/oauth2"
)

type ListerOptions struct {
	Milestone string
	Assignee  string
	Project   string
	Label     []string
}

type Lister struct {
	Options *ListerOptions
}

func (l *ListerOptions) GetGithubOrg() string {
	return strings.Split(l.Project, "/")[0]
}

func (l *ListerOptions) GetGithubRepo() string {
	s := strings.Split(l.Project, "/")
	if len(s) == 1 {
		return s[0]
	}
	return s[1]
}

// So we will want to allow this to be able to take in a specific GH issue id or
// --all.

// Login to GH and get the Issue.
// Then login to jira and create a new issue (attach to an epic if supplied)
//

// Needs a --dryrun flag which will print out what jira issue it will create

// gh2jira genconfig
// gh2jira list --project operator-framework/operator-sdk [--milestone=] [--assignee=]
// gh2jira copy GH# [--dry-run]

func PrintGithubIssue(issue *github.Issue, oneline bool, color bool) {

	// fmt.Printf("%5d %s %+v\n", issue.GetNumber(), issue.GetTitle(), issue.GetMilestone())
	// return

	if oneline {
		if color {
			// print the idea in yellow, then reset the rest of the line
			fmt.Printf("\033[33m%5d\033[0m \033[32m%s\033[0m %s\n", issue.GetNumber(), issue.GetState(), issue.GetTitle())
		} else {
			fmt.Printf("%5d %s %s\n", issue.GetNumber(), issue.GetState(), issue.GetTitle())
		}
	} else {
		// fmt.Println(*issue.ID)
		fmt.Printf("Issue:\t%d\n", issue.GetNumber())
		// fmt.Println(*issue.Title)
		fmt.Printf("State:\t%s\n", issue.GetState())
		if issue.GetAssignee() != nil {
			fmt.Printf("Assignee:\t%s\n", *issue.GetAssignee().Login)
		}

		// NOTE: This should be the jira body
		// fmt.Printf("Title:\t%s\n", issue.GetTitle())
		fmt.Printf("\n   %s\n\n", issue.GetTitle())
		// fmt.Printf("Body:\n\t%s\n", issue.GetBody())

		// Look through the labels
		// important soon should be Urgent
		// important long term should be Medium
		// fmt.Println(issue.Labels)
	}
}

func (l *Lister) GetIssue(issueNum int) (*github.Issue, error) {
	token, err := getToken()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	issue, _, err := client.Issues.Get(context.Background(), l.Options.GetGithubOrg(),
		l.Options.GetGithubRepo(), issueNum)

	if err != nil {
		fmt.Println(err.Error())
	}
	return issue, nil
}

func getToken() (string, error) {
	token, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		return "", fmt.Errorf("please supply your GITHUB_TOKEN")
	}
	return token, nil
}

func (l *Lister) ListIssues() error {
	// If no options, at least make it non-nil to avoid any issues later.
	if l.Options == nil {
		l.Options = &ListerOptions{}
	}

	token, err := getToken()
	if err != nil {
		return err
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
		Milestone:   l.Options.Milestone,
		Assignee:    l.Options.Assignee,
		Labels:      l.Options.Label,
	}

	var allIssues []*github.Issue

	for {
		issues, resp, err := client.Issues.ListByRepo(context.Background(),
			l.Options.GetGithubOrg(), l.Options.GetGithubRepo(), opt)

		if err != nil {
			fmt.Println(err.Error())
		}

		allIssues = append(allIssues, issues...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	for _, issue := range allIssues {
		if issue.IsPullRequest() {
			// We have a PR, skipping
			continue
		}
		PrintGithubIssue(issue, true, true)
	}

	return nil
}
