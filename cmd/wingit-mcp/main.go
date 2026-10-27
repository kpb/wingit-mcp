package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/kpb/wingit-mcp/internal/prompts"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/kpb/wingit-mcp/internal/ebird"
	mcpi "github.com/kpb/wingit-mcp/internal/mcp"
	"github.com/kpb/wingit-mcp/internal/tools"
	it "github.com/kpb/wingit-mcp/internal/types"
)

type recentObservationLoader func(context.Context, tools.TargetArgs) ([]it.RecentObservation, error)

type liveClientProvider struct {
	once    sync.Once
	factory func() (*ebird.Client, error)
	client  *ebird.Client
	err     error
}

func newLiveClientProvider(factory func() (*ebird.Client, error)) *liveClientProvider {
	return &liveClientProvider{factory: factory}
}

func (p *liveClientProvider) get() (*ebird.Client, error) {
	p.once.Do(func() {
		p.client, p.err = p.factory()
	})
	return p.client, p.err
}

type recentObservationSource struct {
	offlinePath string
	liveClient  *liveClientProvider
}

func (s recentObservationSource) load(ctx context.Context, args tools.TargetArgs) ([]it.RecentObservation, error) {
	if s.offlinePath != "" {
		return ebird.LoadRecentNearby(s.offlinePath)
	}

	client, err := s.liveClient.get()
	if err != nil {
		return nil, err
	}
	lat, lon, err := parseLatLon(args.Location)
	if err != nil {
		return nil, err
	}
	return client.RecentNearby(ctx, lat, lon, args.RadiusKm, args.DaysBack, 0)
}

func buildTargetChecklist(ctx context.Context, args tools.TargetArgs, seen map[string]struct{}, loadRecent recentObservationLoader) (tools.TargetResult, error) {
	args, err := tools.NormalizeTargetArgs(args)
	if err != nil {
		return tools.TargetResult{}, err
	}

	recent, err := loadRecent(ctx, args)
	if err != nil {
		return tools.TargetResult{}, err
	}

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

	return tools.BuildTargetChecklist(ctx, args, seen, engineRecent)
}

func parseLatLon(input string) (float64, float64, error) {
	parts := strings.Split(strings.TrimSpace(input), ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("location must be \"lat,lon\" (example: \"42.47,-76.45\")")
	}
	lat, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("latitude must be a number between -90 and 90 (example: \"42.47\"): %w", err)
	}
	if math.IsNaN(lat) || math.IsInf(lat, 0) {
		return 0, 0, fmt.Errorf("latitude must be a finite number between -90 and 90")
	}
	lon, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("longitude must be a number between -180 and 180 (example: \"-76.45\"): %w", err)
	}
	if math.IsNaN(lon) || math.IsInf(lon, 0) {
		return 0, 0, fmt.Errorf("longitude must be a finite number between -180 and 180")
	}
	if lat < -90 || lat > 90 {
		return 0, 0, fmt.Errorf("latitude %v out of range; must be between -90 and 90", lat)
	}
	if lon < -180 || lon > 180 {
		return 0, 0, fmt.Errorf("longitude %v out of range; must be between -180 and 180", lon)
	}
	return lat, lon, nil
}

func main() {
	// IMPORTANT: stdio servers must not write to stdout; use stderr for logs
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
	liveClient := newLiveClientProvider(func() (*ebird.Client, error) {
		return ebird.NewClientFromEnv()
	})

	s := mcp.NewServer(&mcp.Implementation{
		Name:    "wingit-mcp",
		Version: "0.1.0",
	}, nil)

	// Register prompts before tools so the host sees them on initialize.
	prompts.Register(s)
	mcpi.RegisterResources(s, pc)

	// Register the target_checklist tool.
	// The SDK infers JSON Schema for input/output from the types you use.
	mcp.AddTool(s, &mcp.Tool{
		Name:        "target_checklist",
		Description: "Return likely new lifers near a location by comparing recent eBird observations with your personal history.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args tools.TargetArgs) (*mcp.CallToolResult, any, error) {
		recentPath := os.Getenv("WINGIT_RECENT_JSON")
		out, err := buildTargetChecklist(ctx, args, seen, func(ctx context.Context, normalized tools.TargetArgs) ([]it.RecentObservation, error) {
			source := recentObservationSource{offlinePath: recentPath, liveClient: liveClient}
			rows, err := source.load(ctx, normalized)
			if err != nil {
				if recentPath != "" {
					// Fail softly: log and continue with empty recent data.
					logger.Printf("WARN: LoadRecentNearby(%q): %v (continuing with empty recent)", recentPath, err)
					return nil, nil
				}
				return nil, err
			}
			return rows, nil
		})
		if err != nil {
			return nil, nil, err
		}

		// Human-friendly text summary for host UI.
		summary := "WingIt-MCP: no candidate lifers"
		if n := len(out.Targets); n > 0 {
			top := out.Targets[0].CommonName
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

	// Run the server on stdio transport.
	if err := s.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		logger.Printf("server failed: %v", err)
	}
}
