package commands

import (
	"net/http"

	"github.com/neetozone/neeto-cli-commons/output"
	"github.com/spf13/cobra"
)

var deactivatedTeamMembersCmd = &cobra.Command{
	Use:   "deactivated-team-members",
	Short: "Manage deactivated team members",
}

var deactivatedTeamMembersReactivateCmd = &cobra.Command{
	Use:   "reactivate <id>",
	Short: "Reactivate a team member",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return sendAndPrint(cmd, http.MethodPatch, apiPath("deactivated_team_members", args[0], "reactivate"), map[string]any{}, printResource)
	},
}

func init() {
	deactivatedTeamMembersCmd.AddCommand(
		listCmd("list", "List deactivated team members", 0, workspaceResource("deactivated_team_members"), "team_members", []output.Breadcrumb{
			{Label: "Reactivate", Command: "neetoform deactivated-team-members reactivate <id>"},
		}),
		deactivatedTeamMembersReactivateCmd,
	)
	register(func(root *cobra.Command) { root.AddCommand(deactivatedTeamMembersCmd) })
}
