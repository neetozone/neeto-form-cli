package commands

import "testing"

func TestSaveWithNothingToSendIsRefusedBeforeAnyRequest(t *testing.T) {
	got, _, err := runCLI(t, `{}`, "forms", "conditions", "create", "form-1")
	if err == nil || err.Error() != "Nothing to send. Pass --json-file or set at least one flag; see --help." {
		t.Fatalf("error = %v", err)
	}
	if got.path != "" {
		t.Errorf("request was made to %q", got.path)
	}
}

func TestBodyWithNoFieldsHasNothingToSend(t *testing.T) {
	if !hasNothingToSend(map[string]any{}) {
		t.Error("an empty body was sent")
	}
}

func TestBodyWithAFilledResourceAndAnEmptyExtraIsSent(t *testing.T) {
	if hasNothingToSend(map[string]any{"theme": map[string]any{"name": "Brand"}, "meta": map[string]any{}}) {
		t.Error("a body with a filled resource was refused")
	}
}
