package commands

import "github.com/spf13/cobra"

var ipRestrictionsCmd = &cobra.Command{
	Use:   "ip-restrictions",
	Short: "View IP restrictions",
}

func init() {
	ipRestrictionsCmd.AddCommand(showCmd("show", "Show IP restriction status and allowed ranges", 0, workspaceResource("ip_restriction")))
	register(func(root *cobra.Command) { root.AddCommand(ipRestrictionsCmd) })
}
