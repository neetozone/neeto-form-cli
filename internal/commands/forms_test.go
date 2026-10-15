package commands

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/neetozone/neeto-cli-commons/cli"
	"github.com/neetozone/neeto-cli-commons/config"
	product "github.com/neetozone/neeto-form-cli"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

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

func TestValidateOrder_Valid(t *testing.T) {
	for _, order := range []string{"asc", "desc"} {
		if err := validateOrder(order); err != nil {
			t.Errorf("validateOrder(%q) error = %v, want nil", order, err)
		}
	}
}

func TestValidateOrder_Invalid(t *testing.T) {
	err := validateOrder("newest")
	if err == nil {
		t.Fatal("validateOrder(\"newest\") expected error, got nil")
	}

	want := `Invalid --order "newest". Valid values: asc, desc.`
	if err.Error() != want {
		t.Errorf("validateOrder(\"newest\") error = %q, want %q", err.Error(), want)
	}
}

type recordedRequest struct {
	method string
	path   string
	query  url.Values
	body   map[string]interface{}
}

var (
	testApp     *cli.App
	testAppOnce sync.Once
)

func runCLI(t *testing.T, response string, args ...string) (recordedRequest, string, error) {
	t.Helper()

	var got recordedRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = recordedRequest{method: r.Method, path: r.URL.Path, query: r.URL.Query()}
		if r.Body != nil {
			_ = json.NewDecoder(r.Body).Decode(&got.body)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(response))
	}))
	t.Cleanup(server.Close)

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("NEETOFORM_BASE_URL", server.URL)
	t.Setenv("NEETO_SKIP_VERSION_CHECK", "1")

	configDir := filepath.Join(home, ".config", "neetoform")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}
	creds := `{"credentials":[{"subdomain":"acme","email":"oliver.smith@example.com","session_token":"token"}]}`
	if err := os.WriteFile(filepath.Join(configDir, "auth.json"), []byte(creds), 0o600); err != nil {
		t.Fatal(err)
	}

	testAppOnce.Do(func() {
		cfg, err := config.Parse(product.ConfigYAML)
		if err != nil {
			t.Fatal(err)
		}
		testApp = cli.New(*cfg)
		Register(testApp)
	})

	out := &bytes.Buffer{}
	testApp.Printer.Out = out
	root := testApp.Root()
	root.SetArgs(append(args, "--json"))

	_, err := root.ExecuteC()
	resetFlags(root)
	return got, out.String(), err
}

func resetFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		_ = f.Value.Set(f.DefValue)
		f.Changed = false
	})
	for _, sub := range cmd.Commands() {
		resetFlags(sub)
	}
}

func TestFormsListSendsSearchAndOrder(t *testing.T) {
	got, _, err := runCLI(t, `{"forms":[],"pagination":{}}`, "forms", "list", "--search", "event sign up", "--order", "desc")
	if err != nil {
		t.Fatal(err)
	}
	if got.path != "/api/external/v2/forms" {
		t.Errorf("path = %q", got.path)
	}
	if got.query.Get("search") != "event sign up" || got.query.Get("order") != "desc" {
		t.Errorf("query = %v", got.query)
	}
}

func TestFormsListOmitsEmptySearchAndOrder(t *testing.T) {
	got, _, err := runCLI(t, `{"forms":[],"pagination":{}}`, "forms", "list")
	if err != nil {
		t.Fatal(err)
	}
	if got.query.Has("search") || got.query.Has("order") {
		t.Errorf("query = %v, want no search or order", got.query)
	}
}

func TestFormsListRejectsInvalidOrder(t *testing.T) {
	got, _, err := runCLI(t, `{}`, "forms", "list", "--order", "newest")
	if err == nil || err.Error() != `Invalid --order "newest". Valid values: asc, desc.` {
		t.Fatalf("error = %v", err)
	}
	if got.path != "" {
		t.Errorf("request was made to %q", got.path)
	}
}

func TestFormsShowFetchesFormWithQuestions(t *testing.T) {
	got, out, err := runCLI(t, `{"id":"abc","title":"Contact sales","questions":[{"slug":"work-email"}]}`, "forms", "show", "abc")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != http.MethodGet || got.path != "/api/external/v2/forms/abc" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
	if !strings.Contains(out, "work-email") {
		t.Errorf("output missing questions: %s", out)
	}
}

func TestFormsShowRequiresFormID(t *testing.T) {
	if _, _, err := runCLI(t, `{}`, "forms", "show"); err == nil {
		t.Fatal("expected error for missing form id")
	}
}

func TestFormsListAcceptsUppercaseOrder(t *testing.T) {
	got, _, err := runCLI(t, `{"forms":[],"pagination":{}}`, "forms", "list", "--order", "DESC")
	if err != nil {
		t.Fatal(err)
	}
	if got.query.Get("order") != "desc" {
		t.Errorf("query = %v", got.query)
	}
}
