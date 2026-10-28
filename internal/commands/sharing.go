package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

var shareEmbedTypes = []string{"inline", "floating-popup", "element-popup", "iframe"}

var shareCmd = &cobra.Command{
	Use:   "share",
	Short: "View how a form is shared",
}

var shareShowCmd = &cobra.Command{
	Use:   "show <form-id>",
	Short: "Show a form's public link and embed codes",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		embed := flagString(cmd, "embed")
		if err := validateEnum("embed", embed, shareEmbedTypes); err != nil {
			return err
		}

		data, err := fetch(cmd, formPath(args[0], "share"), nil)
		if err != nil {
			return err
		}

		if embed == "" {
			printResource(data, nil)
			return nil
		}

		var share struct {
			EmbedCodes map[string]string `json:"embed_codes"`
		}
		if err := json.Unmarshal(data, &share); err != nil {
			return err
		}
		snippet := share.EmbedCodes[strings.ReplaceAll(embed, "-", "_")]
		if app.Printer.Quiet {
			_, err := fmt.Fprintln(app.Printer.Out, snippet)
			return err
		}
		printMessage(snippet)
		return nil
	},
}

var publicLinkCmd = &cobra.Command{
	Use:   "public-link",
	Short: "Manage a form's public link",
}

var publicLinkUpdateCmd = &cobra.Command{
	Use:   "update <form-id>",
	Short: "Change a form's public link",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return sendAndPrint(cmd, http.MethodPatch, formPath(args[0], "public_link"), flagValues(cmd, "public-link-id"), printResource)
	},
}

var prefilledLinksCmd = &cobra.Command{
	Use:   "prefilled-links",
	Short: "Build links that open a form with answers filled in",
}

var prefilledLinksCreateCmd = &cobra.Command{
	Use:   "create <form-id>",
	Short: "Build a prefilled link for a form",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pairs, _ := cmd.Flags().GetStringArray("answer")
		parsed, err := parseKeyValues(pairs, "answer", "<question-id>=<value>")
		if err != nil {
			return err
		}
		answers := []map[string]any{}
		for _, pair := range parsed {
			answers = append(answers, map[string]any{"question_id": pair.key, "value": pair.value})
		}

		if raw := flagString(cmd, "answers-json"); raw != "" {
			var extra []map[string]any
			if err := json.Unmarshal([]byte(raw), &extra); err != nil {
				return fmt.Errorf("Invalid --answers-json: %w", err)
			}
			answers = append(answers, extra...)
		}

		body := flagValues(cmd, "discount-code")
		body["answers"] = answers

		return sendAndPrint(cmd, http.MethodPost, formPath(args[0], "prefilled_link"), body, printActionResult)
	},
}

var emailInvitationsCmd = &cobra.Command{
	Use:   "email-invitations",
	Short: "Email people an invitation to fill in a form",
}

var emailInvitationsCreateCmd = &cobra.Command{
	Use:   "create <form-id>",
	Short: "Send email invitations for a form",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		body := flagValues(cmd,
			"emails[]",
			"subject=email_subject",
			"body=email_body",
			"button-text=start_survey_button_text",
			"reply-to=reply_to_email",
			"from=send_from_email",
			"send-at",
		)

		return sendAndPrint(cmd, http.MethodPost, formPath(args[0], "email_invitations"), body, printActionResult)
	},
}

func init() {
	shareShowCmd.Flags().String("embed", "", "Print only one embed snippet: inline, floating-popup, element-popup or iframe")
	shareCmd.AddCommand(shareShowCmd)

	publicLinkUpdateCmd.Flags().String("public-link-id", "", "New public link id (letters, numbers, - and _). Omit for a random one. The old link stops working at once")
	publicLinkCmd.AddCommand(publicLinkUpdateCmd)

	prefilledLinksCreateCmd.Flags().StringArray("answer", nil, "Answer as <question-id>=<value>; repeat for more answers")
	prefilledLinksCreateCmd.Flags().String("answers-json", "", `JSON array of answers for compound values, e.g. '[{"question_id":"<id>","value":{"country":"US","number":"2025550123"}}]'`)
	prefilledLinksCreateCmd.Flags().String("discount-code", "", "Discount code to add to the link")
	prefilledLinksCmd.AddCommand(prefilledLinksCreateCmd)

	emailInvitationsCreateCmd.Flags().String("emails", "", "Comma-separated email addresses to invite")
	emailInvitationsCreateCmd.Flags().String("subject", "", "Email subject")
	emailInvitationsCreateCmd.Flags().String("body", "", "Email body in HTML")
	emailInvitationsCreateCmd.Flags().String("button-text", "", "Text of the button that opens the form")
	emailInvitationsCreateCmd.Flags().String("reply-to", "", "Reply-to email address")
	emailInvitationsCreateCmd.Flags().String("from", "", "Sender email address")
	emailInvitationsCreateCmd.Flags().String("send-at", "", "Send at this ISO 8601 date-time, e.g. 2026-11-01T10:00:00Z. Omit to send now")
	markFlagsRequired(emailInvitationsCreateCmd, "emails")
	emailInvitationsCmd.AddCommand(emailInvitationsCreateCmd)

	formsCmd.AddCommand(shareCmd, publicLinkCmd, prefilledLinksCmd, emailInvitationsCmd)
}
