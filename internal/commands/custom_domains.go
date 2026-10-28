package commands

import (
	"net/http"

	"github.com/neetozone/neeto-cli-commons/output"
	"github.com/spf13/cobra"
)

var customDomainsCmd = &cobra.Command{
	Use:   "custom-domains",
	Short: "Manage custom domains",
}

var customDomainsCreateCmd = &cobra.Command{
	Use:   "create <hostname>",
	Short: "Add a custom domain",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := send(cmd, http.MethodPost, apiPath("custom_domains"), map[string]any{"hostname": args[0]})
		if err != nil {
			return err
		}
		printActionResult(data, []output.Breadcrumb{
			{Label: "Check status", Command: "neetoform custom-domains show " + args[0]},
		})
		return nil
	},
}

func init() {
	customDomainsCmd.AddCommand(
		listCmd("list", "List custom domains", 0, workspaceResource("custom_domains"), "custom_domains", []output.Breadcrumb{
			{Label: "Show", Command: "neetoform custom-domains show <hostname>"},
		}),
		showCmd("show <hostname>", "Show a custom domain and its DNS records", 1, workspaceResource("custom_domains")),
		customDomainsCreateCmd,
		deleteCmd("delete <hostname>", "Remove a custom domain", 1, workspaceResource("custom_domains"), "Custom domain removed."),
	)
	register(func(root *cobra.Command) { root.AddCommand(customDomainsCmd) })
}
