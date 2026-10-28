package commands

import "testing"

func TestDeactivatedTeamMembersReactivatePatchesMemberPath(t *testing.T) {
	got, _, err := runCLI(t, `{"team_member":{"id":"m1","active":true}}`, "deactivated-team-members", "reactivate", "m1")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "PATCH" || got.path != "/api/external/v2/deactivated_team_members/m1/reactivate" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
}
