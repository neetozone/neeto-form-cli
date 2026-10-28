package commands

import (
	"github.com/neetozone/neeto-cli-commons/output"
	"github.com/spf13/cobra"
)

var apiKeysCmd = &cobra.Command{
	Use:   "api-keys",
	Short: "Manage API keys",
}

func init() {
	list := listCmd("list", "List API keys", 0, workspaceResource("api_keys"), "api_keys", []output.Breadcrumb{
		{Label: "Revoke", Command: "neetoform api-keys delete <id>"},
	}, "search")
	list.Flags().String("search", "", "Filter by key label")
	apiKeysCmd.AddCommand(list, deleteCmd("delete <id>", "Revoke an API key", 1, workspaceResource("api_keys"), "API key revoked."))
	register(func(root *cobra.Command) { root.AddCommand(apiKeysCmd) })
}
