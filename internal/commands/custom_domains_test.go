package commands

import "testing"

func TestCustomDomainsShowEscapesHostname(t *testing.T) {
	got, _, err := runCLI(t, `{"custom_domain":{"hostname":"forms.example.com"}}`, "custom-domains", "show", "forms.example.com")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "GET" || got.path != "/api/external/v2/custom_domains/forms.example.com" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
}

func TestCustomDomainsCreateSendsHostname(t *testing.T) {
	got, _, err := runCLI(t, `{"custom_domain":{"hostname":"forms.example.com"}}`, "custom-domains", "create", "forms.example.com")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "POST" || got.path != "/api/external/v2/custom_domains" || got.body["hostname"] != "forms.example.com" {
		t.Errorf("request = %s %s %v", got.method, got.path, got.body)
	}
}
