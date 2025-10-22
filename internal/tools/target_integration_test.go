package tools

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/kpb/wingit-mcp/internal/ebird"
)

func Test_integration_target_checklist_with_fixtures(t *testing.T) {
	t.Parallel()

	// Read fixtures from package-local testdata/
	pcPath := filepath.Join("testdata", "personal_checklist_example.json")
	rnPath := filepath.Join("testdata", "recent_nearby_example.json")

	pc, err := ebird.LoadPersonalChecklist(pcPath)
	if err != nil {
		t.Fatalf("load personal checklist: %v", err)
	}
	recent, err := ebird.LoadRecentNearby(rnPath)
	if err != nil {
		t.Fatalf("load recent nearby: %v", err)
	}

	seen := ebird.BuildPersonalSeenSet(pc)

	// Adapt []types.RecentObservation -> []RecentObs (engine type for BuildTargetChecklist)
	recs := make([]RecentObs, 0, len(recent))
	for _, r := range recent {
		recs = append(recs, RecentObs{
			SpeciesCode: r.SpeciesCode,
			CommonName:  r.CommonName,
			SciName:     r.SciName,
			LocName:     r.LocName,
			LocID:       r.LocID,
			ObsDt:       r.ObsDt,
			HeardOnly:   r.HeardOnly,
		})
	}

	out, err := BuildTargetChecklist(context.Background(), targetArgs{
		Location:         "35.6870,-105.9378",
		RadiusKm:         20,
		DaysBack:         7,
		IncludeHeardOnly: false, // heard-only should be excluded
		MinFrequency:     0.0,   // keep everything in the sample
		MaxSpecies:       10,
	}, seen, recs)
	if err != nil {
		t.Fatalf("BuildTargetChecklist error: %v", err)
	}

	// Expect:
	// - "clanut" excluded (already seen in personal checklist)
	// - "caltow" excluded (heard-only, IncludeHeardOnly=false)
	// - Remaining lifers: "lewo" then "pinsis" (stable order)
	if len(out.Targets) < 2 {
		t.Fatalf("need at least two targets, got %d", len(out.Targets))
	}
	got := []string{out.Targets[0].SpeciesCode, out.Targets[1].SpeciesCode}
	want := []string{"lewo", "pinsis"}
	if got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("targets = %v, want %v", got, want)
	}

	// Sanity on filter echo
	if out.Filters.Location == "" || out.Filters.RadiusKm == 0 || out.Filters.DaysBack == 0 {
		t.Fatalf("filters not echoed as expected: %+v", out.Filters)
	}
}
