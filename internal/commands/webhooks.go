package commands

import (
	"net/http"

	"github.com/neetozone/neeto-cli-commons/output"
	"github.com/spf13/cobra"
)

var validWebhookEvents = []string{"submit"}

var webhooksCmd = &cobra.Command{
	Use:   "webhooks",
	Short: "Manage workspace webhooks",
}

var formWebhooksCmd = &cobra.Command{
	Use:   "webhooks",
	Short: "Manage a form's webhooks",
}

func listWebhooks(cmd *cobra.Command, path string) error {
	data, err := fetch(cmd, path, paginationParams(cmd))
	if err != nil {
		return err
	}
	printList(data, "webhooks", []output.Breadcrumb{
		{Label: "Show", Command: "neetoform webhooks show <id>"},
	})
	return nil
}

func createWebhook(cmd *cobra.Command, path string) error {
	body := flagValues(cmd, "endpoint", "events[]", "secret")
	if events, ok := body["events"].([]string); ok {
		for _, event := range events {
			if err := validateEnum("events", event, validWebhookEvents); err != nil {
				return err
			}
		}
	}
	data, err := send(cmd, http.MethodPost, path, map[string]any{"webhook": body})
	if err != nil {
		return err
	}
	printActionResult(data, nil)
	return nil
}

func addWebhookCreateFlags(cmd *cobra.Command) {
	cmd.Flags().String("endpoint", "", "URL that receives the webhook calls")
	cmd.Flags().String("events", "", "Comma-separated events to fire on: submit. All events when omitted")
	cmd.Flags().String("secret", "", "Secret used to sign calls; never shown again")
	markFlagsRequired(cmd, "endpoint")
}

var webhooksListCmd = &cobra.Command{
	Use:   "list",
	Short: "List workspace webhooks",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listWebhooks(cmd, apiPath("webhooks"))
	},
}

var webhooksShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show a webhook",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := fetch(cmd, apiPath("webhooks", args[0]), nil)
		if err != nil {
			return err
		}
		printResource(data, nil)
		return nil
	},
}

var webhooksCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a workspace webhook",
	RunE: func(cmd *cobra.Command, args []string) error {
		return createWebhook(cmd, apiPath("webhooks"))
	},
}

var webhooksDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a webhook",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return removeAndConfirm(cmd, apiPath("webhooks", args[0]), nil, "Webhook deleted.")
	},
}

var formWebhooksListCmd = &cobra.Command{
	Use:   "list <form-id>",
	Short: "List a form's webhooks",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return listWebhooks(cmd, formPath(args[0], "webhooks"))
	},
}

var formWebhooksCreateCmd = &cobra.Command{
	Use:   "create <form-id>",
	Short: "Create a webhook for a form",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return createWebhook(cmd, formPath(args[0], "webhooks"))
	},
}

func init() {
	addPaginationFlags(webhooksListCmd)
	addPaginationFlags(formWebhooksListCmd)
	addWebhookCreateFlags(webhooksCreateCmd)
	addWebhookCreateFlags(formWebhooksCreateCmd)

	webhooksCmd.AddCommand(webhooksListCmd, webhooksShowCmd, webhooksCreateCmd, webhooksDeleteCmd)
	formWebhooksCmd.AddCommand(formWebhooksListCmd, formWebhooksCreateCmd)
	formsCmd.AddCommand(formWebhooksCmd)
	register(func(root *cobra.Command) { root.AddCommand(webhooksCmd) })
}
