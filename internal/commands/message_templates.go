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

func messageTemplateBody(cmd *cobra.Command) (map[string]any, error) {
	if err := validateMessageTemplateFlags(cmd); err != nil {
		return nil, err
	}
	return map[string]any{"message_template": flagValues(cmd, "name", "type=template_type", "body", "subject", "status")}, nil
}

var messageTemplatesCreateCmd = saveCmd("create <form-id>", "Create a message template", 1, http.MethodPost,
	formResource("message_templates"), printActionResult, messageTemplateBody)

var messageTemplatesUpdateCmd = saveCmd("update <form-id> <id>", "Update a message template", 2, http.MethodPatch,
	formResource("message_templates"), printResource, messageTemplateBody)

func validateMessageTemplateFlags(cmd *cobra.Command) error {
	if err := validateEnum("type", flagString(cmd, "type"), validMessageTemplateTypes); err != nil {
		return err
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

	messageTemplatesCmd.AddCommand(
		messageTemplatesListCmd,
		showCmd("show <form-id> <id>", "Show a message template", 2, formResource("message_templates")),
		messageTemplatesCreateCmd,
		messageTemplatesUpdateCmd,
		deleteCmd("delete <form-id> <id>", "Delete a message template", 2, formResource("message_templates"), "Message template deleted."),
	)
	formsCmd.AddCommand(messageTemplatesCmd)
}
