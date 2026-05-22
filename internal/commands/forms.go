package commands

import (
	"net/url"

	"github.com/neetozone/neeto-form-cli/internal/output"
	"github.com/spf13/cobra"
)

var formsCmd = &cobra.Command{
	Use:   "forms",
	Short: "Manage forms",
}

var formsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List forms",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		params := paginationParams(cmd)
		if status, _ := cmd.Flags().GetString("status"); status != "" {
			params.Set("status", status)
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
	rootCmd.AddCommand(formsCmd)
}

var _ = url.Values{}
