package commands

import (
	"fmt"
	"slices"
	"strings"

	"github.com/neetozone/neeto-cli-commons/output"
	"github.com/spf13/cobra"
)

var validFormStatuses = []string{"active", "archived", "favorite"}

var formsCmd = &cobra.Command{
	Use:   "forms",
	Short: "Manage forms",
}

func validateFormStatus(status string) error {
	if slices.Contains(validFormStatuses, status) {
		return nil
	}
	return fmt.Errorf("invalid --status %q; valid values: %s", status, strings.Join(validFormStatuses, ", "))
}

var formsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List forms",
	RunE: func(cmd *cobra.Command, args []string) error {
		params := paginationParams(cmd)
		if status, _ := cmd.Flags().GetString("status"); status != "" {
			if err := validateFormStatus(status); err != nil {
				return err
			}
			params.Set("status", status)
		}

		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		data, err := c.Get("/forms", params)
		if err != nil {
			return err
		}

		printList(data, "forms", []output.Breadcrumb{
			{Label: "Submissions", Command: "neetoform forms submissions list <form-id>"},
		})
		return nil
	},
}

func init() {
	addPaginationFlags(formsListCmd)
	formsListCmd.Flags().String("status", "", "Filter by status: active, archived, favorite")

	formsCmd.AddCommand(formsListCmd)
	register(func(root *cobra.Command) { root.AddCommand(formsCmd) })
}
