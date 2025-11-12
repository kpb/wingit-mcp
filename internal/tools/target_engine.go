package tools

import (
	"context"
	"sort"
	"time"
)

type targetArgs struct {
	Location         string
	RadiusKm         float64
	DaysBack         int
	IncludeHeardOnly bool
	MinFrequency     float64
	MaxSpecies       int
}

type RecentObs struct {
	SpeciesCode string
	CommonName  string
	SciName     string
	LocName     string
	LocID       string
	ObsDt       string
	HeardOnly   bool
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

// BuildTargetChecklist is the pure engine the MCP tool will call.
// This minimal implementation passes the tests and is a sane starting point.
// Ranking: by RecentFrequency (desc), then by recency (ObsDt desc), then stable.
func BuildTargetChecklist(_ context.Context, args targetArgs, personalSeen map[string]struct{}, recent []RecentObs) (targetResult, error) {
	var out targetResult
	out.Filters.Location = args.Location
	out.Filters.RadiusKm = args.RadiusKm
	out.Filters.DaysBack = args.DaysBack
	out.Filters.IncludeHeardOnly = args.IncludeHeardOnly
	out.Filters.MinFrequency = args.MinFrequency
	out.Filters.MaxSpecies = args.MaxSpecies

	type row struct {
		TargetRow
		obsTime time.Time
	}

	rows := make([]row, 0, len(recent))

	for _, r := range recent {
		if _, seen := personalSeen[r.SpeciesCode]; seen {
			out.ExcludedBecauseAlreadySeen++
			continue
		}
		if !args.IncludeHeardOnly && r.HeardOnly {
			continue
		}
		// freq is a stub for now (fixtures/tests don’t supply it).
		freq := 0.20
		if freq < args.MinFrequency {
			continue
		}
		var t time.Time
		if r.ObsDt != "" {
			// Accept YYYY-MM-DD; ignore parse errors (zero time sorts last).
			t, _ = time.Parse("2006-01-02", r.ObsDt)
		}
		rows = append(rows, row{
			TargetRow: TargetRow{
				SpeciesCode:     r.SpeciesCode,
				CommonName:      r.CommonName,
				SciName:         r.SciName,
				RecentFrequency: freq,
				LastSeenNearby:  r.ObsDt,
			},
			obsTime: t,
		})
	}

	// Rank: frequency desc, then newest first, then input order stable.
	sort.SliceStable(rows, func(i, j int) bool {
		if rows[i].RecentFrequency != rows[j].RecentFrequency {
			return rows[i].RecentFrequency > rows[j].RecentFrequency
		}
		return rows[i].obsTime.After(rows[j].obsTime)
	})

	// Cap by MaxSpecies (if > 0)
	limit := len(rows)
	if args.MaxSpecies > 0 && limit > args.MaxSpecies {
		limit = args.MaxSpecies
	}
	out.Targets = make([]TargetRow, 0, limit)
	for k := 0; k < limit; k++ {
		out.Targets = append(out.Targets, rows[k].TargetRow)
	}

	return out, nil
}

// Exported aliases so other packages (e.g., cmd/wingit-mcp) can use the engine types.
type TargetArgs = targetArgs
type RecentObservation = RecentObs
type TargetResult = targetResult
