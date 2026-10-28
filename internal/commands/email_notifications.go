package commands

import (
	"net/http"

	"github.com/spf13/cobra"
)

var validEmailNotificationKinds = []string{"admin", "respondent"}

var emailNotificationSpecs = []string{
	"enabled=is_enabled", "subject", "message", "to[]=notify_emails", "cc[]=notify_cc_emails",
	"bcc[]=notify_bcc_emails", "reply-to-question-id", "reply-to", "send-to-question-id",
}

var emailNotificationsCmd = &cobra.Command{
	Use:   "email-notifications",
	Short: "Manage a form's submission emails",
}

var emailNotificationsShowCmd = &cobra.Command{
	Use:   "show <form-id> <admin|respondent>",
	Short: "Show the admin or respondent email settings",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateEmailNotificationKind(args[1]); err != nil {
			return err
		}
		data, err := fetch(cmd, formPath(args[0], "email_notifications", args[1]), nil)
		if err != nil {
			return err
		}
		printResource(data, nil)
		return nil
	},
}

var emailNotificationsUpdateCmd = &cobra.Command{
	Use:   "update <form-id> <admin|respondent>",
	Short: "Update the admin or respondent email settings",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateEmailNotificationKind(args[1]); err != nil {
			return err
		}
		body := map[string]any{"email_notification": flagValues(cmd, emailNotificationSpecs...)}
		return sendAndPrint(cmd, http.MethodPatch, formPath(args[0], "email_notifications", args[1]), body, printResource)
	},
}

func validateEmailNotificationKind(kind string) error {
	return validateEnum("kind", kind, validEmailNotificationKinds)
}

func init() {
	flags := emailNotificationsUpdateCmd.Flags()
	flags.Bool("enabled", false, "Send this email after a submission")
	flags.String("subject", "", "Email subject")
	flags.String("message", "", "HTML email body")
	flags.String("to", "", "Comma-separated admin recipient emails")
	flags.String("cc", "", "Comma-separated admin CC emails")
	flags.String("bcc", "", "Comma-separated admin BCC emails")
	flags.String("reply-to-question-id", "", "UUID of the email question used as reply-to; admin only")
	flags.String("reply-to", "", "Reply-to email address; respondent only")
	flags.String("send-to-question-id", "", "UUID of the email question whose answer receives the email; respondent only")

	emailNotificationsCmd.AddCommand(emailNotificationsShowCmd, emailNotificationsUpdateCmd)
	formsCmd.AddCommand(emailNotificationsCmd)
}
