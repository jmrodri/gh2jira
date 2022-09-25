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
package main

import (
	"os"

	// "github.com/jmrodri/gh2jira/internal/gh"
	// "github.com/jmrodri/gh2jira/internal/jira"
	"github.com/jmrodri/gh2jira/cmd/root"
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
	cmd := root.NewCmd()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
	// if len(os.Args) < 2 {
	//     fmt.Println("Usage: list or clone")
	//     os.Exit(1)
	// }
	//
	// command := os.Args[1]
	//
	// switch command {
	// case "list":
	//     // Pass in options to the List command
	//     gh.ListIssues()
	// case "clone":
	//     // Need github issue number
	//     if len(os.Args) < 3 {
	//         fmt.Println("Usage: clone ISSUE_TO_CLONE")
	//         os.Exit(1)
	//     }
	//     issueNum, _ := strconv.Atoi(os.Args[2])
	//     issue := gh.GetIssue(issueNum)
	//
	//     var dryRun bool
	//     if len(os.Args) > 3 {
	//         dryRun = os.Args[3] == "--dryrun"
	//     }
	//     jira.CloneIssueToJira(issue, dryRun)
	// default:
	//     fmt.Printf("Unsupported command: %s\n", command)
	// }
}
