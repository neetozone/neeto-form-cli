package commands

import (
	"net/http"

	"github.com/spf13/cobra"
)

var accessControlInvitationsCmd = &cobra.Command{
	Use:   "access-control-invitations",
	Short: "Manage who is invited to a form",
}

var accessControlInvitationsListCmd = &cobra.Command{
	Use:   "list <form-id>",
	Short: "List invited emails",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := fetch(cmd, formPath(args[0], "access_control_invitations"), paginationParams(cmd))
		if err != nil {
			return err
		}
		printList(data, "invitations", nil)
		return nil
	},
}

var accessControlInvitationsCreateCmd = &cobra.Command{
	Use:   "create <form-id>",
	Short: "Invite emails to a form",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := send(cmd, http.MethodPost, formPath(args[0], "access_control_invitations"), flagValues(cmd, "emails[]"))
		if err != nil {
			return err
		}
		printActionResult(data, nil)
		return nil
	},
}

var accessControlInvitationsDeleteCmd = &cobra.Command{
	Use:   "delete <form-id>",
	Short: "Remove invited emails from a form",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := send(cmd, http.MethodDelete, formPath(args[0], "access_control_invitations"), flagValues(cmd, "emails[]"))
		if err != nil {
			return err
		}
		printActionResult(data, nil)
		return nil
	},
}

func init() {
	addPaginationFlags(accessControlInvitationsListCmd)
	accessControlInvitationsCreateCmd.Flags().String("emails", "", "Comma-separated emails to invite, 1 to 100. The form access must be set to invited")
	accessControlInvitationsDeleteCmd.Flags().String("emails", "", "Comma-separated emails to remove, 1 to 100")
	markFlagsRequired(accessControlInvitationsCreateCmd, "emails")
	markFlagsRequired(accessControlInvitationsDeleteCmd, "emails")

	accessControlInvitationsCmd.AddCommand(accessControlInvitationsListCmd, accessControlInvitationsCreateCmd, accessControlInvitationsDeleteCmd)
	formsCmd.AddCommand(accessControlInvitationsCmd)
}
