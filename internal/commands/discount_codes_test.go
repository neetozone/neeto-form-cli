package commands

import (
	"encoding/json"
	"testing"
)

func discountCodesJSON(t *testing.T, v any) string {
	t.Helper()
	out, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return string(out)
}

func TestDiscountCodesCreateSendsNestedBody(t *testing.T) {
	got, _, err := runCLI(t, `{"discount_code":{"id":"d1"}}`, "discount-codes", "create",
		"--code", "WELCOME", "--value", "15", "--kind", "fixed", "--currency", "usd",
		"--expires-at", "2026-12-31", "--redeems-limit", "50", "--form-ids", "f1, f2")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "POST" || got.path != "/api/external/v2/discount_codes" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
	want := `{"code":"WELCOME","currency":"usd","expires_at":"2026-12-31","form_ids":["f1","f2"],"kind":"fixed","redeems_limit":50,"value":"15"}`
	if s := discountCodesJSON(t, got.body["discount_code"]); s != want {
		t.Errorf("discount_code = %s, want %s", s, want)
	}
}

func TestDiscountCodesUpdateWithoutFormIDsOmitsThem(t *testing.T) {
	got, _, err := runCLI(t, `{"discount_code":{"id":"d1"}}`, "discount-codes", "update", "d1", "--value", "20")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "PATCH" || got.path != "/api/external/v2/discount_codes/d1" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
	if s := discountCodesJSON(t, got.body["discount_code"]); s != `{"value":"20"}` {
		t.Errorf("discount_code = %s", s)
	}
}

func TestDiscountCodesCreateRejectsInvalidKind(t *testing.T) {
	got, _, err := runCLI(t, `{}`, "discount-codes", "create", "--code", "WELCOME", "--value", "5", "--kind", "flat")
	if err == nil || err.Error() != `Invalid --kind "flat". Valid values: percentage, fixed.` {
		t.Fatalf("error = %v", err)
	}
	if got.path != "" {
		t.Errorf("request was made to %q", got.path)
	}
}

func TestDiscountCodesReportSendsDateRange(t *testing.T) {
	got, _, err := runCLI(t, `{"report":{}}`, "discount-codes", "report", "d1", "--start-date", "2026-10-01", "--end-date", "2026-10-31")
	if err != nil {
		t.Fatal(err)
	}
	if got.path != "/api/external/v2/discount_codes/d1/report" {
		t.Errorf("path = %q", got.path)
	}
	if got.query.Get("start_date") != "2026-10-01" || got.query.Get("end_date") != "2026-10-31" {
		t.Errorf("query = %v", got.query)
	}
}

func TestDiscountCodesFormsListUsesCollectionPath(t *testing.T) {
	got, _, err := runCLI(t, `{"forms":[],"pagination":{}}`, "discount-codes", "forms", "list")
	if err != nil {
		t.Fatal(err)
	}
	if got.path != "/api/external/v2/discount_codes/forms" {
		t.Errorf("path = %q", got.path)
	}
}

func TestDiscountCodesRedemptionsListUsesNestedPath(t *testing.T) {
	got, _, err := runCLI(t, `{"redemptions":[],"pagination":{}}`, "discount-codes", "redemptions", "list", "d1")
	if err != nil {
		t.Fatal(err)
	}
	if got.path != "/api/external/v2/discount_codes/d1/redemptions" {
		t.Errorf("path = %q", got.path)
	}
}

func TestDiscountCodesDeleteSendsDelete(t *testing.T) {
	got, _, err := runCLI(t, `{}`, "discount-codes", "delete", "d1")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "DELETE" || got.path != "/api/external/v2/discount_codes/d1" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
}
