package commands

import (
	"os"
	"path/filepath"
	"testing"
)

func TestThemesCreateSendsNameBaseAndProperties(t *testing.T) {
	got, _, err := runCLI(t, `{"theme":{"id":"t1"}}`, "themes", "create", "--name", "Brand", "--base-theme-id", "b1",
		"--property", "question_color=#123456", "--property", "background_color=#ffffff")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "POST" || got.path != "/api/external/v2/themes" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
	theme := got.body["theme"].(map[string]any)
	properties := theme["properties"].(map[string]any)
	if theme["name"] != "Brand" || theme["base_theme_id"] != "b1" {
		t.Errorf("theme = %v", theme)
	}
	if properties["question_color"] != "#123456" || properties["background_color"] != "#ffffff" {
		t.Errorf("properties = %v", properties)
	}
}

func TestThemesUpdateMergesPropertyFlagsOverJSONFile(t *testing.T) {
	file := filepath.Join(t.TempDir(), "theme.json")
	if err := os.WriteFile(file, []byte(`{"name":"Old","properties":{"question_color":"#111111","button_border_radius":"4"}}`), 0o600); err != nil {
		t.Fatal(err)
	}

	got, _, err := runCLI(t, `{"theme":{"id":"t1"}}`, "themes", "update", "t1", "--json-file", file,
		"--name", "New", "--property", "question_color=#222222")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "PATCH" || got.path != "/api/external/v2/themes/t1" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
	theme := got.body["theme"].(map[string]any)
	properties := theme["properties"].(map[string]any)
	if theme["name"] != "New" || properties["question_color"] != "#222222" || properties["button_border_radius"] != "4" {
		t.Errorf("theme = %v", theme)
	}
}

func TestThemesCreateRejectsPropertyWithoutEquals(t *testing.T) {
	got, _, err := runCLI(t, `{}`, "themes", "create", "--name", "Brand", "--property", "question_color")
	if err == nil || err.Error() != `Invalid --property "question_color". Use key=value.` {
		t.Fatalf("error = %v", err)
	}
	if got.path != "" {
		t.Errorf("request was sent to %q", got.path)
	}
}

func TestFormThemeUpdateAppliesTheme(t *testing.T) {
	got, _, err := runCLI(t, `{"theme":{"id":"t1"}}`, "forms", "theme", "update", "f1", "--theme-id", "t1")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "PUT" || got.path != "/api/external/v2/forms/f1/theme" || got.body["theme_id"] != "t1" {
		t.Errorf("request = %s %s %v", got.method, got.path, got.body)
	}
}
