package commands

import (
	"reflect"
	"testing"
)

const formJSON = `{"form":{"id":"f1","title":"Customer feedback"}}`

func TestFormsCreateSendsNestedFormBody(t *testing.T) {
	got, _, err := runCLI(t, formJSON, "forms", "create", "--title", "Customer feedback", "--style", "one_question_at_a_time", "--created-by-email", "oliver.smith@example.com")
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]interface{}{"form": map[string]interface{}{
		"title": "Customer feedback", "style": "one_question_at_a_time", "created_by_email": "oliver.smith@example.com",
	}}
	if got.method != "POST" || got.path != "/api/external/v2/forms" || !reflect.DeepEqual(got.body, want) {
		t.Errorf("request = %+v", got)
	}
}

func TestFormsCreateRejectsInvalidStyle(t *testing.T) {
	got, _, err := runCLI(t, formJSON, "forms", "create", "--title", "Customer feedback", "--style", "grid")
	if err == nil || err.Error() != `Invalid --style "grid". Valid values: single_page, one_question_at_a_time.` {
		t.Fatalf("error = %v", err)
	}
	if got.path != "" {
		t.Errorf("request was sent: %+v", got)
	}
}

func TestFormsUpdateRenamesForm(t *testing.T) {
	got, _, err := runCLI(t, formJSON, "forms", "update", "f1", "--title", "Event sign up")
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]interface{}{"form": map[string]interface{}{"title": "Event sign up"}}
	if got.method != "PATCH" || got.path != "/api/external/v2/forms/f1" || !reflect.DeepEqual(got.body, want) {
		t.Errorf("request = %+v", got)
	}
}

func TestFormsDeleteRemovesOneForm(t *testing.T) {
	got, _, err := runCLI(t, `{}`, "forms", "delete", "f1")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "DELETE" || got.path != "/api/external/v2/forms/f1" {
		t.Errorf("request = %+v", got)
	}
}

func TestFormsBulkDeleteSendsIDs(t *testing.T) {
	got, _, err := runCLI(t, `{"form_ids":["f1","f2"]}`, "forms", "bulk-delete", "f1", "f2")
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]interface{}{"ids": []interface{}{"f1", "f2"}}
	if got.method != "DELETE" || got.path != "/api/external/v2/forms/bulk_destroy" || !reflect.DeepEqual(got.body, want) {
		t.Errorf("request = %+v", got)
	}
}

func TestFormsArchiveSendsOperationAndIDs(t *testing.T) {
	got, _, err := runCLI(t, `{"forms":[]}`, "forms", "archive", "f1", "f2")
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]interface{}{"ids": []interface{}{"f1", "f2"}, "operation": "archive"}
	if got.method != "PATCH" || got.path != "/api/external/v2/forms/bulk_update" || !reflect.DeepEqual(got.body, want) {
		t.Errorf("request = %+v", got)
	}
}

func TestFormsUnfavoriteSendsOperation(t *testing.T) {
	got, _, err := runCLI(t, `{"forms":[]}`, "forms", "unfavorite", "f1")
	if err != nil {
		t.Fatal(err)
	}
	if got.body["operation"] != "unfavorite" {
		t.Errorf("body = %v", got.body)
	}
}

func TestFormsCloneCallsCloneEndpoint(t *testing.T) {
	got, _, err := runCLI(t, formJSON, "forms", "clone", "f1")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "POST" || got.path != "/api/external/v2/forms/f1/clone" {
		t.Errorf("request = %+v", got)
	}
}

func TestFormsPublishPatchesPublishEndpoint(t *testing.T) {
	got, _, err := runCLI(t, formJSON, "forms", "publish", "f1")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "PATCH" || got.path != "/api/external/v2/forms/f1/publish" {
		t.Errorf("request = %+v", got)
	}
}

func TestFormsResetPatchesResetEndpoint(t *testing.T) {
	got, _, err := runCLI(t, formJSON, "forms", "reset", "f1")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "PATCH" || got.path != "/api/external/v2/forms/f1/reset" {
		t.Errorf("request = %+v", got)
	}
}

func TestFormsCountsFetchesCounts(t *testing.T) {
	got, _, err := runCLI(t, `{"total":12,"active":9,"archived":3,"favorite":2}`, "forms", "counts")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "GET" || got.path != "/api/external/v2/forms/counts" {
		t.Errorf("request = %+v", got)
	}
}

func TestFormsTransferSendsIDsAndDestination(t *testing.T) {
	got, _, err := runCLI(t, `{"form_ids":["f1","f2"],"destination_workspace":"acme-eu"}`, "forms", "transfer", "f1", "f2", "--to", "acme-eu")
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]interface{}{"ids": []interface{}{"f1", "f2"}, "destination_workspace": "acme-eu"}
	if got.method != "POST" || got.path != "/api/external/v2/forms/transfer" || !reflect.DeepEqual(got.body, want) {
		t.Errorf("request = %+v", got)
	}
}
