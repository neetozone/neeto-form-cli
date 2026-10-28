package commands

import (
	"net/url"

	"github.com/neetozone/neeto-cli-commons/output"
	"github.com/spf13/cobra"
)

var templatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "Browse form templates",
}

var templatesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		params := queryValues(cmd, paginationParams(cmd), "search", "category")
		data, err := fetch(cmd, apiPath("templates"), params)
		if err != nil {
			return err
		}
		printList(data, "templates", []output.Breadcrumb{
			{Label: "Show", Command: "neetoform templates show <id-or-slug>"},
			{Label: "Use", Command: "neetoform templates use <template-id>"},
		})
		return nil
	},
}

var templatesShowCmd = &cobra.Command{
	Use:   "show <id-or-slug>",
	Short: "Show a template",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := fetch(cmd, apiPath("templates", args[0]), url.Values{})
		if err != nil {
			return err
		}
		printResource(data, nil)
		return nil
	},
}

var templatesUseCmd = &cobra.Command{
	Use:   "use <template-id>",
	Short: "Create a form from a template",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		body := flagValues(cmd, "created-by-email")
		data, err := send(cmd, "POST", apiPath("templates", args[0], "forms"), body)
		if err != nil {
			return err
		}
		printActionResult(data, nil)
		return nil
	},
}

func init() {
	addPaginationFlags(templatesListCmd)
	templatesListCmd.Flags().String("search", "", "Filter by title text")
	templatesListCmd.Flags().String("category", "", "Filter by category name")
	templatesUseCmd.Flags().String("created-by-email", "", "Email of the active team member who owns the new form; defaults to the signed-in user")

	templatesCmd.AddCommand(templatesListCmd, templatesShowCmd, templatesUseCmd)
	register(func(root *cobra.Command) { root.AddCommand(templatesCmd) })
}
