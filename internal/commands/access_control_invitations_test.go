package commands

import "testing"

func TestAccessControlInvitationsCreateSendsEmails(t *testing.T) {
	got, _, err := runCLI(t, `{"invited_emails":[]}`, "forms", "access-control-invitations", "create", "f1",
		"--emails", "oliver.smith@example.com,sam.smith@example.com")
	if err != nil {
		t.Fatal(err)
	}
	emails := got.body["emails"].([]interface{})
	if got.method != "POST" || got.path != "/api/external/v2/forms/f1/access_control_invitations" ||
		len(emails) != 2 || emails[0] != "oliver.smith@example.com" {
		t.Errorf("request = %+v", got)
	}
}

func TestAccessControlInvitationsDeleteSendsEmailsInBody(t *testing.T) {
	got, _, err := runCLI(t, `{"invitations_count":0}`, "forms", "access-control-invitations", "delete", "f1", "--emails", "oliver.smith@example.com")
	if err != nil {
		t.Fatal(err)
	}
	emails := got.body["emails"].([]interface{})
	if got.method != "DELETE" || got.path != "/api/external/v2/forms/f1/access_control_invitations" ||
		len(emails) != 1 || emails[0] != "oliver.smith@example.com" {
		t.Errorf("request = %+v", got)
	}
}

func TestAccessControlInvitationsListSendsPagination(t *testing.T) {
	got, _, err := runCLI(t, `{"invitations":[],"pagination":{}}`, "forms", "access-control-invitations", "list", "f1", "--page-size", "10")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "GET" || got.query.Get("page_size") != "10" {
		t.Errorf("request = %+v", got)
	}
}
