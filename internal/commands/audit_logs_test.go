package commands

import "testing"

func TestAuditLogsListSendsDateRange(t *testing.T) {
	got, _, err := runCLI(t, `{"audit_logs":[],"pagination":{}}`, "audit-logs", "list", "--created-after", "2026-10-01", "--created-before", "2026-10-28")
	if err != nil {
		t.Fatal(err)
	}
	if got.path != "/api/external/v2/audit_logs" {
		t.Errorf("path = %q", got.path)
	}
	if got.query.Get("created_after") != "2026-10-01" || got.query.Get("created_before") != "2026-10-28" {
		t.Errorf("query = %v", got.query)
	}
}
