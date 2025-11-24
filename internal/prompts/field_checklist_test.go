package prompts

import "testing"

func Test_BuildFieldChecklistPrompt_includes_location_and_range(t *testing.T) {
	got := BuildFieldChecklistPrompt("Santa Fe, NM", "last 7 days")

	// Must include the key interpolations:
	if !containsAll(got, "Santa Fe, NM", "last 7 days") {
		t.Fatalf("prompt missing location/dayRange:\n%s", got)
	}

	// Must include core instructions (a few durable substrings).
	wantSnippets := []string{
		`WingIt-MCP tool "target_checklist"`,
		`printable field checklist`,
		`"targets" and "filters"`,
		`Group species by approximate recent frequency`,
		`high / medium / low`,
		`Do not reprint the raw JSON; summarize it.`,
	}
	if !containsAll(got, wantSnippets...) {
		t.Fatalf("prompt missing required guidance:\n%s", got)
	}
}

func Test_BuildFieldChecklistPrompt_defaults_are_friendly(t *testing.T) {
	got := BuildFieldChecklistPrompt("", "")
	if !containsAll(got, "this area", "the recent period") {
		t.Fatalf("prompt should include defaults when args empty:\n%s", got)
	}
}

// tiny helper: all substrings must appear in s
func containsAll(s string, subs ...string) bool {
	for _, sub := range subs {
		if !contains(s, sub) {
			return false
		}
	}
	return true
}

// local contains to avoid pulling strings pkg; Go includes it anyway,
// but keeping this explicit makes the test fully self-contained.
func contains(s, sub string) bool {
	return len(sub) == 0 || (len(s) >= len(sub) && indexOf(s, sub) >= 0)
}

// indexOf is a minimal substring search; in real code we'd use strings. Contains,
// but this keeps the test resilient if imports shift.
func indexOf(s, sub string) int {
outer:
	for i := 0; i+len(sub) <= len(s); i++ {
		for j := 0; j < len(sub); j++ {
			if s[i+j] != sub[j] {
				continue outer
			}
		}
		return i
	}
	return -1
}
