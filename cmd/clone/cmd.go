package clone

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jmrodri/gh2jira/internal/gh"
	"github.com/jmrodri/gh2jira/internal/jira"
)

var dryRun bool

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
				issue := lister.GetIssue(issueId)
				cloner.Clone(issue, dryRun)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dryrun", false, "display what we would do without cloning")

	return cmd
}
