package mcp

import (
	"context"
	"encoding/json"

	it "github.com/kpb/wingit-mcp/internal/types"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

func RegisterResources(s *sdk.Server, pc *it.PersonalChecklist) {
	const personalURI = "wingit://personal-checklist"

	s.AddResource(&sdk.Resource{
		URI:         personalURI,
		MIMEType:    "application/json",
		Name:        "Personal eBird Checklist (normalized)",
		Description: "Read-only view of the user's normalized personal eBird export.",
	}, func(ctx context.Context, req *sdk.ReadResourceRequest) (*sdk.ReadResourceResult, error) {
		payload := struct {
			Meta           any               `json:"meta"`
			SpeciesIndex   []it.SpeciesIndex `json:"speciesIndex,omitempty"`
			CountSightings int               `json:"countSightings"`
		}{
			Meta:           pc.Meta,
			SpeciesIndex:   pc.SpeciesIndex,
			CountSightings: len(pc.Sightings),
		}

		buf, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			return nil, err
		}

		return &sdk.ReadResourceResult{
			Contents: []*sdk.ResourceContents{
				{
					URI:      personalURI,
					MIMEType: "application/json",
					Text:     string(buf),
				},
			},
		}, nil
	})
}
