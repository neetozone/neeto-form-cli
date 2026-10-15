package commands

import "testing"

func TestTeamMembersListSendsSortAndOrder(t *testing.T) {
	got, _, err := runCLI(t, `{"team_members":[],"pagination":{}}`, "team-members", "list", "--sort", "email", "--order", "desc")
	if err != nil {
		t.Fatal(err)
	}
	if got.query.Get("sort") != "email" || got.query.Get("order") != "desc" {
		t.Errorf("query = %v", got.query)
	}
}

func TestTeamMembersListSortsByCreationTimeWhenOnlyOrderIsGiven(t *testing.T) {
	got, _, err := runCLI(t, `{"team_members":[],"pagination":{}}`, "team-members", "list", "--order", "desc")
	if err != nil {
		t.Fatal(err)
	}
	if got.query.Get("sort") != "created_at" || got.query.Get("order") != "desc" {
		t.Errorf("query = %v", got.query)
	}
}

func TestTeamMembersCreateSendsInvitedByOnlyWhenSet(t *testing.T) {
	got, _, err := runCLI(t, `{}`, "team-members", "create", "--emails", "sam.smith@example.com", "--role", "Standard", "--invited-by", "oliver.smith@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "POST" || got.body["invited_by"] != "oliver.smith@example.com" {
		t.Errorf("request = %s %v", got.method, got.body)
	}

	got, _, err = runCLI(t, `{}`, "team-members", "create", "--emails", "sam.smith@example.com", "--role", "Standard")
	if err != nil {
		t.Fatal(err)
	}
	if _, present := got.body["invited_by"]; present {
		t.Errorf("invited_by sent without the flag: %v", got.body)
	}
}
