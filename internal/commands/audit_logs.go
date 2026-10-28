package commands

import "github.com/spf13/cobra"

var auditLogsCmd = &cobra.Command{
	Use:   "audit-logs",
	Short: "View the workspace audit log",
}

func init() {
	list := listCmd("list", "List audit log entries", 0, workspaceResource("audit_logs"), "audit_logs", nil, "created-after", "created-before")
	list.Flags().String("created-after", "", "Only entries created on or after this ISO 8601 date or date-time")
	list.Flags().String("created-before", "", "Only entries created before this ISO 8601 date or date-time")
	auditLogsCmd.AddCommand(list)
	register(func(root *cobra.Command) { root.AddCommand(auditLogsCmd) })
}
