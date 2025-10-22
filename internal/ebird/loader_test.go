package ebird

import (
	"path/filepath"
	"testing"
)

func Test_loaders_and_seen_set(t *testing.T) {
	pc, err := LoadPersonalChecklist(filepath.Join("testdata", "personal_checklist_example.json"))
	if err != nil {
		t.Fatalf("LoadPersonalChecklist: %v", err)
	}
	recent, err := LoadRecentNearby(filepath.Join("testdata", "recent_nearby_example.json"))
	if err != nil {
		t.Fatalf("LoadRecentNearby: %v", err)
	}
	seen := BuildPersonalSeenSet(pc)
	if len(seen) == 0 || len(recent) == 0 {
		t.Fatalf("unexpected seen=%d recent=%d", len(seen), len(recent))
	}
}
