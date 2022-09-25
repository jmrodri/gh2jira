package clone

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clone",
		Short: "Clone given Github issues to Jira",
		Long:  "Clone given Github issues to Jira",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("clone called")
			return nil
		},
	}

	return cmd
}
