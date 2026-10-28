package commands

import (
	"net/url"
	"testing"
)

func TestSubmissionsListSendsDateFiltersAndOrder(t *testing.T) {
	got, _, err := runCLI(t, `{"submissions":[],"pagination":{}}`,
		"forms", "submissions", "list", "form-1",
		"--submitted-after", "2026-10-12", "--submitted-before", "2026-10-13T00:00:00+05:30", "--order", "desc")
	if err != nil {
		t.Fatal(err)
	}
	if got.path != "/api/external/v2/forms/form-1/submissions" {
		t.Errorf("path = %q", got.path)
	}
	want := url.Values{
		"submitted_after":  {"2026-10-12"},
		"submitted_before": {"2026-10-13T00:00:00+05:30"},
		"order":            {"desc"},
	}
	if got.query.Encode() != want.Encode() {
		t.Errorf("query = %v, want %v", got.query, want)
	}
}

func TestSubmissionsListOmitsUnsetFilters(t *testing.T) {
	got, _, err := runCLI(t, `{"submissions":[],"pagination":{}}`, "forms", "submissions", "list", "form-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(got.query) != 0 {
		t.Errorf("query = %v, want empty", got.query)
	}
}

func TestSubmissionsListRejectsInvalidOrder(t *testing.T) {
	got, _, err := runCLI(t, `{}`, "forms", "submissions", "list", "form-1", "--order", "newest")
	if err == nil || err.Error() != `Invalid --order "newest". Valid values: asc, desc.` {
		t.Fatalf("error = %v", err)
	}
	if got.path != "" {
		t.Errorf("request was made to %q", got.path)
	}
}

func TestSubmissionsListSendsStatusAndSearch(t *testing.T) {
	got, _, err := runCLI(t, `{"submissions":[],"pagination":{}}`,
		"forms", "submissions", "list", "form-1", "--status", "incomplete", "--search", "Oliver")
	if err != nil {
		t.Fatal(err)
	}
	want := url.Values{"status": {"incomplete"}, "search": {"Oliver"}}
	if got.query.Encode() != want.Encode() {
		t.Errorf("query = %v, want %v", got.query, want)
	}
}

func TestSubmissionsListRejectsInvalidStatus(t *testing.T) {
	got, _, err := runCLI(t, `{}`, "forms", "submissions", "list", "form-1", "--status", "done")
	if err == nil || err.Error() != `Invalid --status "done". Valid values: completed, incomplete, all.` {
		t.Fatalf("error = %v", err)
	}
	if got.path != "" {
		t.Errorf("request was made to %q", got.path)
	}
}

func TestSubmissionsShowGetsSubmission(t *testing.T) {
	got, _, err := runCLI(t, `{"id":"sub-1"}`, "forms", "submissions", "show", "form-1", "sub-1")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "GET" || got.path != "/api/external/v2/forms/form-1/submissions/sub-1" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
}

func TestSubmissionsDeleteDeletesSubmission(t *testing.T) {
	got, _, err := runCLI(t, `{"trashed_count":1}`, "forms", "submissions", "delete", "form-1", "sub-1")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "DELETE" || got.path != "/api/external/v2/forms/form-1/submissions/sub-1" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
}

func TestSubmissionsBulkDeleteSendsIDs(t *testing.T) {
	got, _, err := runCLI(t, `{"trashed_count":2}`, "forms", "submissions", "bulk-delete", "form-1", "--ids", "a,b")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "DELETE" || got.path != "/api/external/v2/forms/form-1/submissions/bulk_destroy" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
	if ids, _ := got.body["ids"].([]interface{}); len(ids) != 2 || ids[0] != "a" || ids[1] != "b" {
		t.Errorf("body = %v", got.body)
	}
}

func TestSubmissionsSendResponsesSendsMessage(t *testing.T) {
	got, _, err := runCLI(t, `{"queued":true}`, "forms", "submissions", "send-responses", "form-1",
		"--ids", "a,b", "--email-question-id", "q-1", "--subject", "Thanks", "--body", "<p>Hi</p>",
		"--cc-emails", "sam.smith@example.com", "--bcc-emails", "oliver.smith@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "POST" || got.path != "/api/external/v2/forms/form-1/submissions/send_responses" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
	b := got.body
	if b["email_question_id"] != "q-1" || b["subject"] != "Thanks" || b["body"] != "<p>Hi</p>" {
		t.Errorf("body = %v", b)
	}
	if cc, _ := b["cc_emails"].([]interface{}); len(cc) != 1 || cc[0] != "sam.smith@example.com" {
		t.Errorf("cc_emails = %v", b["cc_emails"])
	}
	if bcc, _ := b["bcc_emails"].([]interface{}); len(bcc) != 1 || bcc[0] != "oliver.smith@example.com" {
		t.Errorf("bcc_emails = %v", b["bcc_emails"])
	}
	if ids, _ := b["ids"].([]interface{}); len(ids) != 2 {
		t.Errorf("ids = %v", b["ids"])
	}
}

func TestSubmissionsExportSendsOptions(t *testing.T) {
	got, _, err := runCLI(t, `{"status":"queued"}`, "forms", "submissions", "export", "form-1",
		"--type", "xlsx", "--ids", "a", "--recipient-email", "sam.smith@example.com", "--include-meta-data")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "POST" || got.path != "/api/external/v2/forms/form-1/exports" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
	b := got.body
	if b["export_type"] != "xlsx" || b["recipient_email"] != "sam.smith@example.com" || b["include_meta_data"] != true {
		t.Errorf("body = %v", b)
	}
	if ids, _ := b["submission_ids"].([]interface{}); len(ids) != 1 || ids[0] != "a" {
		t.Errorf("submission_ids = %v", b["submission_ids"])
	}
	if _, ok := b["include_hidden_fields"]; ok {
		t.Errorf("unset flag was sent: %v", b)
	}
}

func TestSubmissionsExportRejectsInvalidType(t *testing.T) {
	got, _, err := runCLI(t, `{}`, "forms", "submissions", "export", "form-1", "--type", "docx")
	if err == nil || err.Error() != `Invalid --type "docx". Valid values: csv, xlsx, pdf.` {
		t.Fatalf("error = %v", err)
	}
	if got.path != "" {
		t.Errorf("request was made to %q", got.path)
	}
}
