package commands

import (
	"net/http"

	"github.com/neetozone/neeto-cli-commons/output"
	"github.com/spf13/cobra"
)

var validSubmissionStatuses = []string{"completed", "incomplete", "all"}

var validExportTypes = []string{"csv", "xlsx", "pdf"}

var submissionsCmd = &cobra.Command{
	Use:   "submissions",
	Short: "Manage form submissions",
}

var submissionsListCmd = &cobra.Command{
	Use:   "list <form-id>",
	Short: "List submissions for a form",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateEnum("status", flagString(cmd, "status"), validSubmissionStatuses); err != nil {
			return err
		}

		params := queryValues(cmd, paginationParams(cmd), "status", "search", "submitted-after", "submitted-before")
		if err := setOrderParam(cmd, params); err != nil {
			return err
		}

		data, err := fetch(cmd, formPath(args[0], "submissions"), params)
		if err != nil {
			return err
		}

		printList(data, "submissions", []output.Breadcrumb{
			{Label: "Show", Command: "neetoform forms submissions show <form-id> <submission-id>"},
		})
		return nil
	},
}

var submissionsDeleteCmd = deleteCmd("delete <form-id> <submission-id>", "Move a submission to the trash", 2,
	func(args []string) string { return formPath(args[0], "submissions", args[1]) }, "Submission moved to the trash.")

var submissionsBulkDeleteCmd = &cobra.Command{
	Use:   "bulk-delete <form-id>",
	Short: "Move several submissions to the trash",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return sendAndPrint(cmd, http.MethodDelete, formPath(args[0], "submissions", "bulk_destroy"), flagValues(cmd, "ids[]"), printActionResult)
	},
}

var submissionsSendResponsesCmd = &cobra.Command{
	Use:   "send-responses <form-id>",
	Short: "Email a message to the respondents of submissions",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		body := flagValues(cmd, "ids[]", "email-question-id", "subject", "body", "cc-emails[]=cc_emails", "bcc-emails[]=bcc_emails")
		return sendAndPrint(cmd, http.MethodPost, formPath(args[0], "submissions", "send_responses"), body, printActionResult)
	},
}

var submissionsExportCmd = &cobra.Command{
	Use:   "export <form-id>",
	Short: "Email an export of submissions",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateEnum("type", flagString(cmd, "type"), validExportTypes); err != nil {
			return err
		}
		body := flagValues(cmd, "type=export_type", "ids[]=submission_ids", "recipient-email",
			"include-hidden-fields", "include-meta-data", "include-question-descriptions")
		return sendAndPrint(cmd, http.MethodPost, formPath(args[0], "exports"), body, printActionResult)
	},
}

func init() {
	addPaginationFlags(submissionsListCmd)
	submissionsListCmd.Flags().String("submitted-after", "", "Only submissions at or after this ISO 8601 date (YYYY-MM-DD) or date-time")
	submissionsListCmd.Flags().String("submitted-before", "", "Only submissions before this ISO 8601 date (YYYY-MM-DD) or date-time")
	submissionsListCmd.Flags().String("order", "", "Sort by submission time: asc (oldest first, default) or desc")

	submissionsListCmd.Flags().String("status", "", "Which submissions to list: completed (default), incomplete or all")
	submissionsListCmd.Flags().String("search", "", "Only submissions whose answers match this text")

	submissionsBulkDeleteCmd.Flags().String("ids", "", "Comma-separated submission IDs to move to the trash")
	markFlagsRequired(submissionsBulkDeleteCmd, "ids")

	submissionsSendResponsesCmd.Flags().String("ids", "", "Comma-separated submission IDs to email, up to 100")
	submissionsSendResponsesCmd.Flags().String("email-question-id", "", "ID of the email question whose answer is the recipient address")
	submissionsSendResponsesCmd.Flags().String("subject", "", "Email subject")
	submissionsSendResponsesCmd.Flags().String("body", "", "Email body as HTML")
	submissionsSendResponsesCmd.Flags().String("cc-emails", "", "Comma-separated CC addresses")
	submissionsSendResponsesCmd.Flags().String("bcc-emails", "", "Comma-separated BCC addresses")
	markFlagsRequired(submissionsSendResponsesCmd, "ids", "email-question-id", "subject", "body")

	submissionsExportCmd.Flags().String("type", "", "Export format: csv, xlsx or pdf")
	submissionsExportCmd.Flags().String("ids", "", "Comma-separated submission IDs to export (default: all completed submissions)")
	submissionsExportCmd.Flags().String("recipient-email", "", "Team member email that receives the download link (default: you)")
	submissionsExportCmd.Flags().Bool("include-hidden-fields", false, "Include hidden fields")
	submissionsExportCmd.Flags().Bool("include-meta-data", false, "Include submission meta data")
	submissionsExportCmd.Flags().Bool("include-question-descriptions", false, "Include question descriptions")
	markFlagsRequired(submissionsExportCmd, "type")

	submissionsCmd.AddCommand(submissionsListCmd,
		showCmd("show <form-id> <submission-id>", "Show a submission", 2, formResource("submissions")),
		submissionsDeleteCmd, submissionsBulkDeleteCmd,
		submissionsSendResponsesCmd, submissionsExportCmd)
	formsCmd.AddCommand(submissionsCmd)
}
