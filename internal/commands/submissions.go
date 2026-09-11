package commands

import (
	"fmt"

	"github.com/neetozone/neeto-cli-commons/output"
	"github.com/spf13/cobra"
)

var submissionsCmd = &cobra.Command{
	Use:   "submissions",
	Short: "View form submissions",
}

var submissionsListCmd = &cobra.Command{
	Use:   "list <form-id>",
	Short: "List submissions for a form",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		data, err := c.Get(fmt.Sprintf("/forms/%s/submissions", args[0]), paginationParams(cmd))
		if err != nil {
			return err
		}

		printList(data, "submissions", []output.Breadcrumb{
			{Label: "Form", Command: "neetoform forms list"},
		})
		return nil
	},
}

func init() {
	addPaginationFlags(submissionsListCmd)

	submissionsCmd.AddCommand(submissionsListCmd)
	formsCmd.AddCommand(submissionsCmd)
}
