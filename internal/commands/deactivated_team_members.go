package commands

import (
	"github.com/neetozone/neeto-cli-commons/output"
	"github.com/spf13/cobra"
)

var deactivatedTeamMembersCmd = &cobra.Command{
	Use:   "deactivated-team-members",
	Short: "Manage deactivated team members",
}

var deactivatedTeamMembersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List deactivated team members",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := fetch(cmd, apiPath("deactivated_team_members"), paginationParams(cmd))
		if err != nil {
			return err
		}
		printList(data, "team_members", []output.Breadcrumb{
			{Label: "Reactivate", Command: "neetoform deactivated-team-members reactivate <id>"},
		})
		return nil
	},
}

var deactivatedTeamMembersReactivateCmd = &cobra.Command{
	Use:   "reactivate <id>",
	Short: "Reactivate a team member",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := send(cmd, "PATCH", apiPath("deactivated_team_members", args[0], "reactivate"), map[string]any{})
		if err != nil {
			return err
		}
		printResource(data, nil)
		return nil
	},
}

func init() {
	addPaginationFlags(deactivatedTeamMembersListCmd)
	deactivatedTeamMembersCmd.AddCommand(deactivatedTeamMembersListCmd, deactivatedTeamMembersReactivateCmd)
	register(func(root *cobra.Command) { root.AddCommand(deactivatedTeamMembersCmd) })
}
