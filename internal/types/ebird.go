package ebird

// Personal index (from your CSV export → normalized JSON)
type SpeciesIndex struct {
	SpeciesCode     string   `json:"speciesCode"`
	CommonName      string   `json:"commonName"`
	SciName         string   `json:"sciName"`
	FirstSeen       string   `json:"firstSeen"`
	LastSeen        string   `json:"lastSeen"`
	TotalChecklists int      `json:"totalChecklists"`
	Locations       []string `json:"locations"`
}

// Minimal “recent obs” row you can fill from eBird API
type RecentObservation struct {
	SpeciesCode string `json:"speciesCode"`
	CommonName  string `json:"comName"`
	SciName     string `json:"sciName"`
	LocName     string `json:"locName"`
	LocID       string `json:"locId"`
	ObsDt       string `json:"obsDt"`
	HeardOnly   bool   `json:"howr,omitempty"`
}

// You’ll add: CSV loader, nearest-hotspots fetcher, frequency calculator, etc.
