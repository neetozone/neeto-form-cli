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
