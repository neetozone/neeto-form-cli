package commands

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func writeQuestionFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "q.json")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestQuestionsListSendsPagination(t *testing.T) {
	got, _, err := runCLI(t, `{"questions":[],"pagination":{}}`, "forms", "questions", "list", "f-1", "--page", "2", "--page-size", "5")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "GET" || got.path != "/api/external/v2/forms/f-1/questions" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
	if got.query.Get("page") != "2" || got.query.Get("page_size") != "5" {
		t.Errorf("query = %v", got.query)
	}
}

func TestQuestionsCreateMergesFlagsOverFile(t *testing.T) {
	file := writeQuestionFile(t, `{"label":"Old","options":[{"label":"Basic"}]}`)
	got, _, err := runCLI(t, `{"question":{}}`, "forms", "questions", "create", "f-1", "--kind", "dropdown", "--position", "2", "--label", "Plan", "--required", "--json-file", file)
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "POST" || got.path != "/api/external/v2/forms/f-1/questions" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
	want := map[string]any{"question": map[string]any{
		"kind":     "dropdown",
		"position": float64(2),
		"attributes": map[string]any{
			"label":       "Plan",
			"is_required": true,
			"options":     []any{map[string]any{"label": "Basic"}},
		},
	}}
	if !reflect.DeepEqual(got.body, want) {
		t.Errorf("body = %v, want %v", got.body, want)
	}
}

func TestQuestionsBulkCreateSendsFileBody(t *testing.T) {
	file := writeQuestionFile(t, `{"questions":[{"kind":"phone"}]}`)
	got, _, err := runCLI(t, `{"questions":[]}`, "forms", "questions", "bulk-create", "f-1", "--json-file", file)
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "POST" || got.path != "/api/external/v2/forms/f-1/questions/bulk_create" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
	want := map[string]any{"questions": []any{map[string]any{"kind": "phone"}}}
	if !reflect.DeepEqual(got.body, want) {
		t.Errorf("body = %v, want %v", got.body, want)
	}
}

func TestQuestionsUpdateSendsAttributes(t *testing.T) {
	got, _, err := runCLI(t, `{"question":{}}`, "forms", "questions", "update", "f-1", "q-1", "--required=false", "--placeholder", "you@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "PATCH" || got.path != "/api/external/v2/forms/f-1/questions/q-1" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
	want := map[string]any{"question": map[string]any{"attributes": map[string]any{"is_required": false, "placeholder": "you@example.com"}}}
	if !reflect.DeepEqual(got.body, want) {
		t.Errorf("body = %v, want %v", got.body, want)
	}
}

func TestQuestionsCloneSendsPost(t *testing.T) {
	got, _, err := runCLI(t, `{"question":{}}`, "forms", "questions", "clone", "f-1", "q-1")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "POST" || got.path != "/api/external/v2/forms/f-1/questions/q-1/clone" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
}

func TestQuestionsReorderSendsIds(t *testing.T) {
	got, _, err := runCLI(t, `{"question_ids":[]}`, "forms", "questions", "reorder", "f-1", "--ids", "a, b,c")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "PATCH" || got.path != "/api/external/v2/forms/f-1/questions/reorder" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
	want := map[string]any{"question_ids": []any{"a", "b", "c"}}
	if !reflect.DeepEqual(got.body, want) {
		t.Errorf("body = %v, want %v", got.body, want)
	}
}

func TestQuestionsDeleteSendsDelete(t *testing.T) {
	got, _, err := runCLI(t, `{}`, "forms", "questions", "delete", "f-1", "q-1")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "DELETE" || got.path != "/api/external/v2/forms/f-1/questions/q-1" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
}
