package commands

import (
	"net/url"

	"github.com/spf13/cobra"
)

var insightsCmd = &cobra.Command{
	Use:   "insights <form-id>",
	Short: "Show form insights",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		params := queryValues(cmd, url.Values{}, "start-date", "end-date", "time-zone")
		data, err := fetch(cmd, formPath(args[0], "insights"), params)
		if err != nil {
			return err
		}
		printResource(data, nil)
		return nil
	},
}

func init() {
	insightsCmd.Flags().String("start-date", "", "First day to include (YYYY-MM-DD); use with --end-date")
	insightsCmd.Flags().String("end-date", "", "Last day to include (YYYY-MM-DD); use with --start-date")
	insightsCmd.Flags().String("time-zone", "", "IANA time zone such as Asia/Kolkata (default: your time zone)")
	formsCmd.AddCommand(insightsCmd)
}
