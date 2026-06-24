package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestPaginationParams_DefaultsToFirstPage(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	addPaginationFlags(cmd)

	params := paginationParams(cmd)

	if got := params.Get("page_number"); got != "1" {
		t.Errorf("page_number = %q, want 1", got)
	}
	if got := params.Get("page_size"); got != "" {
		t.Errorf("page_size = %q, want empty", got)
	}
}

func TestPaginationParams_RespectsExplicitFlags(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	addPaginationFlags(cmd)
	if err := cmd.Flags().Set("page", "2"); err != nil {
		t.Fatalf("Set(page) error = %v", err)
	}
	if err := cmd.Flags().Set("page-size", "50"); err != nil {
		t.Fatalf("Set(page-size) error = %v", err)
	}

	params := paginationParams(cmd)

	if got := params.Get("page_number"); got != "2" {
		t.Errorf("page_number = %q, want 2", got)
	}
	if got := params.Get("page_size"); got != "50" {
		t.Errorf("page_size = %q, want 50", got)
	}
}

func TestReadJSONFile_Valid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "data.json")
	if err := os.WriteFile(path, []byte(`{"name":"test","count":42}`), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	result, err := readJSONFile(path)
	if err != nil {
		t.Fatalf("readJSONFile() error = %v", err)
	}

	if result["name"] != "test" {
		t.Errorf("name = %v, want test", result["name"])
	}
	if result["count"] != float64(42) {
		t.Errorf("count = %v, want 42", result["count"])
	}
}

func TestReadJSONFile_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte(`{not valid`), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	_, err := readJSONFile(path)
	if err == nil {
		t.Error("readJSONFile() expected error for invalid JSON")
	}
}

func TestReadJSONFile_NotFound(t *testing.T) {
	_, err := readJSONFile("/nonexistent/file.json")
	if err == nil {
		t.Error("readJSONFile() expected error for missing file")
	}
}
