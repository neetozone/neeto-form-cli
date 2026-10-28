package commands

import "github.com/spf13/cobra"

var ipRestrictionsCmd = &cobra.Command{
	Use:   "ip-restrictions",
	Short: "View IP restrictions",
}

var ipRestrictionsShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show IP restriction status and allowed ranges",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := fetch(cmd, apiPath("ip_restriction"), nil)
		if err != nil {
			return err
		}
		printResource(data, nil)
		return nil
	},
}

func init() {
	ipRestrictionsCmd.AddCommand(ipRestrictionsShowCmd)
	register(func(root *cobra.Command) { root.AddCommand(ipRestrictionsCmd) })
}
