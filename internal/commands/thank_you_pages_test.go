package commands

import "testing"

func TestThankYouPagesCreateSendsRedirectFields(t *testing.T) {
	got, _, err := runCLI(t, `{"thank_you_page":{}}`, "forms", "thank-you-pages", "create", "f1",
		"--kind", "redirect_to_url", "--redirect-url", "https://example.com/pricing", "--pass-event-details", "--social-sharing=false", "--primary")
	if err != nil {
		t.Fatal(err)
	}
	inner := got.body["thank_you_page"].(map[string]interface{})
	if got.method != "POST" || got.path != "/api/external/v2/forms/f1/thank_you_pages" ||
		inner["kind"] != "redirect_to_url" || inner["redirect_url"] != "https://example.com/pricing" ||
		inner["pass_event_details_to_redirect_url"] != true || inner["social_sharing_enabled"] != false || inner["primary"] != true {
		t.Errorf("request = %+v", got)
	}
}

func TestThankYouPagesUpdateSendsOnlyPassedFields(t *testing.T) {
	got, _, err := runCLI(t, `{"thank_you_page":{}}`, "forms", "thank-you-pages", "update", "f1", "p1", "--title", "Pricing")
	if err != nil {
		t.Fatal(err)
	}
	inner := got.body["thank_you_page"].(map[string]interface{})
	if got.method != "PATCH" || got.path != "/api/external/v2/forms/f1/thank_you_pages/p1" || len(inner) != 1 || inner["title"] != "Pricing" {
		t.Errorf("request = %+v", got)
	}
}

func TestThankYouPagesCreateRejectsUnknownKind(t *testing.T) {
	got, _, err := runCLI(t, `{}`, "forms", "thank-you-pages", "create", "f1", "--kind", "popup")
	if err == nil || err.Error() != `Invalid --kind "popup". Valid values: custom_message, redirect_to_url.` {
		t.Fatalf("error = %v", err)
	}
	if got.path != "" {
		t.Errorf("request was made to %q", got.path)
	}
}

func TestThankYouPagesDeleteSendsDeleteToPage(t *testing.T) {
	got, _, err := runCLI(t, `{"id":"p1","deleted":true}`, "forms", "thank-you-pages", "delete", "f1", "p1")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "DELETE" || got.path != "/api/external/v2/forms/f1/thank_you_pages/p1" {
		t.Errorf("request = %+v", got)
	}
}
