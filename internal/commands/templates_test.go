package commands

import "testing"

func TestTemplatesListSendsSearchAndCategory(t *testing.T) {
	got, _, err := runCLI(t, `{"templates":[],"pagination":{}}`, "templates", "list", "--search", "feedback", "--category", "Events")
	if err != nil {
		t.Fatal(err)
	}
	if got.path != "/api/external/v2/templates" {
		t.Errorf("path = %q", got.path)
	}
	if got.query.Get("search") != "feedback" || got.query.Get("category") != "Events" {
		t.Errorf("query = %v", got.query)
	}
}

func TestTemplatesUsePostsCreatedByEmail(t *testing.T) {
	got, _, err := runCLI(t, `{"form":{"id":"f1"}}`, "templates", "use", "customer-feedback", "--created-by-email", "oliver.smith@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "POST" || got.path != "/api/external/v2/templates/customer-feedback/forms" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
	if got.body["created_by_email"] != "oliver.smith@example.com" {
		t.Errorf("body = %v", got.body)
	}
}

func TestTemplatesUseWithoutEmailSendsEmptyBody(t *testing.T) {
	got, _, err := runCLI(t, `{"form":{"id":"f1"}}`, "templates", "use", "t1")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := got.body["created_by_email"]; ok {
		t.Errorf("body = %v, want no created_by_email", got.body)
	}
}
