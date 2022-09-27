/*
Copyright © 2022 jesus m. rodriguez jmrodri@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package jira

import (
	"fmt"
	"os"
	"strings"

	"github.com/andygrunwald/go-jira"
	gojira "github.com/andygrunwald/go-jira"
	"github.com/google/go-github/v47/github"
)

type Cloner struct {
}

func getToken() string {
	token, ok := os.LookupEnv("JIRA_TOKEN")
	if !ok {
		fmt.Println("please supply your JIRA_TOKEN")
		os.Exit(1)
	}
	return token
}

func getWebURL(url string) string {
	// https://api.github.com/repos/operator-framework/operator-sdk/issues/3447
	// https://github.com/operator-framework/operator-sdk/issues/3447
	if url == "" {
		return url
	}
	return strings.Replace(strings.Replace(url, "api.github.com", "github.com", 1), "repos/", "", 1)
}

func (c *Cloner) Clone(issue *github.Issue, project string, dryRun bool) {
	token := getToken()

	tp := gojira.BearerAuthTransport{
		Token: token,
	}

	jiraClient, err := gojira.NewClient(tp.Client(), "https://issues.redhat.com")
	if err != nil {
		panic(err)
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
				Key: project,
			},
			Summary: fmt.Sprintf("[UPSTREAM] %s #%d", issue.GetTitle(), issue.GetNumber()),
		},
	}

	if dryRun {
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
			fmt.Println(err)
			os.Exit(2)
		}

		if daIssue != nil {
			fmt.Printf("Issue cloned; see %s\n",
				fmt.Sprintf("https://issues.redhat.com/browse/%s", daIssue.Key))
		}
	}
}
