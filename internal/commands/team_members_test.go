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

func TestTeamMembersListOmitsSortAndOrderByDefault(t *testing.T) {
	got, _, err := runCLI(t, `{"team_members":[],"pagination":{}}`, "team-members", "list")
	if err != nil {
		t.Fatal(err)
	}
	if got.query.Has("sort") || got.query.Has("order") {
		t.Errorf("query = %v, want no sort or order", got.query)
	}
}

func TestTeamMembersListRejectsInvalidSort(t *testing.T) {
	got, _, err := runCLI(t, `{}`, "team-members", "list", "--sort", "name")
	want := `Invalid --sort "name". Valid values: created_at, updated_at, email, first_name, last_name.`
	if err == nil || err.Error() != want {
		t.Fatalf("error = %v", err)
	}
	if got.path != "" {
		t.Errorf("request was made to %q", got.path)
	}
}

func TestTeamMembersListRejectsInvalidOrder(t *testing.T) {
	got, _, err := runCLI(t, `{}`, "team-members", "list", "--order", "up")
	if err == nil || err.Error() != `Invalid --order "up". Valid values: asc, desc.` {
		t.Fatalf("error = %v", err)
	}
	if got.path != "" {
		t.Errorf("request was made to %q", got.path)
	}
}

func TestTeamMembersCreateSendsInvitedBy(t *testing.T) {
	got, _, err := runCLI(t, `{}`, "team-members", "create", "--emails", "sam.smith@example.com", "--role", "Standard", "--invited-by", "oliver.smith@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "POST" || got.body["invited_by"] != "oliver.smith@example.com" {
		t.Errorf("request = %s %v", got.method, got.body)
	}
}

func TestTeamMembersCreateOmitsInvitedByByDefault(t *testing.T) {
	got, _, err := runCLI(t, `{}`, "team-members", "create", "--emails", "sam.smith@example.com", "--role", "Standard")
	if err != nil {
		t.Fatal(err)
	}
	if _, present := got.body["invited_by"]; present {
		t.Errorf("invited_by sent without the flag: %v", got.body)
	}
}
