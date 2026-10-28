package commands

import (
	"reflect"
	"strings"
	"testing"
)

func TestShareShowRequestsShareDetails(t *testing.T) {
	got, _, err := runCLI(t, `{"id":"f1","public_link":"https://acme.neetoform.com/feedback"}`, "forms", "share", "show", "f1")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "GET" || got.path != "/api/external/v2/forms/f1/share" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
}

func TestShareShowPrintsOnlyChosenEmbed(t *testing.T) {
	_, out, err := runCLI(t, `{"embed_codes":{"floating_popup":"<div>popup</div>","iframe":"<iframe></iframe>"}}`, "forms", "share", "show", "f1", "--embed", "floating-popup")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `\u003cdiv\u003epopup`) || strings.Contains(out, "iframe") {
		t.Errorf("output = %q", out)
	}
}

func TestShareShowEmbedWithQuietPrintsOnlyTheSnippet(t *testing.T) {
	_, out, err := runCLI(t, `{"embed_codes":{"iframe":"<iframe></iframe>"}}`, "forms", "share", "show", "f1", "--embed", "iframe", "--quiet")
	if err != nil {
		t.Fatal(err)
	}
	if out != "<iframe></iframe>\n" {
		t.Errorf("output = %q", out)
	}
}

func TestShareShowRejectsInvalidEmbed(t *testing.T) {
	got, _, err := runCLI(t, `{}`, "forms", "share", "show", "f1", "--embed", "banner")
	if err == nil || err.Error() != `Invalid --embed "banner". Valid values: inline, floating-popup, element-popup, iframe.` {
		t.Fatalf("error = %v", err)
	}
	if got.path != "" {
		t.Errorf("request was sent to %q", got.path)
	}
}

func TestPublicLinkUpdateSendsPublicLinkID(t *testing.T) {
	got, _, err := runCLI(t, `{"id":"f1"}`, "forms", "public-link", "update", "f1", "--public-link-id", "customer-feedback")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "PATCH" || got.path != "/api/external/v2/forms/f1/public_link" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
	if got.body["public_link_id"] != "customer-feedback" {
		t.Errorf("body = %v", got.body)
	}
}

func TestPublicLinkUpdateWithoutFlagSendsNoPublicLinkID(t *testing.T) {
	got, _, err := runCLI(t, `{"id":"f1"}`, "forms", "public-link", "update", "f1")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := got.body["public_link_id"]; ok {
		t.Errorf("body = %v, want no public_link_id", got.body)
	}
}

func TestPrefilledLinksCreateSendsAnswersAndDiscountCode(t *testing.T) {
	got, _, err := runCLI(t, `{"url":"https://acme.neetoform.com/feedback?email=oliver%40example.com"}`,
		"forms", "prefilled-links", "create", "f1",
		"--answer", "q1=oliver@example.com",
		"--answer", "q2=a=b",
		"--answers-json", `[{"question_id":"q3","value":{"country":"US","number":"2025550123"}}]`,
		"--discount-code", "SAVE10")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "POST" || got.path != "/api/external/v2/forms/f1/prefilled_link" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
	want := []any{
		map[string]any{"question_id": "q1", "value": "oliver@example.com"},
		map[string]any{"question_id": "q2", "value": "a=b"},
		map[string]any{"question_id": "q3", "value": map[string]any{"country": "US", "number": "2025550123"}},
	}
	if !reflect.DeepEqual(got.body["answers"], want) {
		t.Errorf("answers = %v", got.body["answers"])
	}
	if got.body["discount_code"] != "SAVE10" {
		t.Errorf("body = %v", got.body)
	}
}

func TestPrefilledLinksCreateRejectsMalformedAnswer(t *testing.T) {
	got, _, err := runCLI(t, `{}`, "forms", "prefilled-links", "create", "f1", "--answer", "q1")
	if err == nil || err.Error() != `Invalid --answer "q1". Use <question-id>=<value>.` {
		t.Fatalf("error = %v", err)
	}
	if got.path != "" {
		t.Errorf("request was sent to %q", got.path)
	}
}

func TestEmailInvitationsCreateSendsEmailsAndOptions(t *testing.T) {
	got, _, err := runCLI(t, `{"form_id":"f1"}`,
		"forms", "email-invitations", "create", "f1",
		"--emails", "oliver@example.com, sam@example.com",
		"--subject", "Please reply",
		"--body", "<p>Hi</p>",
		"--button-text", "Start",
		"--reply-to", "reply@example.com",
		"--from", "sender@example.com",
		"--send-at", "2026-11-01T10:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "POST" || got.path != "/api/external/v2/forms/f1/email_invitations" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
	want := map[string]any{
		"emails":                   []any{"oliver@example.com", "sam@example.com"},
		"email_subject":            "Please reply",
		"email_body":               "<p>Hi</p>",
		"start_survey_button_text": "Start",
		"reply_to_email":           "reply@example.com",
		"send_from_email":          "sender@example.com",
		"send_at":                  "2026-11-01T10:00:00Z",
	}
	if !reflect.DeepEqual(got.body, want) {
		t.Errorf("body = %v", got.body)
	}
}

func TestEmailInvitationsCreateRequiresEmails(t *testing.T) {
	got, _, err := runCLI(t, `{}`, "forms", "email-invitations", "create", "f1")
	if err == nil {
		t.Fatal("expected error")
	}
	if got.path != "" {
		t.Errorf("request was sent to %q", got.path)
	}
}
