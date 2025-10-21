package tools

import (
	"context"
	"reflect"
	"testing"
	"time"
)

// ---- Test-time minimal types (match your production types) ----

type targetArgs struct {
	Location         string
	RadiusKm         float64
	DaysBack         int
	IncludeHeardOnly bool
	MinFrequency     float64
	MaxSpecies       int
}

type TargetRow struct {
	SpeciesCode     string
	CommonName      string
	SciName         string
	RecentFrequency float64
	LastSeenNearby  string
}

type targetResult struct {
	Targets []TargetRow
	Filters struct {
		Location         string
		RadiusKm         float64
		DaysBack         int
		IncludeHeardOnly bool
		MinFrequency     float64
		MaxSpecies       int
	}
	ExcludedBecauseAlreadySeen int
}

// RecentObs is the minimal input we need from an eBird “recent nearby” query.
type RecentObs struct {
	SpeciesCode string
	CommonName  string
	SciName     string
	LocName     string
	LocID       string
	ObsDt       string
	HeardOnly   bool
}

// BuildTargetChecklist is the pure “engine” the tool handler will call.
// Implement this in production code; tests here define the contract.
func BuildTargetChecklist(ctx context.Context, args targetArgs, personalSeen map[string]struct{}, recent []RecentObs) (targetResult, error) {
	// NOTE: this is only a placeholder so the test compiles if you paste it verbatim.
	// Replace with real logic. The tests below define expected behavior.
	var out targetResult
	out.Filters.Location = args.Location
	out.Filters.RadiusKm = args.RadiusKm
	out.Filters.DaysBack = args.DaysBack
	out.Filters.IncludeHeardOnly = args.IncludeHeardOnly
	out.Filters.MinFrequency = args.MinFrequency
	out.Filters.MaxSpecies = args.MaxSpecies

	// Extremely naive: keep items that aren't seen and (optionally) not heard-only.
	for _, r := range recent {
		if _, seen := personalSeen[r.SpeciesCode]; seen {
			out.ExcludedBecauseAlreadySeen++
			continue
		}
		if !args.IncludeHeardOnly && r.HeardOnly {
			continue
		}
		out.Targets = append(out.Targets, TargetRow{
			SpeciesCode:     r.SpeciesCode,
			CommonName:      r.CommonName,
			SciName:         r.SciName,
			RecentFrequency: 0.20, // stubbed; real impl will compute
			LastSeenNearby:  r.ObsDt,
		})
		if args.MaxSpecies > 0 && len(out.Targets) >= args.MaxSpecies {
			break
		}
	}
	return out, nil
}

// ---- Tests ----

func Test_build_target_checklist_filters_seen_and_heard_only(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC().Format("2006-01-02")

	args := targetArgs{
		Location:         "35.6870,-105.9378",
		RadiusKm:         20,
		DaysBack:         7,
		IncludeHeardOnly: false, // heard-only should be filtered OUT
		MinFrequency:     0.05,
		MaxSpecies:       10,
	}

	personalSeen := map[string]struct{}{
		"clanut": {}, // already seen Clark's Nutcracker
	}

	recent := []RecentObs{
		{SpeciesCode: "clanut", CommonName: "Clark's Nutcracker", SciName: "Nucifraga columbiana", ObsDt: now, HeardOnly: false},
		{SpeciesCode: "lewo", CommonName: "Lewis's Woodpecker", SciName: "Melanerpes lewis", ObsDt: now, HeardOnly: false},
		{SpeciesCode: "caltow", CommonName: "Canyon Towhee", SciName: "Melozone fusca", ObsDt: now, HeardOnly: true}, // should be dropped
	}

	got, err := BuildTargetChecklist(context.Background(), args, personalSeen, recent)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Expect only the lifer (lewo) to remain.
	wantTargets := []TargetRow{
		{SpeciesCode: "lewo", CommonName: "Lewis's Woodpecker", SciName: "Melanerpes lewis", RecentFrequency: 0.20, LastSeenNearby: now},
	}
	if !reflect.DeepEqual(got.Targets, wantTargets) {
		t.Fatalf("targets mismatch\n got: %#v\nwant: %#v", got.Targets, wantTargets)
	}

	// One species excluded because it was already seen.
	if got.ExcludedBecauseAlreadySeen != 1 {
		t.Fatalf("excludedBecauseAlreadySeen = %d, want 1", got.ExcludedBecauseAlreadySeen)
	}

	// Filters should be echoed back verbatim.
	if got.Filters.Location != args.Location ||
		got.Filters.RadiusKm != args.RadiusKm ||
		got.Filters.DaysBack != args.DaysBack ||
		got.Filters.IncludeHeardOnly != args.IncludeHeardOnly ||
		got.Filters.MinFrequency != args.MinFrequency ||
		got.Filters.MaxSpecies != args.MaxSpecies {
		t.Fatalf("filters echo mismatch: %#v", got.Filters)
	}
}

func Test_build_target_checklist_respects_max_species_cap(t *testing.T) {
	t.Parallel()

	now := "2025-10-06"
	args := targetArgs{
		Location:     "Santa Fe, NM",
		RadiusKm:     25,
		DaysBack:     3,
		MinFrequency: 0.0,
		MaxSpecies:   2, // cap at 2
	}

	personalSeen := map[string]struct{}{} // no lifers seen yet
	recent := []RecentObs{
		{SpeciesCode: "lewo", CommonName: "Lewis's Woodpecker", SciName: "Melanerpes lewis", ObsDt: now},
		{SpeciesCode: "clanut", CommonName: "Clark's Nutcracker", SciName: "Nucifraga columbiana", ObsDt: now},
		{SpeciesCode: "pinsis", CommonName: "Pine Siskin", SciName: "Spinus pinus", ObsDt: now},
	}

	got, err := BuildTargetChecklist(context.Background(), args, personalSeen, recent)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Targets) != 2 {
		t.Fatalf("len(targets) = %d, want 2", len(got.Targets))
	}
	// Order should preserve input order (stable) unless you add ranking later.
	wantOrder := []string{"lewo", "clanut"}
	if got.Targets[0].SpeciesCode != wantOrder[0] || got.Targets[1].SpeciesCode != wantOrder[1] {
		t.Fatalf("order mismatch, got %v,%v want %v,%v",
			got.Targets[0].SpeciesCode, got.Targets[1].SpeciesCode, wantOrder[0], wantOrder[1])
	}
}
