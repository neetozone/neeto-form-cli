package commands

import "github.com/spf13/cobra"

var rolesCmd = &cobra.Command{
	Use:   "roles",
	Short: "View organization roles",
}

var rolesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List organization roles",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := fetch(cmd, apiPath("roles"), paginationParams(cmd))
		if err != nil {
			return err
		}
		printList(data, "roles", nil)
		return nil
	},
}

func init() {
	addPaginationFlags(rolesListCmd)
	rolesCmd.AddCommand(rolesListCmd)
	register(func(root *cobra.Command) { root.AddCommand(rolesCmd) })
}
