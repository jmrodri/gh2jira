package list

import (
	"fmt"

	"github.com/jmrodri/gh2jira/internal/gh"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	lister := gh.Lister{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Github issues",
		Long:  "List Github issues filtered by milestone, assignee, or label",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("list called")
			fmt.Println(len(args))
			lister.ListIssues()
			return nil
		},
	}

	l := gh.ListerOptions{}

	cmd.Flags().StringVar(&l.Milestone, "milestone", "", "milestone")
	cmd.Flags().StringVar(&l.Assignee, "assignee", "", "assignee")
	cmd.Flags().StringVar(&l.Label, "label", "", "label")

	return cmd
}
