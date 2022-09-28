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
	"fmt"

	"github.com/google/go-github/v47/github"
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
