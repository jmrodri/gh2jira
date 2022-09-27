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

	cmd.Flags().StringVar(&lo.Milestone, "milestone", "", "the milestone ID from the url, not the display name")
	cmd.Flags().StringVar(&lo.Assignee, "assignee", "", "username of the issue is assigned")
	cmd.Flags().StringVar(&lo.Project, "project", "operator-framework/operator-sdk",
		"Github project to list e.g. ORG/REPO")
	cmd.Flags().StringSliceVar(&lo.Label, "label", nil, "label i.e. --label \"documentation,bug\" or --label doc --label bug")

	return cmd
}
