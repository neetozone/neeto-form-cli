package commands

func init() {
	insightsCmd := showCmd("insights <form-id>", "Show form insights", 1, formResource("insights"), "start-date", "end-date", "time-zone")
	insightsCmd.Flags().String("start-date", "", "First day to include (YYYY-MM-DD); use with --end-date")
	insightsCmd.Flags().String("end-date", "", "Last day to include (YYYY-MM-DD); use with --start-date")
	insightsCmd.Flags().String("time-zone", "", "IANA time zone such as Asia/Kolkata (default: your time zone)")
	formsCmd.AddCommand(insightsCmd)
}
