package tools

import (
	"context"
	"fmt"
	"sort"
	"strings"
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

// Exported aliases so other packages (cmd/wingit-mcp) can use engine types.
type TargetArgs = targetArgs
type RecentObservation = RecentObs
type TargetResult = targetResult

const (
	defaultRadiusKm   = 20.0
	defaultDaysBack   = 7
	defaultMaxSpecies = 40
)

// BuildTargetChecklist is the pure engine the MCP tool will call.
// This minimal implementation passes the tests and is a sane starting point.
// Ranking: by RecentFrequency (desc), then by recency (ObsDt desc), then stable.
func BuildTargetChecklist(_ context.Context, args targetArgs, personalSeen map[string]struct{}, recent []RecentObs) (targetResult, error) {
	var out targetResult

	// Hard validation: we require *some* location string.
	if strings.TrimSpace(args.Location) == "" {
		return out, fmt.Errorf("location is required")
	}

	// Soft validation: normalize obviously bad numeric inputs.
	args = normalizeArgs(args)

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

		// For now, frequency is a stubbed constant until we wire real stats.
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

// normalizeArgs clamps obviously bad numeric values to sane defaults.
func normalizeArgs(a targetArgs) targetArgs {
	if a.RadiusKm <= 0 {
		a.RadiusKm = defaultRadiusKm
	}
	if a.DaysBack <= 0 {
		a.DaysBack = defaultDaysBack
	}
	if a.MaxSpecies <= 0 {
		a.MaxSpecies = defaultMaxSpecies
	}
	if a.MinFrequency < 0 {
		a.MinFrequency = 0
	}
	if a.MinFrequency > 1 {
		a.MinFrequency = 1
	}
	return a
}
