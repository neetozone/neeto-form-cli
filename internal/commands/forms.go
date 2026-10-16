package commands

import (
	"fmt"
	"net/url"
	"slices"
	"strings"

	"github.com/neetozone/neeto-cli-commons/output"
	"github.com/spf13/cobra"
)

var validFormStatuses = []string{"active", "archived", "favorite"}

var validOrders = []string{"asc", "desc"}

var formsCmd = &cobra.Command{
	Use:   "forms",
	Short: "Manage forms",
}

func validateFormStatus(status string) error {
	if slices.Contains(validFormStatuses, status) {
		return nil
	}
	return fmt.Errorf("Invalid --status %q. Valid values: %s.", status, strings.Join(validFormStatuses, ", "))
}

func validateOrder(order string) error {
	if slices.Contains(validOrders, order) {
		return nil
	}
	return fmt.Errorf("Invalid --order %q. Valid values: %s.", order, strings.Join(validOrders, ", "))
}

func setOrderParam(cmd *cobra.Command, params url.Values) error {
	order, _ := cmd.Flags().GetString("order")
	order = strings.ToLower(order)
	if order == "" {
		return nil
	}
	if err := validateOrder(order); err != nil {
		return err
	}
	params.Set("order", order)
	return nil
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
		if search, _ := cmd.Flags().GetString("search"); search != "" {
			params.Set("search", search)
		}
		if err := setOrderParam(cmd, params); err != nil {
			return err
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
			{Label: "Show", Command: "neetoform forms show <form-id>"},
			{Label: "Submissions", Command: "neetoform forms submissions list <form-id>"},
		})
		return nil
	},
}

var formsShowCmd = &cobra.Command{
	Use:   "show <form-id>",
	Short: "Show a form and its questions",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		data, err := c.Get(fmt.Sprintf("/forms/%s", url.PathEscape(args[0])), nil)
		if err != nil {
			return err
		}

		printResource(data, []output.Breadcrumb{
			{Label: "Submissions", Command: "neetoform forms submissions list <form-id>"},
		})
		return nil
	},
}

func init() {
	addPaginationFlags(formsListCmd)
	formsListCmd.Flags().String("status", "", "Filter by status: active, archived, favorite")

	formsListCmd.Flags().String("search", "", "Words that must all appear in the form title")
	formsListCmd.Flags().String("order", "", "Sort by creation time: asc (oldest first, default) or desc")

	formsCmd.AddCommand(formsListCmd)
	formsCmd.AddCommand(formsShowCmd)
	register(func(root *cobra.Command) { root.AddCommand(formsCmd) })
}
