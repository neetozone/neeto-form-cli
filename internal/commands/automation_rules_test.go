package commands

import (
	"net/http"
	"reflect"
	"testing"
)

func TestAutomationRulesListSendsSearch(t *testing.T) {
	got, _, err := runCLI(t, `{"automation_rules":[],"pagination":{}}`, "forms", "automation-rules", "list", "form-1", "--search", "sales")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != http.MethodGet || got.path != "/api/external/v2/forms/form-1/automation_rules" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
	if got.query.Get("search") != "sales" {
		t.Errorf("query = %v", got.query)
	}
}

func TestAutomationRulesCreateTakesNameFromFile(t *testing.T) {
	file := writeJSONFile(t, `{"name":"From file"}`)
	got, _, err := runCLI(t, `{"automation_rule":{"id":"r1"}}`, "forms", "automation-rules", "create", "form-1", "--json-file", file)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]any{"automation_rule": map[string]any{"name": "From file"}}
	if !reflect.DeepEqual(got.body, want) {
		t.Errorf("body = %v", got.body)
	}
}

func TestAutomationRulesCreateMergesFlagsOverFile(t *testing.T) {
	file := writeJSONFile(t, `{"name":"From file","actions":[{"name":"email_to_form_creator"}],"execution_delay":{"value":3}}`)
	got, _, err := runCLI(t, `{"automation_rule":{"id":"r1"}}`, "forms", "automation-rules", "create", "form-1",
		"--json-file", file, "--name", "Email sales", "--status", "disabled", "--performer", "any",
		"--events", "submission_created, custom_field_updated", "--delay-unit", "hours")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != http.MethodPost || got.path != "/api/external/v2/forms/form-1/automation_rules" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
	want := map[string]any{"automation_rule": map[string]any{
		"name":            "Email sales",
		"status":          "disabled",
		"performer":       "any",
		"events":          []any{"submission_created", "custom_field_updated"},
		"actions":         []any{map[string]any{"name": "email_to_form_creator"}},
		"execution_delay": map[string]any{"value": float64(3), "unit": "hours"},
	}}
	if !reflect.DeepEqual(got.body, want) {
		t.Errorf("body = %v, want %v", got.body, want)
	}
}

func TestAutomationRulesCreateRequiresName(t *testing.T) {
	got, _, err := runCLI(t, `{}`, "forms", "automation-rules", "create", "form-1")
	if err == nil {
		t.Fatal("expected error")
	}
	if got.path != "" {
		t.Errorf("request was made to %q", got.path)
	}
}

func TestAutomationRulesUpdateSendsOnlyChangedFlags(t *testing.T) {
	got, _, err := runCLI(t, `{"automation_rule":{"id":"r1"}}`, "forms", "automation-rules", "update", "form-1", "r1", "--delay-value", "5", "--delay-unit", "minutes")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != http.MethodPatch || got.path != "/api/external/v2/forms/form-1/automation_rules/r1" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
	want := map[string]any{"automation_rule": map[string]any{"execution_delay": map[string]any{"value": float64(5), "unit": "minutes"}}}
	if !reflect.DeepEqual(got.body, want) {
		t.Errorf("body = %v, want %v", got.body, want)
	}
}

func TestAutomationRulesRejectsInvalidStatus(t *testing.T) {
	got, _, err := runCLI(t, `{}`, "forms", "automation-rules", "update", "form-1", "r1", "--status", "paused")
	if err == nil || err.Error() != `Invalid --status "paused". Valid values: enabled, disabled.` {
		t.Fatalf("error = %v", err)
	}
	if got.path != "" {
		t.Errorf("request was made to %q", got.path)
	}
}

func TestAutomationRulesRejectsInvalidDelayUnit(t *testing.T) {
	_, _, err := runCLI(t, `{}`, "forms", "automation-rules", "update", "form-1", "r1", "--delay-unit", "weeks")
	if err == nil || err.Error() != `Invalid --delay-unit "weeks". Valid values: minutes, hours, days.` {
		t.Fatalf("error = %v", err)
	}
}

func TestAutomationRulesRejectsInvalidEvent(t *testing.T) {
	_, _, err := runCLI(t, `{}`, "forms", "automation-rules", "update", "form-1", "r1", "--events", "submission_deleted")
	if err == nil || err.Error() != `Invalid --events "submission_deleted". Valid values: submission_created, custom_field_updated.` {
		t.Fatalf("error = %v", err)
	}
}

func TestAutomationRulesReorderSendsRuleIds(t *testing.T) {
	got, _, err := runCLI(t, `{"rule_ids":["r2","r1"]}`, "forms", "automation-rules", "reorder", "form-1", "--ids", "r2,r1")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != http.MethodPatch || got.path != "/api/external/v2/forms/form-1/automation_rules/reorder" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
	want := map[string]any{"rule_ids": []any{"r2", "r1"}}
	if !reflect.DeepEqual(got.body, want) {
		t.Errorf("body = %v, want %v", got.body, want)
	}
}
