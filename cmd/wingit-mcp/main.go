package main

import (
	"context"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
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

	s := mcp.NewServer(&mcp.Implementation{
		Name:    "wingit-mcp",
		Version: "0.1.0",
	}, nil)

	// Register the target_checklist tool.
	// The SDK infers JSON Schema for input/output from the types you use. :contentReference[oaicite:2]{index=2}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "target_checklist",
		Description: "Return likely new lifers near a location by comparing recent eBird observations with your personal history.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args targetArgs) (*mcp.CallToolResult, any, error) {
		// TODO: wire to personal eBird CSV+index and eBird API client.
		// For now, return a tiny, deterministic stub so hosts can demo it.

		out := targetResult{
			Targets: []struct {
				SpeciesCode     string  `json:"speciesCode"`
				CommonName      string  `json:"commonName"`
				SciName         string  `json:"sciName"`
				RecentFrequency float64 `json:"recentFrequency"`
				LastSeenNearby  string  `json:"lastSeenNearby,omitempty"`
			}{
				{"clanut", "Clark's Nutcracker", "Nucifraga columbiana", 0.21, "2025-10-06"},
			},
			ExcludedBecauseAlreadySeen: 182,
		}
		out.Filters.Location = args.Location
		out.Filters.RadiusKm = args.RadiusKm
		out.Filters.DaysBack = args.DaysBack
		out.Filters.IncludeHeardOnly = args.IncludeHeardOnly
		out.Filters.MinFrequency = args.MinFrequency
		out.Filters.MaxSpecies = args.MaxSpecies

		// Text content for quick UX in hosts; JSON “out” for structured consumers.
		res := &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "WingIt-MCP: 1 candidate lifer → Clark's Nutcracker (seen nearby 2025-10-06)"},
			},
		}
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
