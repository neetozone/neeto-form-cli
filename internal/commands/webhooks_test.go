package commands

import "testing"

func TestFormWebhooksCreateNestsParamsUnderWebhook(t *testing.T) {
	got, _, err := runCLI(t, `{"webhook":{"id":"w1"}}`, "forms", "webhooks", "create", "f1",
		"--endpoint", "https://example.com/hook", "--events", "submit", "--secret", "s3cret")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "POST" || got.path != "/api/external/v2/forms/f1/webhooks" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
	webhook := got.body["webhook"].(map[string]any)
	events := webhook["events"].([]any)
	if webhook["endpoint"] != "https://example.com/hook" || webhook["secret"] != "s3cret" || len(events) != 1 || events[0] != "submit" {
		t.Errorf("webhook = %v", webhook)
	}
}

func TestWebhooksCreateOmitsEventsWhenNotGiven(t *testing.T) {
	got, _, err := runCLI(t, `{"webhook":{"id":"w1"}}`, "webhooks", "create", "--endpoint", "https://example.com/hook")
	if err != nil {
		t.Fatal(err)
	}
	if got.path != "/api/external/v2/webhooks" {
		t.Errorf("path = %q", got.path)
	}
	webhook := got.body["webhook"].(map[string]any)
	if _, ok := webhook["events"]; ok {
		t.Errorf("webhook = %v, want no events", webhook)
	}
}

func TestWebhooksCreateRejectsUnknownEvent(t *testing.T) {
	got, _, err := runCLI(t, `{}`, "webhooks", "create", "--endpoint", "https://example.com/hook", "--events", "delete")
	if err == nil || err.Error() != `Invalid --events "delete". Valid values: submit.` {
		t.Fatalf("error = %v", err)
	}
	if got.path != "" {
		t.Errorf("request was sent to %q", got.path)
	}
}
