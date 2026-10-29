package commands

import (
	"net/url"

	"github.com/spf13/cobra"
)

var formsOutlineCmd = &cobra.Command{
	Use:   "outline <form-id>",
	Short: "Show a text outline of a form as respondents see it: theme, questions, logic and thank-you page",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		params := url.Values{}
		if draft, _ := cmd.Flags().GetBool("draft"); draft {
			params.Set("draft", "true")
		}
		data, err := fetch(cmd, formPath(args[0], "outline"), params)
		if err != nil {
			return err
		}
		printResource(data, nil)
		return nil
	},
}

func init() {
	formsOutlineCmd.Flags().Bool("draft", false, "Include unpublished changes")
	formsCmd.AddCommand(formsOutlineCmd)
}
