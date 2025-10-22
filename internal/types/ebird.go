// internal/types/ebird.go
package types

type PersonalChecklist struct {
	Meta struct {
		Owner              string `json:"owner"`
		Source             string `json:"source"`
		GeneratedAt        string `json:"generatedAt"`
		TotalObservations  int    `json:"totalObservations"`
		TotalSpecies       int    `json:"totalSpecies"`
		FirstChecklistDate string `json:"firstChecklistDate"`
	} `json:"meta"`
	Sightings    []PersonalSighting `json:"sightings"`
	SpeciesIndex []SpeciesIndex     `json:"speciesIndex"`
}

type PersonalSighting struct {
	SpeciesCode        string  `json:"speciesCode"`
	CommonName         string  `json:"commonName"`
	SciName            string  `json:"sciName"`
	ObsDt              string  `json:"obsDt"`
	LocName            string  `json:"locName"`
	LocID              string  `json:"locId"`
	CountyCode         string  `json:"countyCode"`
	Lat                float64 `json:"lat"`
	Lng                float64 `json:"lng"`
	Count              int     `json:"count"`
	ObsValid           bool    `json:"obsValid"`
	ObsReviewed        bool    `json:"obsReviewed"`
	Media              bool    `json:"media"`
	EnteredAsHeardOnly bool    `json:"enteredAsHeardOnly"`
	ChecklistID        string  `json:"checklistId"`
}

type SpeciesIndex struct {
	SpeciesCode     string   `json:"speciesCode"`
	CommonName      string   `json:"commonName"`
	SciName         string   `json:"sciName"`
	FirstSeen       string   `json:"firstSeen"`
	LastSeen        string   `json:"lastSeen"`
	TotalChecklists int      `json:"totalChecklists"`
	TotalCount      int      `json:"totalCount"`
	Locations       []string `json:"locations"`
}

// Subset of eBird recent observations (matches fixture)
type RecentObservation struct {
	SpeciesCode string `json:"speciesCode"`
	CommonName  string `json:"comName"`
	SciName     string `json:"sciName"`
	LocName     string `json:"locName"`
	LocID       string `json:"locId"`
	ObsDt       string `json:"obsDt"`
	HeardOnly   bool   `json:"howr,omitempty"`
}
