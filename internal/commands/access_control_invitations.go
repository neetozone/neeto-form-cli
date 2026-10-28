package commands

import (
	"net/http"

	"github.com/spf13/cobra"
)

var accessControlInvitationsCmd = &cobra.Command{
	Use:   "access-control-invitations",
	Short: "Manage who is invited to a form",
}

var accessControlInvitationsCreateCmd = &cobra.Command{
	Use:   "create <form-id>",
	Short: "Invite emails to a form",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return sendAndPrint(cmd, http.MethodPost, formPath(args[0], "access_control_invitations"), flagValues(cmd, "emails[]"), printActionResult)
	},
}

var accessControlInvitationsDeleteCmd = &cobra.Command{
	Use:   "delete <form-id>",
	Short: "Remove invited emails from a form",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return sendAndPrint(cmd, http.MethodDelete, formPath(args[0], "access_control_invitations"), flagValues(cmd, "emails[]"), printActionResult)
	},
}

func init() {
	accessControlInvitationsCreateCmd.Flags().String("emails", "", "Comma-separated emails to invite, 1 to 100. The form access must be set to invited")
	accessControlInvitationsDeleteCmd.Flags().String("emails", "", "Comma-separated emails to remove, 1 to 100")
	markFlagsRequired(accessControlInvitationsCreateCmd, "emails")
	markFlagsRequired(accessControlInvitationsDeleteCmd, "emails")

	accessControlInvitationsCmd.AddCommand(
		listCmd("list <form-id>", "List invited emails", 1, formResource("access_control_invitations"), "invitations", nil),
		accessControlInvitationsCreateCmd, accessControlInvitationsDeleteCmd,
	)
	formsCmd.AddCommand(accessControlInvitationsCmd)
}
