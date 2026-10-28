package commands

import "github.com/spf13/cobra"

var auditLogsCmd = &cobra.Command{
	Use:   "audit-logs",
	Short: "View the workspace audit log",
}

var auditLogsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List audit log entries",
	RunE: func(cmd *cobra.Command, args []string) error {
		params := queryValues(cmd, paginationParams(cmd), "created-after", "created-before")
		data, err := fetch(cmd, apiPath("audit_logs"), params)
		if err != nil {
			return err
		}
		printList(data, "audit_logs", nil)
		return nil
	},
}

func init() {
	addPaginationFlags(auditLogsListCmd)
	auditLogsListCmd.Flags().String("created-after", "", "Only entries created on or after this ISO 8601 date or date-time")
	auditLogsListCmd.Flags().String("created-before", "", "Only entries created before this ISO 8601 date or date-time")
	auditLogsCmd.AddCommand(auditLogsListCmd)
	register(func(root *cobra.Command) { root.AddCommand(auditLogsCmd) })
}
