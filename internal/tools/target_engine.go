package tools

import (
	"context"
	"fmt"
	"math"
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
// Ranking: by RecentFrequency (desc), then by recency (ObsDt desc), then stable.
func BuildTargetChecklist(_ context.Context, args targetArgs, personalSeen map[string]struct{}, recent []RecentObs) (targetResult, error) {
	var out targetResult

	var err error
	args, err = NormalizeTargetArgs(args)
	if err != nil {
		return out, err
	}

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

// NormalizeTargetArgs validates and normalizes the shared query used for both
// observation loading and result filters. Non-positive radius/days use defaults;
// values above the eBird API limits are clamped to those limits.
func NormalizeTargetArgs(a TargetArgs) (TargetArgs, error) {
	if strings.TrimSpace(a.Location) == "" {
		return a, fmt.Errorf("location is required")
	}
	if math.IsNaN(a.RadiusKm) || math.IsInf(a.RadiusKm, 0) {
		return a, fmt.Errorf("radiusKm must be a finite number")
	}
	if math.IsNaN(a.MinFrequency) || math.IsInf(a.MinFrequency, 0) {
		return a, fmt.Errorf("minFrequency must be a finite number")
	}

	if a.RadiusKm <= 0 {
		a.RadiusKm = defaultRadiusKm
	} else if a.RadiusKm > 50 {
		a.RadiusKm = 50
	}
	if a.DaysBack <= 0 {
		a.DaysBack = defaultDaysBack
	} else if a.DaysBack > 30 {
		a.DaysBack = 30
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
	return a, nil
}
