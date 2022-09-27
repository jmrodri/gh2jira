package clone

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jmrodri/gh2jira/internal/gh"
	"github.com/jmrodri/gh2jira/internal/jira"
)

var dryRun bool
var project string
var ghproject string

func NewCmd() *cobra.Command {
	cloner := jira.Cloner{}
	lister := gh.Lister{}

	cmd := &cobra.Command{
		Use:   "clone <ISSUE_ID> [ISSUE_ID ...]",
		Short: "Clone given Github issues to Jira",
		Long:  "Clone given Github issues to Jira. WARNING! This will write to your jira instance. Use --dryrun to see what will happen",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, id := range args {
				issueId, _ := strconv.Atoi(id)
				lister.Options = &gh.ListerOptions{
					Project: ghproject,
				}
				issue := lister.GetIssue(issueId)
				cloner.Clone(issue, project, dryRun)
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
