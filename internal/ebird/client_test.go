package ebird

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kpb/wingit-mcp/internal/types"
)

func TestClient_RecentNearby_BuildsRequestAndDecodes(t *testing.T) {
	t.Parallel()

	const token = "test-token"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/v2/data/obs/geo/recent" {
			t.Fatalf("path = %s, want /v2/data/obs/geo/recent", r.URL.Path)
		}
		if got := r.Header.Get("X-eBirdApiToken"); got != token {
			t.Fatalf("X-eBirdApiToken = %q, want %q", got, token)
		}

		q := r.URL.Query()
		if q.Get("lat") == "" || q.Get("lng") == "" {
			t.Fatalf("lat/lng missing: %v", q)
		}
		if q.Get("back") != "7" {
			t.Fatalf("back = %q, want 7", q.Get("back"))
		}
		if q.Get("dist") != "20" {
			t.Fatalf("dist = %q, want 20", q.Get("dist"))
		}
		if q.Get("maxResults") != "123" {
			t.Fatalf("maxResults = %q, want 123", q.Get("maxResults"))
		}

		// Minimal sample aligned with eBird docs fields
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[
			{
				"speciesCode":"cangoo",
				"comName":"Canada Goose",
				"sciName":"Branta canadensis",
				"locId":"L1150539",
				"locName":"Hanshaw Rd. fields",
				"obsDt":"2017-08-23 20:05",
				"howMany":30,
				"lat":42.4663513,
				"lng":-76.4531064,
				"obsValid":true,
				"obsReviewed":false,
				"locationPrivate":false
			}
		]`))
	}))
	t.Cleanup(srv.Close)

	c := NewClient(token,
		WithBaseURL(srv.URL),
		WithHTTPClient(&http.Client{Timeout: 2 * time.Second}),
	)

	ctx := context.Background()
	obs, err := c.RecentNearby(ctx, 42.47, -76.45, 20, 7, 123)
	if err != nil {
		t.Fatalf("RecentNearby error: %v", err)
	}

	if len(obs) != 1 {
		t.Fatalf("len(obs)=%d, want 1", len(obs))
	}

	var _ types.RecentObservation = obs[0]
	if obs[0].SpeciesCode != "cangoo" {
		t.Fatalf("SpeciesCode=%q, want cangoo", obs[0].SpeciesCode)
	}
	if obs[0].CommonName != "Canada Goose" {
		t.Fatalf("CommonName=%q, want Canada Goose", obs[0].CommonName)
	}
	if obs[0].SciName != "Branta canadensis" {
		t.Fatalf("SciName=%q, want Branta canadensis", obs[0].SciName)
	}
	if obs[0].LocID != "L1150539" {
		t.Fatalf("LocID=%q, want L1150539", obs[0].LocID)
	}
	if obs[0].LocName != "Hanshaw Rd. fields" {
		t.Fatalf("LocName=%q, want Hanshaw Rd. fields", obs[0].LocName)
	}
	if obs[0].ObsDt != "2017-08-23 20:05" {
		t.Fatalf("ObsDt=%q, want 2017-08-23 20:05", obs[0].ObsDt)
	}

}

func TestClient_RecentNearby_Unauthorized(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("nope"))
	}))
	t.Cleanup(srv.Close)

	c := NewClient("bad-token", WithBaseURL(srv.URL))
	_, err := c.RecentNearby(context.Background(), 40, -105, 10, 7, 50)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if got := err.Error(); got == "" || err == nil {
		t.Fatalf("expected non-empty error")
	}
}

func TestClient_RecentNearby_HeardOnlyHowr(t *testing.T) {
	t.Parallel()

	const token = "test-token"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("X-eBirdApiToken"); got != token {
			t.Fatalf("X-eBirdApiToken = %q, want %q", got, token)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Include howr=true (heard-only) and make sure it decodes into HeardOnly.
		_, _ = w.Write([]byte(`[
			{
				"speciesCode":"amecro",
				"comName":"American Crow",
				"sciName":"Corvus brachyrhynchos",
				"locName":"Somewhere Nice",
				"locId":"L123",
				"obsDt":"2025-12-15 08:00",
				"howr": true
			}
		]`))
	}))
	t.Cleanup(srv.Close)

	c := NewClient(token, WithBaseURL(srv.URL))

	obs, err := c.RecentNearby(context.Background(), 40.00, -105.00, 10, 7, 25)
	if err != nil {
		t.Fatalf("RecentNearby error: %v", err)
	}
	if len(obs) != 1 {
		t.Fatalf("len(obs)=%d, want 1", len(obs))
	}

	if obs[0].SpeciesCode != "amecro" {
		t.Fatalf("SpeciesCode=%q, want amecro", obs[0].SpeciesCode)
	}
	if obs[0].CommonName != "American Crow" {
		t.Fatalf("CommonName=%q, want American Crow", obs[0].CommonName)
	}
	if !obs[0].HeardOnly {
		t.Fatalf("HeardOnly=%v, want true", obs[0].HeardOnly)
	}
}
