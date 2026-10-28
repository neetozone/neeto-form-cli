package commands

import (
	"net/http"

	"github.com/neetozone/neeto-cli-commons/output"
	"github.com/spf13/cobra"
)

var templatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "Browse form templates",
}

var templatesUseCmd = &cobra.Command{
	Use:   "use <template-id>",
	Short: "Create a form from a template",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		body := flagValues(cmd, "created-by-email")
		return sendAndPrint(cmd, http.MethodPost, apiPath("templates", args[0], "forms"), body, printActionResult)
	},
}

func init() {
	templatesListCmd := listCmd("list", "List templates", 0, workspaceResource("templates"), "templates", []output.Breadcrumb{
		{Label: "Show", Command: "neetoform templates show <id-or-slug>"},
		{Label: "Use", Command: "neetoform templates use <template-id>"},
	}, "search", "category")
	templatesListCmd.Flags().String("search", "", "Filter by title text")
	templatesListCmd.Flags().String("category", "", "Filter by category name")
	templatesUseCmd.Flags().String("created-by-email", "", "Email of the active team member who owns the new form; defaults to the signed-in user")

	templatesCmd.AddCommand(
		templatesListCmd,
		showCmd("show <id-or-slug>", "Show a template", 1, workspaceResource("templates")),
		templatesUseCmd,
	)
	register(func(root *cobra.Command) { root.AddCommand(templatesCmd) })
}
