package commands

import "github.com/spf13/cobra"

var rolesCmd = &cobra.Command{
	Use:   "roles",
	Short: "View organization roles",
}

func init() {
	rolesCmd.AddCommand(listCmd("list", "List organization roles", 0, workspaceResource("roles"), "roles", nil))
	register(func(root *cobra.Command) { root.AddCommand(rolesCmd) })
}
