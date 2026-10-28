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

func newThankYouPageSaveCmd(use, short string, nargs int, method string, render printer) *cobra.Command {
	save := &cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.ExactArgs(nargs),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateEnum("kind", flagString(cmd, "kind"), validThankYouKinds); err != nil {
				return err
			}
			body := map[string]any{"thank_you_page": flagValues(cmd, thankYouPageSpecs...)}
			return sendAndPrint(cmd, method, formResource("thank_you_pages")(args), body, render)
		},
	}
	save.Flags().String("title", "", "Page title")
	save.Flags().String("description", "", "Page description")
	save.Flags().String("kind", "", "custom_message or redirect_to_url")
	save.Flags().String("message", "", "HTML message shown for custom_message")
	save.Flags().String("redirect-url", "", "http(s) URL to redirect to; needed for redirect_to_url")
	save.Flags().Bool("show-redirect-message", false, "Show a message before redirecting")
	save.Flags().String("redirect-message", "", "HTML message shown before redirecting")
	save.Flags().Bool("pass-event-details", false, "Pass submission details to the redirect URL")
	save.Flags().Bool("social-sharing", false, "Show social sharing buttons")
	save.Flags().Bool("show-resubmit-link", false, "Show a link to submit another response")
	save.Flags().String("resubmit-link-text", "", "Text of the resubmit link; needed with --show-resubmit-link")
	save.Flags().Bool("primary", false, "Make this the primary thank-you page")
	return save
}

func init() {
	thankYouPagesCmd.AddCommand(
		listCmd("list <form-id>", "List thank-you pages", 1, formResource("thank_you_pages"), "thank_you_pages", []output.Breadcrumb{
			{Label: "Show", Command: "neetoform forms thank-you-pages show <form-id> <id>"},
		}),
		showCmd("show <form-id> <id>", "Show a thank-you page", 2, formResource("thank_you_pages")),
		newThankYouPageSaveCmd("create <form-id>", "Create a thank-you page", 1, http.MethodPost, printActionResult),
		newThankYouPageSaveCmd("update <form-id> <id>", "Update a thank-you page", 2, http.MethodPatch, printResource),
		deleteCmd("delete <form-id> <id>", "Delete a thank-you page", 2, formResource("thank_you_pages"), "Thank-you page deleted."),
	)
	formsCmd.AddCommand(thankYouPagesCmd)
}
