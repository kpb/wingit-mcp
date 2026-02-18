package main

import "testing"

func TestParseLatLon(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		wantErr  bool
		wantLat  float64
		wantLon  float64
	}{
		{
			name:    "valid",
			input:   "42.47,-76.45",
			wantLat: 42.47,
			wantLon: -76.45,
		},
		{
			name:    "spaces",
			input:   " 42.47 , -76.45 ",
			wantLat: 42.47,
			wantLon: -76.45,
		},
		{
			name:    "missing_comma",
			input:   "42.47 -76.45",
			wantErr: true,
		},
		{
			name:    "bad_lat",
			input:   "nope,-76.45",
			wantErr: true,
		},
		{
			name:    "bad_lon",
			input:   "42.47,nope",
			wantErr: true,
		},
		{
			name:    "lat_out_of_range",
			input:   "91,-76.45",
			wantErr: true,
		},
		{
			name:    "lon_out_of_range",
			input:   "42.47,-181",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			lat, lon, err := parseLatLon(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if lat != tt.wantLat || lon != tt.wantLon {
				t.Fatalf("got lat=%v lon=%v, want lat=%v lon=%v", lat, lon, tt.wantLat, tt.wantLon)
			}
		})
	}
}
