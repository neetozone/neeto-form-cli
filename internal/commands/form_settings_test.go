package commands

import (
	"strings"
	"testing"
)

func TestUniqueSubmissionUpdateSendsMethodAndQuestion(t *testing.T) {
	got, _, err := runCLI(t, `{"unique_submission":{}}`, "forms", "unique-submission", "update", "f1", "--method", "record_track", "--question-id", "q1")
	if err != nil {
		t.Fatal(err)
	}
	inner := got.body["unique_submission"].(map[string]interface{})
	if got.method != "PATCH" || got.path != "/api/external/v2/forms/f1/unique_submission" ||
		inner["submission_tracking_method"] != "record_track" || inner["unique_by_record_id"] != "q1" {
		t.Errorf("request = %+v", got)
	}
}

func TestUniqueSubmissionUpdateRejectsUnknownMethod(t *testing.T) {
	got, _, err := runCLI(t, `{}`, "forms", "unique-submission", "update", "f1", "--method", "magic")
	want := `Invalid --method "magic". Valid values: no_track, cookie_track, ip_track, record_track.`
	if err == nil || err.Error() != want {
		t.Fatalf("error = %v", err)
	}
	if got.path != "" {
		t.Errorf("request was made to %q", got.path)
	}
}

func TestAvailabilityShowReadsSettingPath(t *testing.T) {
	got, _, err := runCLI(t, `{"availability":{}}`, "forms", "availability", "show", "f1")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "GET" || got.path != "/api/external/v2/forms/f1/availability" {
		t.Errorf("request = %+v", got)
	}
}

func TestAvailabilityUpdateSendsLimitAndMessage(t *testing.T) {
	got, _, err := runCLI(t, `{"availability":{}}`, "forms", "availability", "update", "f1",
		"--method", "based_on_limit_and_time", "--deadline", "2030-01-01T00:00:00Z", "--limit", "50", "--closed-message", "Closed.")
	if err != nil {
		t.Fatal(err)
	}
	inner := got.body["availability"].(map[string]interface{})
	if inner["manage_access_method"] != "based_on_limit_and_time" || inner["submission_deadline"] != "2030-01-01T00:00:00Z" ||
		inner["submission_limit"] != float64(50) || inner["custom_closed_message"] != "Closed." {
		t.Errorf("body = %v", got.body)
	}
}

func TestAvailabilityUpdateRejectsUnknownMethod(t *testing.T) {
	got, _, err := runCLI(t, `{}`, "forms", "availability", "update", "f1", "--method", "open")
	if err == nil || err.Error() != `Invalid --method "open". Valid values: enable, disable, based_on_submission_count, based_on_time, based_on_limit_and_time.` {
		t.Fatalf("error = %v", err)
	}
	if got.path != "" {
		t.Errorf("request was made to %q", got.path)
	}
}

func TestAccessControlUpdateSendsRepeatedKeywords(t *testing.T) {
	got, _, err := runCLI(t, `{"access_control":{}}`, "forms", "access-control", "update", "f1",
		"--auth", "email", "--allowed-keywords", "example.com,sam.smith@example.org")
	if err != nil {
		t.Fatal(err)
	}
	inner := got.body["access_control"].(map[string]interface{})
	keywords := inner["allowed_keywords"].([]interface{})
	if got.path != "/api/external/v2/forms/f1/access_control" || inner["auth"] != "email" ||
		len(keywords) != 2 || keywords[0] != "example.com" || keywords[1] != "sam.smith@example.org" {
		t.Errorf("request = %+v", got)
	}
}

func TestAccessControlUpdateRejectsUnknownAuth(t *testing.T) {
	got, _, err := runCLI(t, `{}`, "forms", "access-control", "update", "f1", "--auth", "open")
	if err == nil || err.Error() != `Invalid --auth "open". Valid values: no, basic, email, session, invited.` {
		t.Fatalf("error = %v", err)
	}
	if got.path != "" {
		t.Errorf("request was made to %q", got.path)
	}
}

func TestTrackingCodesUpdateSendsEmptyValueToRemove(t *testing.T) {
	got, _, err := runCLI(t, `{"tracking_codes":{}}`, "forms", "tracking-codes", "update", "f1",
		"--google-analytics", "G-ABC123XYZ", "--google-tag-manager", "", "--meta-pixel", "123456789012345")
	if err != nil {
		t.Fatal(err)
	}
	inner := got.body["tracking_codes"].(map[string]interface{})
	if got.path != "/api/external/v2/forms/f1/tracking_codes" || inner["ga_code"] != "G-ABC123XYZ" ||
		inner["gtm_code"] != "" || inner["pixel_id"] != "123456789012345" {
		t.Errorf("request = %+v", got)
	}
}

func TestNeetoChatUpdateSendsFalse(t *testing.T) {
	got, _, err := runCLI(t, `{"neeto_chat":{}}`, "forms", "neeto-chat", "update", "f1", "--enabled=false")
	if err != nil {
		t.Fatal(err)
	}
	inner := got.body["neeto_chat"].(map[string]interface{})
	if got.path != "/api/external/v2/forms/f1/neeto_chat" || inner["enabled"] != false {
		t.Errorf("request = %+v", got)
	}
}

func TestFormDiscountCodesListSendsPagination(t *testing.T) {
	got, _, err := runCLI(t, `{"discount_codes":[],"pagination":{}}`, "forms", "discount-codes", "list", "f1", "--page", "2")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "GET" || got.path != "/api/external/v2/forms/f1/discount_codes" || got.query.Get("page") != "2" {
		t.Errorf("request = %+v", got)
	}
}

func TestFormDiscountCodesUpdateSendsIds(t *testing.T) {
	got, _, err := runCLI(t, `{"enabled_discount_code_ids":[]}`, "forms", "discount-codes", "update", "f1", "--ids", "d1,d2")
	if err != nil {
		t.Fatal(err)
	}
	ids := got.body["discount_codes"].(map[string]interface{})["ids"].([]interface{})
	if got.method != "PATCH" || len(ids) != 2 || ids[0] != "d1" || ids[1] != "d2" {
		t.Errorf("request = %+v", got)
	}
}

func TestFormDiscountCodesUpdateWithoutIdsIsRefused(t *testing.T) {
	_, _, err := runCLI(t, `{"enabled_discount_code_ids":[]}`, "forms", "discount-codes", "update", "f1")
	if err == nil || !strings.Contains(err.Error(), `required flag(s) "ids" not set`) {
		t.Fatalf("err = %v", err)
	}
}

func TestFormDiscountCodesUpdateWithEmptyIdsClearsAll(t *testing.T) {
	got, _, err := runCLI(t, `{"enabled_discount_code_ids":[]}`, "forms", "discount-codes", "update", "f1", "--ids=")
	if err != nil {
		t.Fatal(err)
	}
	ids := got.body["discount_codes"].(map[string]interface{})["ids"].([]interface{})
	if len(ids) != 0 {
		t.Errorf("body = %v", got.body)
	}
}
