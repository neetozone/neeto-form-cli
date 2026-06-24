package output

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("io.Copy error: %v", err)
	}
	return buf.String()
}

func TestUseJSON_ForceJSON(t *testing.T) {
	ForceJSON = true
	QuietMode = false
	defer func() { ForceJSON = false }()

	if !UseJSON() {
		t.Error("UseJSON() = false, want true when ForceJSON is set")
	}
}

func TestUseJSON_QuietMode(t *testing.T) {
	ForceJSON = false
	QuietMode = true
	defer func() { QuietMode = false }()

	if !UseJSON() {
		t.Error("UseJSON() = false, want true when QuietMode is set")
	}
}

func TestPrintMessage_JSON(t *testing.T) {
	ForceJSON = true
	QuietMode = false
	defer func() { ForceJSON = false }()

	out := captureStdout(t, func() {
		PrintMessage("hello world")
	})

	var parsed map[string]string
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if parsed["message"] != "hello world" {
		t.Errorf("message = %q, want %q", parsed["message"], "hello world")
	}
}

func TestPrintMessage_Plain(t *testing.T) {
	ForceJSON = false
	QuietMode = false

	out := captureStdout(t, func() {
		PrintMessage("hello world")
	})

	trimmed := strings.TrimSpace(out)
	if !strings.Contains(trimmed, "hello world") {
		t.Errorf("output = %q, want it to contain %q", trimmed, "hello world")
	}
}

func TestPrint_QuietMode(t *testing.T) {
	ForceJSON = false
	QuietMode = true
	defer func() { QuietMode = false }()

	data := json.RawMessage(`[{"id":1}]`)
	out := captureStdout(t, func() {
		Print(data, nil)
	})

	trimmed := strings.TrimSpace(out)
	if trimmed != `[{"id":1}]` {
		t.Errorf("quiet output = %q, want raw data", trimmed)
	}
}

func TestPrint_JSONEnvelope(t *testing.T) {
	ForceJSON = true
	QuietMode = false
	defer func() { ForceJSON = false }()

	data := json.RawMessage(`{"name":"test"}`)
	breadcrumbs := []Breadcrumb{{Label: "details", Command: "app show 1"}}

	out := captureStdout(t, func() {
		Print(data, breadcrumbs)
	})

	var envelope Envelope
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &envelope); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	var parsed map[string]string
	if err := json.Unmarshal(envelope.Data, &parsed); err != nil {
		t.Fatalf("envelope data is not valid JSON: %v", err)
	}
	if parsed["name"] != "test" {
		t.Errorf("data.name = %q, want %q", parsed["name"], "test")
	}
	if len(envelope.Breadcrumbs) != 1 || envelope.Breadcrumbs[0].Label != "details" {
		t.Errorf("breadcrumbs = %v, want [{details app show 1}]", envelope.Breadcrumbs)
	}
}

func TestPrintWithPagination_QuietMode(t *testing.T) {
	ForceJSON = false
	QuietMode = true
	defer func() { QuietMode = false }()

	data := json.RawMessage(`[{"id":1}]`)
	pagination := json.RawMessage(`{"current_page_number":1,"total_pages":3,"total_records":25}`)

	out := captureStdout(t, func() {
		PrintWithPagination(data, pagination, nil)
	})

	trimmed := strings.TrimSpace(out)
	if trimmed != `[{"id":1}]` {
		t.Errorf("quiet output = %q, want raw data without pagination", trimmed)
	}
}

func TestPrintPretty_SubmissionsShowResponsesAndUserAgent(t *testing.T) {
	data := json.RawMessage(`[
		{
			"id": "3ad844c3-2785-4ddf-87a4-b48195fc6360",
			"created_at": "2026-02-09T05:20:59.464Z",
			"field_values": [{"id": "fv1", "value": "x"}],
			"user_agent": {"name": "Chrome", "operating_system": "macOS", "ip_address": "1.2.3.4"},
			"responses": [
				{"id": "r1", "label": "Email", "kind": "email", "value": "foo@bar.com"},
				{"id": "r2", "label": "Full Name", "kind": "text", "value": "John Doe"}
			]
		}
	]`)

	out := captureStdout(t, func() {
		printPretty(data)
	})

	for _, want := range []string{
		"3ad844c3-2785-4ddf-87a4-b48195fc6360",
		"RESPONSES", "Email", "foo@bar.com", "Full Name", "John Doe",
		"USER AGENT", "Chrome", "macOS", "1.2.3.4",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("pretty submissions output is missing %q.\nGot:\n%s", want, out)
		}
	}
}

func TestPrintPretty_FormsStillRenderAsTable(t *testing.T) {
	data := json.RawMessage(`[
		{"id": "f1", "title": "Contact form", "state": "published", "submissions_count": 3, "created_at": "2026-02-09T05:20:59.464Z"},
		{"id": "f2", "title": "Survey", "state": "draft", "submissions_count": 0, "created_at": "2026-02-10T05:20:59.464Z"}
	]`)

	out := captureStdout(t, func() {
		printPretty(data)
	})

	if !strings.Contains(out, "TITLE") {
		t.Errorf("forms list should keep the table header TITLE.\nGot:\n%s", out)
	}
	if !strings.Contains(out, "Contact form") || !strings.Contains(out, "Survey") {
		t.Errorf("forms list should show both rows.\nGot:\n%s", out)
	}
}

func TestPrintWithPagination_JSONEnvelope(t *testing.T) {
	ForceJSON = true
	QuietMode = false
	defer func() { ForceJSON = false }()

	data := json.RawMessage(`[{"id":1}]`)
	pagination := json.RawMessage(`{"current_page_number":1,"total_pages":3}`)

	out := captureStdout(t, func() {
		PrintWithPagination(data, pagination, nil)
	})

	var envelope Envelope
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &envelope); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if envelope.Pagination == nil {
		t.Error("pagination should be present in envelope")
	}
}
