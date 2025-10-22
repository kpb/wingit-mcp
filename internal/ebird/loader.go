// internal/ebird/loader.go
package ebird

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	it "github.com/kpb/wingit-mcp/internal/types"
)

func LoadPersonalChecklist(path string) (*it.PersonalChecklist, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read personal checklist: %w", err)
	}
	var pc it.PersonalChecklist
	dec := json.NewDecoder(bytes.NewReader(b))
	if err := dec.Decode(&pc); err != nil {
		return nil, fmt.Errorf("decode personal checklist: %w", err)
	}
	return &pc, nil
}

func LoadRecentNearby(path string) ([]it.RecentObservation, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read recent nearby: %w", err)
	}
	var rows []it.RecentObservation
	dec := json.NewDecoder(bytes.NewReader(b))
	if err := dec.Decode(&rows); err != nil {
		return nil, fmt.Errorf("decode recent nearby: %w", err)
	}
	return rows, nil
}

func BuildPersonalSeenSet(pc *it.PersonalChecklist) map[string]struct{} {
	seen := make(map[string]struct{}, len(pc.SpeciesIndex))
	if len(pc.SpeciesIndex) > 0 {
		for _, s := range pc.SpeciesIndex {
			if s.SpeciesCode != "" {
				seen[s.SpeciesCode] = struct{}{}
			}
		}
		return seen
	}
	for _, s := range pc.Sightings {
		if s.SpeciesCode != "" {
			seen[s.SpeciesCode] = struct{}{}
		}
	}
	return seen
}
