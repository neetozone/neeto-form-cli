package commands

import "testing"

func TestValidateFormStatus_Valid(t *testing.T) {
	for _, status := range []string{"active", "archived", "favorite"} {
		if err := validateFormStatus(status); err != nil {
			t.Errorf("validateFormStatus(%q) error = %v, want nil", status, err)
		}
	}
}

func TestValidateFormStatus_Invalid(t *testing.T) {
	err := validateFormStatus("pinned")
	if err == nil {
		t.Fatal("validateFormStatus(\"pinned\") expected error, got nil")
	}

	want := `Invalid --status "pinned". Valid values: active, archived, favorite.`
	if err.Error() != want {
		t.Errorf("validateFormStatus(\"pinned\") error = %q, want %q", err.Error(), want)
	}
}
