package commands

import (
	"github.com/neetozone/neeto-cli-commons/output"
	"github.com/spf13/cobra"
)

var apiKeysCmd = &cobra.Command{
	Use:   "api-keys",
	Short: "Manage API keys",
}

var apiKeysListCmd = &cobra.Command{
	Use:   "list",
	Short: "List API keys",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := fetch(cmd, apiPath("api_keys"), queryValues(cmd, paginationParams(cmd), "search"))
		if err != nil {
			return err
		}
		printList(data, "api_keys", []output.Breadcrumb{
			{Label: "Revoke", Command: "neetoform api-keys delete <id>"},
		})
		return nil
	},
}

var apiKeysDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Revoke an API key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return removeAndConfirm(cmd, apiPath("api_keys", args[0]), nil, "API key revoked.")
	},
}

func init() {
	addPaginationFlags(apiKeysListCmd)
	apiKeysListCmd.Flags().String("search", "", "Filter by key label")
	apiKeysCmd.AddCommand(apiKeysListCmd, apiKeysDeleteCmd)
	register(func(root *cobra.Command) { root.AddCommand(apiKeysCmd) })
}
