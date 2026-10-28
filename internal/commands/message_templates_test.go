package commands

import "testing"

func TestMessageTemplatesListSendsTypeAndStatus(t *testing.T) {
	got, _, err := runCLI(t, `{"message_templates":[],"pagination":{}}`, "forms", "message-templates", "list", "f1", "--type", "sms", "--status", "active")
	if err != nil {
		t.Fatal(err)
	}
	if got.path != "/api/external/v2/forms/f1/message_templates" || got.query.Get("template_type") != "sms" || got.query.Get("status") != "active" {
		t.Errorf("request = %+v", got)
	}
}

func TestMessageTemplatesCreateSendsTemplate(t *testing.T) {
	got, _, err := runCLI(t, `{"message_template":{}}`, "forms", "message-templates", "create", "f1",
		"--name", "Welcome", "--type", "email", "--subject", "Hello", "--body", "<p>Hi Oliver Smith</p>")
	if err != nil {
		t.Fatal(err)
	}
	inner := got.body["message_template"].(map[string]interface{})
	if got.method != "POST" || inner["name"] != "Welcome" || inner["template_type"] != "email" ||
		inner["subject"] != "Hello" || inner["body"] != "<p>Hi Oliver Smith</p>" {
		t.Errorf("request = %+v", got)
	}
}

func TestMessageTemplatesCreateRejectsUnknownType(t *testing.T) {
	got, _, err := runCLI(t, `{}`, "forms", "message-templates", "create", "f1", "--name", "A", "--type", "whatsapp", "--body", "x")
	if err == nil || err.Error() != `Invalid --type "whatsapp". Valid values: email, sms.` {
		t.Fatalf("error = %v", err)
	}
	if got.path != "" {
		t.Errorf("request was made to %q", got.path)
	}
}

func TestMessageTemplatesUpdateRejectsUnknownStatus(t *testing.T) {
	got, _, err := runCLI(t, `{}`, "forms", "message-templates", "update", "f1", "t1", "--status", "paused")
	if err == nil || err.Error() != `Invalid --status "paused". Valid values: active, inactive.` {
		t.Fatalf("error = %v", err)
	}
	if got.path != "" {
		t.Errorf("request was made to %q", got.path)
	}
}

func TestMessageTemplatesUpdateSendsOnlyPassedFields(t *testing.T) {
	got, _, err := runCLI(t, `{"message_template":{}}`, "forms", "message-templates", "update", "f1", "t1", "--status", "inactive")
	if err != nil {
		t.Fatal(err)
	}
	inner := got.body["message_template"].(map[string]interface{})
	if got.method != "PATCH" || got.path != "/api/external/v2/forms/f1/message_templates/t1" || len(inner) != 1 || inner["status"] != "inactive" {
		t.Errorf("request = %+v", got)
	}
}

func TestMessageTemplatesDeleteSendsDeleteToTemplate(t *testing.T) {
	got, _, err := runCLI(t, `{"id":"t1","deleted":true}`, "forms", "message-templates", "delete", "f1", "t1")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "DELETE" || got.path != "/api/external/v2/forms/f1/message_templates/t1" {
		t.Errorf("request = %+v", got)
	}
}
