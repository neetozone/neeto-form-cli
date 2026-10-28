package commands

import "testing"

func TestAPIKeysListSendsSearch(t *testing.T) {
	got, _, err := runCLI(t, `{"api_keys":[],"pagination":{}}`, "api-keys", "list", "--search", "CRM")
	if err != nil {
		t.Fatal(err)
	}
	if got.path != "/api/external/v2/api_keys" || got.query.Get("search") != "CRM" {
		t.Errorf("request = %s %v", got.path, got.query)
	}
}

func TestAPIKeysDeleteRevokesKey(t *testing.T) {
	got, _, err := runCLI(t, `{"notice":"The API key is successfully deleted."}`, "api-keys", "delete", "k1")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "DELETE" || got.path != "/api/external/v2/api_keys/k1" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
}
