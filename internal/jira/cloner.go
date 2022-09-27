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

package jira

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/andygrunwald/go-jira"
	gojira "github.com/andygrunwald/go-jira"
	"github.com/google/go-github/v47/github"
)

type Option func(*ClonerConfig) error

type ClonerConfig struct {
	client  *http.Client
	dryRun  bool
	project string
	jiraURL string
}

func (c *ClonerConfig) setDefaults() error {
	if c.client == nil {
		token, err := c.getToken()
		if err != nil {
			return err
		}

		tp := gojira.BearerAuthTransport{
			Token: token,
		}
		c.client = tp.Client()
	}
	if c.jiraURL == "" {
		c.jiraURL = "https://issues.redhat.com"
	}
	return nil
}

func (c *ClonerConfig) getToken() (string, error) {
	token, ok := os.LookupEnv("JIRA_TOKEN")
	if !ok {
		return "", fmt.Errorf("please supply your JIRA_TOKEN")
	}
	return token, nil
}

func WithClient(cl *http.Client) Option {
	return func(c *ClonerConfig) error {
		c.client = cl
		return nil
	}
}

func WithDryRun(dr bool) Option {
	return func(c *ClonerConfig) error {
		c.dryRun = dr
		return nil
	}
}

func WithProject(p string) Option {
	return func(c *ClonerConfig) error {
		c.project = p
		return nil
	}
}

func WithJiraURL(j string) Option {
	return func(c *ClonerConfig) error {
		c.jiraURL = j
		return nil
	}
}

func getWebURL(url string) string {
	// https://api.github.com/repos/operator-framework/operator-sdk/issues/3447
	// https://github.com/operator-framework/operator-sdk/issues/3447
	if url == "" {
		return url
	}
	return strings.Replace(strings.Replace(url, "api.github.com", "github.com", 1), "repos/", "", 1)
}

func Clone(issue *github.Issue, opts ...Option) error {
	config := ClonerConfig{}
	for _, opt := range opts {
		if err := opt(&config); err != nil {
			return err
		}
	}

	jiraClient, err := gojira.NewClient(config.client, config.jiraURL)
	if err != nil {
		return err
	}

	ji := jira.Issue{
		Fields: &gojira.IssueFields{
			// Assignee: &gojira.User{
			//     Name: "myuser",
			// },
			// Reporter: &gojira.User{
			//     Name: "youruser",
			// },
			Description: fmt.Sprintf("%s\n\nUpstream Github issue: %s\n", issue.GetBody(), getWebURL(issue.GetURL())),
			Type: gojira.IssueType{
				Name: "Story",
			},
			Project: gojira.Project{
				Key: config.project,
			},
			Summary: fmt.Sprintf("[UPSTREAM] %s #%d", issue.GetTitle(), issue.GetNumber()),
		},
	}

	if config.dryRun {
		fmt.Println("\n############# DRY RUN MODE #############")
		fmt.Printf("Cloning issue #%d to jira project board: %s\n\n", issue.GetNumber(), ji.Fields.Project.Key)
		fmt.Printf("Summary: %s\n", ji.Fields.Summary)
		fmt.Printf("Type: %s\n", ji.Fields.Type.Name)
		fmt.Println("Description:")
		fmt.Printf("%s\n", ji.Fields.Description)
		fmt.Println("\n############# DRY RUN MODE #############")
	} else {
		fmt.Printf("Cloning issue #%d to jira project board: %s\n\n", issue.GetNumber(), ji.Fields.Project.Key)
		daIssue, _, err := jiraClient.Issue.Create(&ji)
		if err != nil {
			return err
		}

		if daIssue != nil {
			fmt.Printf("Issue cloned; see %s\n",
				fmt.Sprintf("https://issues.redhat.com/browse/%s", daIssue.Key))
		}
	}

	return nil
}
