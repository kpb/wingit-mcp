package prompts

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Register attaches the "field_checklist" prompt to the MCP server.
func Register(s *mcp.Server) {
	prompt := &mcp.Prompt{
		Name:        "field_checklist",
		Description: "Format WingIt target_checklist results as a printable field checklist.",
		Arguments: []*mcp.PromptArgument{
			{Name: "location", Description: "Name of the birding location or area.", Required: true},
			{Name: "dayRange", Description: "Time window label (e.g., 'last 7 days').", Required: false},
		},
	}

	handler := func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		loc := req.Params.Arguments["location"]
		day := req.Params.Arguments["dayRange"]

		text := BuildFieldChecklistPrompt(loc, day)

		return &mcp.GetPromptResult{
			Description: prompt.Description,
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

	s.AddPrompt(prompt, handler)
}
