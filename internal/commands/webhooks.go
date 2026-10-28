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

func createWebhook(cmd *cobra.Command, path string) error {
	if err := validateEnumList("events", splitCSV(flagString(cmd, "events")), validWebhookEvents); err != nil {
		return err
	}
	body := flagValues(cmd, "endpoint", "events[]", "secret")
	return sendAndPrint(cmd, http.MethodPost, path, map[string]any{"webhook": body}, printActionResult)
}

func init() {
	scopes := []struct {
		group       *cobra.Command
		argsUse     string
		nargs       int
		path        argsPath
		listShort   string
		createShort string
	}{
		{webhooksCmd, "", 0, workspaceResource("webhooks"), "List workspace webhooks", "Create a workspace webhook"},
		{formWebhooksCmd, " <form-id>", 1, formResource("webhooks"), "List a form's webhooks", "Create a webhook for a form"},
	}
	for _, scope := range scopes {
		list := listCmd("list"+scope.argsUse, scope.listShort, scope.nargs, scope.path, "webhooks", []output.Breadcrumb{
			{Label: "Show", Command: "neetoform webhooks show <id>"},
		})
		create := &cobra.Command{
			Use:   "create" + scope.argsUse,
			Short: scope.createShort,
			Args:  exactArgs(scope.nargs),
			RunE: func(cmd *cobra.Command, args []string) error {
				return createWebhook(cmd, scope.path(args[:scope.nargs]))
			},
		}
		create.Flags().String("endpoint", "", "URL that receives the webhook calls")
		create.Flags().String("events", "", "Comma-separated events to fire on: submit. All events when omitted")
		create.Flags().String("secret", "", "Secret used to sign calls; never shown again")
		markFlagsRequired(create, "endpoint")
		scope.group.AddCommand(list, create)
	}

	webhooksCmd.AddCommand(
		showCmd("show <id>", "Show a webhook", 1, workspaceResource("webhooks")),
		deleteCmd("delete <id>", "Delete a webhook", 1, workspaceResource("webhooks"), "Webhook deleted."),
	)
	formsCmd.AddCommand(formWebhooksCmd)
	register(func(root *cobra.Command) { root.AddCommand(webhooksCmd) })
}
