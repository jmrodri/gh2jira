package list

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "list github issues",
		Long:  `This command blah blah blah`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			fmt.Println("LIST COMMAND")
			return nil
		},
	}

}
