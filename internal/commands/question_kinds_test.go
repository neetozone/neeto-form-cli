package commands

import "testing"

func TestQuestionKindsListUsesUnderscoredPath(t *testing.T) {
	got, _, err := runCLI(t, `{"kinds":[]}`, "question-kinds", "list")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "GET" || got.path != "/api/external/v2/question_kinds" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
}
