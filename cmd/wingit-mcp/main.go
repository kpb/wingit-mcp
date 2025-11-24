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

func registerPrompts(s *mcp.Server) {
	prompt := &mcp.Prompt{
		Name:        "field_checklist",
		Description: "Format WingIt target_checklist results as a printable field checklist.",
		Arguments: []*mcp.PromptArgument{
			{
				Name:        "location",
				Description: "Name of the birding location or general area.",
				Required:    true,
			},
			{
				Name:        "dayRange",
				Description: "How far back the recent observations go (e.g. 'last 7 days').",
				Required:    false,
			},
		},
	}

	promptHandler := func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		args := req.Params.Arguments
		loc := args["location"]
		dayRange := args["dayRange"]

		if loc == "" {
			loc = "this area"
		}
		if dayRange == "" {
			dayRange = "the recent period"
		}

		text := fmt.Sprintf(
			`You are a birding assistant. The user has just called the WingIt-MCP tool "target_checklist" to get likely new lifers near %s for %s.

Using the tool output provided in this conversation (JSON with "targets" and "filters"), produce a concise, printable field checklist:

- Focus only on likely lifers (the "targets" array).
- Group species by approximate recent frequency (high / medium / low) based on "recentFrequency".
- For each species, show: common name, scientific name, and a short note like "seen recently at <locName>" if present.
- Keep it compact, suitable for printing or quick reference in the field.
- Do not reprint the raw JSON; summarize it.

If there are no targets, explain that there are no likely new lifers for this query and suggest broadening radius or daysBack.`, loc, dayRange)

		return &mcp.GetPromptResult{
			Description: "Format WingIt target_checklist results as a printable field checklist.",
			Messages: []*mcp.PromptMessage{
				{
					Role: "user",
					Content: &mcp.TextContent{
						Text: text,
					},
				},
			},
		}, nil
	}

	s.AddPrompt(prompt, promptHandler)
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

	// Register prompts before tools so the host sees them on initialize.
	registerPrompts(s)

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
