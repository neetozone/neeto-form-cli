package commands

import (
	"net/http"

	"github.com/spf13/cobra"
)

var validFormStyles = []string{"single_page", "one_question_at_a_time"}

var formBulkOperations = []struct{ verb, short string }{
	{"archive", "Archive forms"},
	{"unarchive", "Unarchive forms"},
	{"enable", "Enable forms"},
	{"disable", "Disable forms"},
	{"favorite", "Favorite forms"},
	{"unfavorite", "Unfavorite forms"},
}

var formsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a blank form",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateEnum("style", flagString(cmd, "style"), validFormStyles); err != nil {
			return err
		}

		return sendAndPrint(cmd, http.MethodPost, apiPath("forms"), map[string]any{
			"form": flagValues(cmd, "title", "style", "created-by-email"),
		}, printActionResult)
	},
}

var formsUpdateCmd = &cobra.Command{
	Use:   "update <form-id>",
	Short: "Rename a form",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return sendAndPrint(cmd, http.MethodPatch, formPath(args[0]), map[string]any{
			"form": flagValues(cmd, "title"),
		}, printResource)
	},
}

var formsBulkDeleteCmd = &cobra.Command{
	Use:   "bulk-delete <form-id>...",
	Short: "Delete several forms at once",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return sendAndPrint(cmd, http.MethodDelete, apiPath("forms", "bulk_destroy"), map[string]any{"ids": args}, printActionResult)
	},
}

func newFormBulkOperationCmd(verb, short string) *cobra.Command {
	return &cobra.Command{
		Use:   verb + " <form-id>...",
		Short: short,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := send(cmd, http.MethodPatch, apiPath("forms", "bulk_update"), map[string]any{
				"ids":       args,
				"operation": verb,
			})
			if err != nil {
				return err
			}

			printList(data, "forms", nil)
			return nil
		},
	}
}

func newFormActionCmd(verb, short, method string, segment string) *cobra.Command {
	return &cobra.Command{
		Use:   verb + " <form-id>",
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return sendAndPrint(cmd, method, formPath(args[0], segment), nil, printActionResult)
		},
	}
}

var formsCountsCmd = &cobra.Command{
	Use:   "counts",
	Short: "Count forms by status",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := fetch(cmd, apiPath("forms", "counts"), nil)
		if err != nil {
			return err
		}

		printResource(data, nil)
		return nil
	},
}

var formsTransferCmd = &cobra.Command{
	Use:   "transfer <form-id>...",
	Short: "Transfer forms to another workspace",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		body := flagValues(cmd, "to=destination_workspace")
		body["ids"] = args

		return sendAndPrint(cmd, http.MethodPost, apiPath("forms", "transfer"), body, printActionResult)
	},
}

func init() {
	formsCreateCmd.Flags().String("title", "", "Form title, up to 80 characters")
	formsCreateCmd.Flags().String("style", "", "Form style: single_page (default) or one_question_at_a_time")
	formsCreateCmd.Flags().String("created-by-email", "", "Email of the form owner; defaults to the signed-in user")
	markFlagsRequired(formsCreateCmd, "title")

	formsUpdateCmd.Flags().String("title", "", "New form title, up to 80 characters")
	markFlagsRequired(formsUpdateCmd, "title")

	formsTransferCmd.Flags().String("to", "", "Subdomain of the destination workspace, where you are an admin")
	markFlagsRequired(formsTransferCmd, "to")

	formsCmd.AddCommand(
		formsCreateCmd,
		formsUpdateCmd,
		deleteCmd("delete <form-id>", "Delete a form", 1, func(args []string) string { return formPath(args[0]) }, "Form deleted."),
		formsBulkDeleteCmd,
	)
	for _, operation := range formBulkOperations {
		formsCmd.AddCommand(newFormBulkOperationCmd(operation.verb, operation.short))
	}
	formsCmd.AddCommand(
		newFormActionCmd("clone", "Clone a form", http.MethodPost, "clone"),
		newFormActionCmd("publish", "Publish a form's draft", http.MethodPatch, "publish"),
		newFormActionCmd("reset", "Reset a form's draft to its published version", http.MethodPatch, "reset"),
		formsCountsCmd,
		formsTransferCmd,
	)
}
