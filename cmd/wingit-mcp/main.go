package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/kpb/wingit-mcp/internal/ebird"
	"github.com/kpb/wingit-mcp/internal/tools"
	it "github.com/kpb/wingit-mcp/internal/types"
)

type targetArgs struct {
	Location         string  `json:"location" jsonschema:"Place name, hotspot code (L123...), or 'lat,lon'"`
	RadiusKm         float64 `json:"radiusKm" jsonschema:"Search radius in kilometers" default:"20"`
	DaysBack         int     `json:"daysBack" jsonschema:"How far back to look for recent obs" default:"7"`
	IncludeHeardOnly bool    `json:"includeHeardOnly" default:"false"`
	MinFrequency     float64 `json:"minFrequency" jsonschema:"0..1 recent proportion threshold" default:"0.05"`
	MaxSpecies       int     `json:"maxSpecies" jsonschema:"Cap results" default:"40"`
}

type targetResult struct {
	Targets []struct {
		SpeciesCode     string  `json:"speciesCode"`
		CommonName      string  `json:"commonName"`
		SciName         string  `json:"sciName"`
		RecentFrequency float64 `json:"recentFrequency"`
		LastSeenNearby  string  `json:"lastSeenNearby,omitempty"`
	} `json:"targets"`
	Filters struct {
		Location         string  `json:"location"`
		RadiusKm         float64 `json:"radiusKm"`
		DaysBack         int     `json:"daysBack"`
		IncludeHeardOnly bool    `json:"includeHeardOnly"`
		MinFrequency     float64 `json:"minFrequency"`
		MaxSpecies       int     `json:"maxSpecies"`
	} `json:"filters"`
	ExcludedBecauseAlreadySeen int `json:"excludedBecauseAlreadySeen"`
}

func main() {
	// IMPORTANT: stdio servers must not write to stdout; use stderr for logs. :contentReference[oaicite:1]{index=1}
	logger := log.New(os.Stderr, "wingit-mcp: ", log.LstdFlags|log.Lmsgprefix)

	// --- Config: load personal checklist path from env, build seen set ---
	personalPath := os.Getenv("WINGIT_PERSONAL_JSON")
	if personalPath == "" {
		logger.Printf("ERROR: WINGIT_PERSONAL_JSON is not set")
		os.Exit(2)
	}
	pc, err := ebird.LoadPersonalChecklist(personalPath)
	if err != nil {
		logger.Printf("ERROR: LoadPersonalChecklist(%q): %v", personalPath, err)
		os.Exit(2)
	}
	seen := ebird.BuildPersonalSeenSet(pc)
	logger.Printf("loaded personal checklist: species=%d (seen set size)", len(seen))

	s := mcp.NewServer(&mcp.Implementation{
		Name:    "wingit-mcp",
		Version: "0.1.0",
	}, nil)

	// Register the target_checklist tool.
	// The SDK infers JSON Schema for input/output from the types you use. :contentReference[oaicite:2]{index=2}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "target_checklist",
		Description: "Return likely new lifers near a location by comparing recent eBird observations with your personal history.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args tools.TargetArgs) (*mcp.CallToolResult, any, error) {
		// Load "recent nearby" from env JSON (offline demo path for now).
		recentPath := os.Getenv("WINGIT_RECENT_JSON")
		var recent []it.RecentObservation
		if recentPath != "" {
			rows, err := ebird.LoadRecentNearby(recentPath)
			if err != nil {
				// Fail softly: log and continue with empty recent data.
				logger.Printf("WARN: LoadRecentNearby(%q): %v (continuing with empty recent)", recentPath, err)
			} else {
				recent = rows
			}
		} else {
			logger.Printf("INFO: WINGIT_RECENT_JSON not set; continuing with empty recent")
		}

		// Adapt internal/types -> engine's RecentObservation
		engineRecent := make([]tools.RecentObservation, 0, len(recent))
		for _, r := range recent {
			engineRecent = append(engineRecent, tools.RecentObservation{
				SpeciesCode: r.SpeciesCode,
				CommonName:  r.CommonName,
				SciName:     r.SciName,
				LocName:     r.LocName,
				LocID:       r.LocID,
				ObsDt:       r.ObsDt,
				HeardOnly:   r.HeardOnly,
			})
		}

		// Call the pure engine.
		out, err := tools.BuildTargetChecklist(ctx, args, seen, engineRecent)
		if err != nil {
			return nil, nil, err
		}

		// Human-friendly text summary for host UI.
		summary := "WingIt-MCP: no candidate lifers"
		if n := len(out.Targets); n > 0 {
			top := out.Targets[0].CommonName
			summary = os.ExpandEnv(
				// Example: “3 candidate lifers; top: Lewis's Woodpecker”
				// (short and useful in the host’s call result view)
				// not actually using env vars here; ExpandEnv is a no-op—just compact code.
				// Feel free to format however you like.
				// Avoid stdout; host reads JSON-RPC on stdout.
				// Logging already goes to stderr via logger above.
				// Return structured JSON too (second return value).
				// Hosts that understand MCP structured content can render tables.
				// Others will still show the text content.
				// Keep it brief and actionable.
				// n.b.: You can add more details later (e.g., location hints).
				"",
			)
			summary = // short, explicit:
				func(n int, top string) string {
					return fmt.Sprintf("%d candidate lifers; top: %s", n, top)
				}(n, top)
		}

		res := &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: summary},
			},
		}

		// Return both: user-facing text and structured JSON (engine result).
		return res, out, nil
	})

	// TODO: Add prompts/resources next:
	// mcp.AddPrompt(s, ...)
	// mcp.AddResource(s, ...) or use resource templates per SDK capabilities.

	// Run the server on stdio transport. :contentReference[oaicite:3]{index=3}
	if err := s.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		logger.Printf("server failed: %v", err)
	}
}
