package commands

import (
	"encoding/json"
	"testing"
)

const outlineResponse = `{"outline":"Riverside Roasters (published, single page)\n1. Work email (email, required)"}`

func TestFormsOutlinePrintsTheOutlineText(t *testing.T) {
	got, out, err := runCLI(t, outlineResponse, "forms", "outline", "form-1")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "GET" || got.path != "/api/external/v2/forms/form-1/outline" || got.query.Has("draft") {
		t.Errorf("request = %s %s %v", got.method, got.path, got.query)
	}
	var printed struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal([]byte(out), &printed); err != nil {
		t.Fatalf("output = %q: %v", out, err)
	}
	if printed.Message != "Riverside Roasters (published, single page)\n1. Work email (email, required)" {
		t.Errorf("message = %q", printed.Message)
	}
}

func TestFormsOutlineDraftAsksForUnpublishedChanges(t *testing.T) {
	got, _, err := runCLI(t, outlineResponse, "forms", "outline", "form-1", "--draft")
	if err != nil {
		t.Fatal(err)
	}
	if got.query.Get("draft") != "true" {
		t.Errorf("query = %v", got.query)
	}
}
