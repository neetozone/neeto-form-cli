package commands

import "testing"

func TestPageSizeWithoutPageRequestsTheFirstPage(t *testing.T) {
	got, _, err := runCLI(t, `{"forms":[],"pagination":{}}`, "forms", "list", "--page-size", "2")
	if err != nil {
		t.Fatal(err)
	}
	if got.query.Get("page") != "1" || got.query.Get("page_size") != "2" {
		t.Errorf("query = %v, want page 1 and page_size 2", got.query)
	}
}

func TestNoPaginationFlagsSendNoPage(t *testing.T) {
	got, _, err := runCLI(t, `{"forms":[],"pagination":{}}`, "forms", "list")
	if err != nil {
		t.Fatal(err)
	}
	if got.query.Has("page") || got.query.Has("page_size") {
		t.Errorf("query = %v, want no paging", got.query)
	}
}
