package commands

import "testing"

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
	want := map[string]string{
		"submitted_after":  "2026-10-12",
		"submitted_before": "2026-10-13T00:00:00+05:30",
		"order":            "desc",
	}
	for key, value := range want {
		if got.query.Get(key) != value {
			t.Errorf("%s = %q, want %q", key, got.query.Get(key), value)
		}
	}
}

func TestSubmissionsListOmitsUnsetFilters(t *testing.T) {
	got, _, err := runCLI(t, `{"submissions":[],"pagination":{}}`, "forms", "submissions", "list", "form-1")
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"submitted_after", "submitted_before", "order"} {
		if got.query.Has(key) {
			t.Errorf("unexpected %s in %v", key, got.query)
		}
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
