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

var customDomainsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List custom domains",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := fetch(cmd, apiPath("custom_domains"), paginationParams(cmd))
		if err != nil {
			return err
		}
		printList(data, "custom_domains", []output.Breadcrumb{
			{Label: "Show", Command: "neetoform custom-domains show <hostname>"},
		})
		return nil
	},
}

var customDomainsShowCmd = &cobra.Command{
	Use:   "show <hostname>",
	Short: "Show a custom domain and its DNS records",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := fetch(cmd, apiPath("custom_domains", args[0]), nil)
		if err != nil {
			return err
		}
		printResource(data, nil)
		return nil
	},
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

var customDomainsDeleteCmd = &cobra.Command{
	Use:   "delete <hostname>",
	Short: "Remove a custom domain",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return removeAndConfirm(cmd, apiPath("custom_domains", args[0]), nil, "Custom domain removed.")
	},
}

func init() {
	addPaginationFlags(customDomainsListCmd)
	customDomainsCmd.AddCommand(customDomainsListCmd, customDomainsShowCmd, customDomainsCreateCmd, customDomainsDeleteCmd)
	register(func(root *cobra.Command) { root.AddCommand(customDomainsCmd) })
}
