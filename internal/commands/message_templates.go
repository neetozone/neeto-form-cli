package commands

import (
	"net/http"

	"github.com/neetozone/neeto-cli-commons/output"
	"github.com/spf13/cobra"
)

var (
	validMessageTemplateTypes    = []string{"email", "sms"}
	validMessageTemplateStatuses = []string{"active", "inactive"}
)

var messageTemplatesCmd = &cobra.Command{
	Use:   "message-templates",
	Short: "Manage a form's email and SMS templates",
}

var messageTemplatesListCmd = &cobra.Command{
	Use:   "list <form-id>",
	Short: "List message templates",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateMessageTemplateFlags(cmd); err != nil {
			return err
		}
		params := queryValues(cmd, paginationParams(cmd), "type=template_type", "status")
		data, err := fetch(cmd, formPath(args[0], "message_templates"), params)
		if err != nil {
			return err
		}
		printList(data, "message_templates", []output.Breadcrumb{
			{Label: "Show", Command: "neetoform forms message-templates show <form-id> <id>"},
		})
		return nil
	},
}

var messageTemplatesShowCmd = &cobra.Command{
	Use:   "show <form-id> <id>",
	Short: "Show a message template",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := fetch(cmd, formPath(args[0], "message_templates", args[1]), nil)
		if err != nil {
			return err
		}
		printResource(data, nil)
		return nil
	},
}

var messageTemplatesCreateCmd = &cobra.Command{
	Use:   "create <form-id>",
	Short: "Create a message template",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateMessageTemplateFlags(cmd); err != nil {
			return err
		}
		body := map[string]any{"message_template": flagValues(cmd, "name", "type=template_type", "body", "subject", "status")}
		data, err := send(cmd, http.MethodPost, formPath(args[0], "message_templates"), body)
		if err != nil {
			return err
		}
		printActionResult(data, nil)
		return nil
	},
}

var messageTemplatesUpdateCmd = &cobra.Command{
	Use:   "update <form-id> <id>",
	Short: "Update a message template",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateMessageTemplateFlags(cmd); err != nil {
			return err
		}
		body := map[string]any{"message_template": flagValues(cmd, "name", "status", "subject", "body")}
		data, err := send(cmd, http.MethodPatch, formPath(args[0], "message_templates", args[1]), body)
		if err != nil {
			return err
		}
		printResource(data, nil)
		return nil
	},
}

var messageTemplatesDeleteCmd = &cobra.Command{
	Use:   "delete <form-id> <id>",
	Short: "Delete a message template",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return removeAndConfirm(cmd, formPath(args[0], "message_templates", args[1]), nil, "Message template deleted.")
	},
}

func validateMessageTemplateFlags(cmd *cobra.Command) error {
	if cmd.Flags().Lookup("type") != nil {
		if err := validateEnum("type", flagString(cmd, "type"), validMessageTemplateTypes); err != nil {
			return err
		}
	}
	return validateEnum("status", flagString(cmd, "status"), validMessageTemplateStatuses)
}

func init() {
	addPaginationFlags(messageTemplatesListCmd)
	messageTemplatesListCmd.Flags().String("type", "", "Only templates of this type: email or sms")
	messageTemplatesListCmd.Flags().String("status", "", "Only templates with this status: active or inactive")

	messageTemplatesCreateCmd.Flags().String("name", "", "Template name")
	messageTemplatesCreateCmd.Flags().String("type", "", "email or sms")
	messageTemplatesCreateCmd.Flags().String("body", "", "Template body; HTML for email")
	messageTemplatesCreateCmd.Flags().String("subject", "", "Email subject; needed for email templates")
	messageTemplatesCreateCmd.Flags().String("status", "", "active or inactive")
	markFlagsRequired(messageTemplatesCreateCmd, "name", "type", "body")

	messageTemplatesUpdateCmd.Flags().String("name", "", "Template name")
	messageTemplatesUpdateCmd.Flags().String("status", "", "active or inactive")
	messageTemplatesUpdateCmd.Flags().String("subject", "", "Email subject")
	messageTemplatesUpdateCmd.Flags().String("body", "", "Template body; HTML for email")

	messageTemplatesCmd.AddCommand(messageTemplatesListCmd, messageTemplatesShowCmd, messageTemplatesCreateCmd, messageTemplatesUpdateCmd, messageTemplatesDeleteCmd)
	formsCmd.AddCommand(messageTemplatesCmd)
}
