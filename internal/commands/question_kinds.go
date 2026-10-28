package commands

import "github.com/spf13/cobra"

var questionKindsCmd = &cobra.Command{
	Use:   "question-kinds",
	Short: "View question kinds",
}

var questionKindsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List question kinds",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := fetch(cmd, apiPath("question_kinds"), nil)
		if err != nil {
			return err
		}

		printList(data, "kinds", nil)
		return nil
	},
}

func init() {
	register(func(root *cobra.Command) { root.AddCommand(questionKindsCmd) })

	questionKindsCmd.AddCommand(questionKindsListCmd)
}
