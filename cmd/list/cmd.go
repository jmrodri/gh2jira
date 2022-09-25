package list

import (
	"github.com/jmrodri/gh2jira/internal/gh"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	lo := gh.ListerOptions{}
	lister := gh.Lister{
		Options: &lo,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Github issues",
		Long:  "List Github issues filtered by milestone, assignee, or label",
		RunE: func(cmd *cobra.Command, args []string) error {
			lister.ListIssues()
			return nil
		},
	}

	cmd.Flags().StringVar(&lo.Milestone, "milestone", "", "milestone")
	cmd.Flags().StringVar(&lo.Assignee, "assignee", "", "assignee")
	cmd.Flags().StringVar(&lo.Label, "label", "", "label")

	return cmd
}
