package commands

import (
	"net/url"
	"testing"
)

func TestInsightsSendsRangeAndTimeZone(t *testing.T) {
	got, _, err := runCLI(t, `{"visit_count":120}`, "forms", "insights", "form-1",
		"--start-date", "2026-10-01", "--end-date", "2026-10-31", "--time-zone", "Asia/Kolkata")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "GET" || got.path != "/api/external/v2/forms/form-1/insights" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
	want := url.Values{"start_date": {"2026-10-01"}, "end_date": {"2026-10-31"}, "time_zone": {"Asia/Kolkata"}}
	if got.query.Encode() != want.Encode() {
		t.Errorf("query = %v, want %v", got.query, want)
	}
}
