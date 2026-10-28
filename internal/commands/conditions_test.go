package commands

import (
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func writeJSONFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "payload.json")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestConditionsListSendsPagination(t *testing.T) {
	got, _, err := runCLI(t, `{"conditions":[],"pagination":{}}`, "forms", "conditions", "list", "form-1", "--page", "2", "--page-size", "10")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != http.MethodGet || got.path != "/api/external/v2/forms/form-1/conditions" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
	if got.query.Get("page") != "2" || got.query.Get("page_size") != "10" {
		t.Errorf("query = %v", got.query)
	}
}

func TestConditionsCreateWrapsFileAndEnabledFlag(t *testing.T) {
	file := writeJSONFile(t, `{"actions":{"items":[{"action_type":"hide","fields":["q1"]}]},"enabled":true}`)
	got, _, err := runCLI(t, `{"condition":{"id":"c1"}}`, "forms", "conditions", "create", "form-1", "--json-file", file, "--enabled=false")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != http.MethodPost || got.path != "/api/external/v2/forms/form-1/conditions" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
	want := map[string]any{"condition": map[string]any{
		"enabled": false,
		"actions": map[string]any{"items": []any{map[string]any{"action_type": "hide", "fields": []any{"q1"}}}},
	}}
	if !reflect.DeepEqual(got.body, want) {
		t.Errorf("body = %v, want %v", got.body, want)
	}
}

func TestConditionsUpdatePatchesCondition(t *testing.T) {
	got, _, err := runCLI(t, `{"condition":{"id":"c1"}}`, "forms", "conditions", "update", "form-1", "c1", "--enabled=false")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != http.MethodPatch || got.path != "/api/external/v2/forms/form-1/conditions/c1" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
	want := map[string]any{"condition": map[string]any{"enabled": false}}
	if !reflect.DeepEqual(got.body, want) {
		t.Errorf("body = %v, want %v", got.body, want)
	}
}

func TestConditionsCloneSendsPostToCloneMember(t *testing.T) {
	got, _, err := runCLI(t, `{"condition":{"id":"c2"}}`, "forms", "conditions", "clone", "form-1", "c1")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != http.MethodPost || got.path != "/api/external/v2/forms/form-1/conditions/c1/clone" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
}

func TestConditionsReorderSendsIds(t *testing.T) {
	got, _, err := runCLI(t, `{"condition_ids":["c2","c1"]}`, "forms", "conditions", "reorder", "form-1", "--ids", "c2, c1")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != http.MethodPatch || got.path != "/api/external/v2/forms/form-1/conditions/reorder" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
	want := map[string]any{"condition_ids": []any{"c2", "c1"}}
	if !reflect.DeepEqual(got.body, want) {
		t.Errorf("body = %v, want %v", got.body, want)
	}
}

func TestConditionsDeleteSendsDelete(t *testing.T) {
	got, _, err := runCLI(t, `{"id":"c1","deleted":true}`, "forms", "conditions", "delete", "form-1", "c1")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != http.MethodDelete || got.path != "/api/external/v2/forms/form-1/conditions/c1" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
}
