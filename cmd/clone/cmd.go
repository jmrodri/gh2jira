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

package clone

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jmrodri/gh2jira/internal/gh"
	"github.com/jmrodri/gh2jira/internal/jira"
)

var (
	dryRun    bool
	project   string
	ghproject string
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clone <ISSUE_ID> [ISSUE_ID ...]",
		Short: "Clone given Github issues to Jira",
		Long:  "Clone given Github issues to Jira. WARNING! This will write to your jira instance. Use --dryrun to see what will happen",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, id := range args {
				issueId, _ := strconv.Atoi(id)
				issue, err := gh.GetIssue(issueId, gh.WithProject(ghproject))
				if err != nil {
					return err
				}
				_, err = jira.Clone(issue, jira.WithProject(project), jira.WithDryRun(dryRun))
				if err != nil {
					return nil
				}
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dryrun", false, "display what we would do without cloning")
	cmd.Flags().StringVar(&project, "project", "OSDK", "Jira project to clone to")
	cmd.Flags().StringVar(&ghproject, "github-project", "operator-framework/operator-sdk",
		"Github project to clone from e.g. ORG/REPO")

	return cmd
}
