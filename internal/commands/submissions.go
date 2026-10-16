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

		params := paginationParams(cmd)
		if after, _ := cmd.Flags().GetString("submitted-after"); after != "" {
			params.Set("submitted_after", after)
		}
		if before, _ := cmd.Flags().GetString("submitted-before"); before != "" {
			params.Set("submitted_before", before)
		}
		if err := setOrderParam(cmd, params); err != nil {
			return err
		}

		data, err := c.Get(fmt.Sprintf("/forms/%s/submissions", args[0]), params)
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
	submissionsListCmd.Flags().String("submitted-after", "", "Only submissions at or after this ISO 8601 date (YYYY-MM-DD) or date-time")
	submissionsListCmd.Flags().String("submitted-before", "", "Only submissions before this ISO 8601 date (YYYY-MM-DD) or date-time")
	submissionsListCmd.Flags().String("order", "", "Sort by submission time: asc (oldest first, default) or desc")

	submissionsCmd.AddCommand(submissionsListCmd)
	formsCmd.AddCommand(submissionsCmd)
}
