package commands

import "testing"

func TestEmailNotificationsUpdateSendsAdminRecipients(t *testing.T) {
	got, _, err := runCLI(t, `{"email_notification":{}}`, "forms", "email-notifications", "update", "f1", "admin",
		"--enabled=false", "--subject", "New submission", "--to", "oliver.smith@example.com,sam.smith@example.com", "--cc", "sam.smith@example.com")
	if err != nil {
		t.Fatal(err)
	}
	inner := got.body["email_notification"].(map[string]interface{})
	to := inner["notify_emails"].([]interface{})
	cc := inner["notify_cc_emails"].([]interface{})
	if got.method != "PATCH" || got.path != "/api/external/v2/forms/f1/email_notifications/admin" ||
		inner["is_enabled"] != false || inner["subject"] != "New submission" ||
		len(to) != 2 || to[1] != "sam.smith@example.com" || len(cc) != 1 {
		t.Errorf("request = %+v", got)
	}
}

func TestEmailNotificationsUpdateSendsRespondentFields(t *testing.T) {
	got, _, err := runCLI(t, `{"email_notification":{}}`, "forms", "email-notifications", "update", "f1", "respondent",
		"--reply-to", "oliver.smith@example.com", "--send-to-question-id", "q1")
	if err != nil {
		t.Fatal(err)
	}
	inner := got.body["email_notification"].(map[string]interface{})
	if inner["reply_to"] != "oliver.smith@example.com" || inner["send_to_question_id"] != "q1" {
		t.Errorf("body = %v", got.body)
	}
}

func TestEmailNotificationsShowRejectsUnknownKind(t *testing.T) {
	got, _, err := runCLI(t, `{}`, "forms", "email-notifications", "show", "f1", "team")
	if err == nil || err.Error() != `Invalid kind "team". Valid values: admin, respondent.` {
		t.Fatalf("error = %v", err)
	}
	if got.path != "" {
		t.Errorf("request was made to %q", got.path)
	}
}

func TestEmailNotificationsUpdateRejectsUnknownKind(t *testing.T) {
	got, _, err := runCLI(t, `{}`, "forms", "email-notifications", "update", "f1", "team", "--enabled")
	if err == nil {
		t.Fatal("expected error")
	}
	if got.path != "" {
		t.Errorf("request was made to %q", got.path)
	}
}
