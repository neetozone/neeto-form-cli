package commands

import (
	"net/http"

	"github.com/neetozone/neeto-cli-commons/output"
	"github.com/spf13/cobra"
)

var validThankYouKinds = []string{"custom_message", "redirect_to_url"}

var thankYouPageSpecs = []string{
	"title", "description", "kind", "message", "redirect-url", "show-redirect-message",
	"redirect-message", "pass-event-details=pass_event_details_to_redirect_url",
	"social-sharing=social_sharing_enabled", "show-resubmit-link", "resubmit-link-text", "primary",
}

var thankYouPagesCmd = &cobra.Command{
	Use:   "thank-you-pages",
	Short: "Manage a form's thank-you pages",
}

var thankYouPagesListCmd = &cobra.Command{
	Use:   "list <form-id>",
	Short: "List thank-you pages",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := fetch(cmd, formPath(args[0], "thank_you_pages"), paginationParams(cmd))
		if err != nil {
			return err
		}
		printList(data, "thank_you_pages", []output.Breadcrumb{
			{Label: "Show", Command: "neetoform forms thank-you-pages show <form-id> <id>"},
		})
		return nil
	},
}

var thankYouPagesShowCmd = &cobra.Command{
	Use:   "show <form-id> <id>",
	Short: "Show a thank-you page",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := fetch(cmd, formPath(args[0], "thank_you_pages", args[1]), nil)
		if err != nil {
			return err
		}
		printResource(data, nil)
		return nil
	},
}

var thankYouPagesCreateCmd = &cobra.Command{
	Use:   "create <form-id>",
	Short: "Create a thank-you page",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateEnum("kind", flagString(cmd, "kind"), validThankYouKinds); err != nil {
			return err
		}
		data, err := send(cmd, http.MethodPost, formPath(args[0], "thank_you_pages"), map[string]any{"thank_you_page": flagValues(cmd, thankYouPageSpecs...)})
		if err != nil {
			return err
		}
		printActionResult(data, nil)
		return nil
	},
}

var thankYouPagesUpdateCmd = &cobra.Command{
	Use:   "update <form-id> <id>",
	Short: "Update a thank-you page",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateEnum("kind", flagString(cmd, "kind"), validThankYouKinds); err != nil {
			return err
		}
		data, err := send(cmd, http.MethodPatch, formPath(args[0], "thank_you_pages", args[1]), map[string]any{"thank_you_page": flagValues(cmd, thankYouPageSpecs...)})
		if err != nil {
			return err
		}
		printResource(data, nil)
		return nil
	},
}

var thankYouPagesDeleteCmd = &cobra.Command{
	Use:   "delete <form-id> <id>",
	Short: "Delete a thank-you page",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return removeAndConfirm(cmd, formPath(args[0], "thank_you_pages", args[1]), nil, "Thank-you page deleted.")
	},
}

func addThankYouPageFlags(cmd *cobra.Command) {
	cmd.Flags().String("title", "", "Page title")
	cmd.Flags().String("description", "", "Page description")
	cmd.Flags().String("kind", "", "custom_message or redirect_to_url")
	cmd.Flags().String("message", "", "HTML message shown for custom_message")
	cmd.Flags().String("redirect-url", "", "http(s) URL to redirect to; needed for redirect_to_url")
	cmd.Flags().Bool("show-redirect-message", false, "Show a message before redirecting")
	cmd.Flags().String("redirect-message", "", "HTML message shown before redirecting")
	cmd.Flags().Bool("pass-event-details", false, "Pass submission details to the redirect URL")
	cmd.Flags().Bool("social-sharing", false, "Show social sharing buttons")
	cmd.Flags().Bool("show-resubmit-link", false, "Show a link to submit another response")
	cmd.Flags().String("resubmit-link-text", "", "Text of the resubmit link; needed with --show-resubmit-link")
	cmd.Flags().Bool("primary", false, "Make this the primary thank-you page")
}

func init() {
	addPaginationFlags(thankYouPagesListCmd)
	addThankYouPageFlags(thankYouPagesCreateCmd)
	addThankYouPageFlags(thankYouPagesUpdateCmd)

	thankYouPagesCmd.AddCommand(thankYouPagesListCmd, thankYouPagesShowCmd, thankYouPagesCreateCmd, thankYouPagesUpdateCmd, thankYouPagesDeleteCmd)
	formsCmd.AddCommand(thankYouPagesCmd)
}
