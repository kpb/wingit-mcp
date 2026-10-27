package main

import (
	"context"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kpb/wingit-mcp/internal/ebird"
	"github.com/kpb/wingit-mcp/internal/tools"
	it "github.com/kpb/wingit-mcp/internal/types"
)

func TestParseLatLon(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		wantLat float64
		wantLon float64
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
		{
			name:    "nan_lat",
			input:   "NaN,-76.45",
			wantErr: true,
		},
		{
			name:    "infinite_lon",
			input:   "42.47,+Inf",
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

func TestBuildTargetChecklistUsesNormalizedQueryForFetchAndFilters(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		radiusKm   float64
		daysBack   int
		wantDist   string
		wantBack   string
		wantRadius float64
		wantDays   int
	}{
		{name: "defaults", wantDist: "20", wantBack: "7", wantRadius: 20, wantDays: 7},
		{name: "explicit", radiusKm: 12.5, daysBack: 4, wantDist: "12.5", wantBack: "4", wantRadius: 12.5, wantDays: 4},
		{name: "non-positive defaults", radiusKm: -1, daysBack: -2, wantDist: "20", wantBack: "7", wantRadius: 20, wantDays: 7},
		{name: "API maximums", radiusKm: 75, daysBack: 60, wantDist: "50", wantBack: "30", wantRadius: 50, wantDays: 30},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if got := r.URL.Query().Get("dist"); got != tt.wantDist {
					t.Errorf("dist = %q, want %q", got, tt.wantDist)
				}
				if got := r.URL.Query().Get("back"); got != tt.wantBack {
					t.Errorf("back = %q, want %q", got, tt.wantBack)
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`[]`))
			}))
			defer server.Close()

			client := ebird.NewClient("token", ebird.WithBaseURL(server.URL), ebird.WithMinInterval(0))
			args := tools.TargetArgs{
				Location: "42.47,-76.45",
				RadiusKm: tt.radiusKm,
				DaysBack: tt.daysBack,
			}
			out, err := buildTargetChecklist(context.Background(), args, nil, func(ctx context.Context, normalized tools.TargetArgs) ([]it.RecentObservation, error) {
				lat, lon, err := parseLatLon(normalized.Location)
				if err != nil {
					return nil, err
				}
				return client.RecentNearby(ctx, lat, lon, normalized.RadiusKm, normalized.DaysBack, 0)
			})
			if err != nil {
				t.Fatalf("buildTargetChecklist() error = %v", err)
			}
			if out.Filters.RadiusKm != tt.wantRadius {
				t.Errorf("Filters.RadiusKm = %v, query dist = %s", out.Filters.RadiusKm, tt.wantDist)
			}
			if out.Filters.DaysBack != tt.wantDays {
				t.Errorf("Filters.DaysBack = %v, query back = %s", out.Filters.DaysBack, tt.wantBack)
			}
		})
	}
}

func TestBuildTargetChecklistRejectsInvalidNumericArgsBeforeLoading(t *testing.T) {
	t.Parallel()

	for _, args := range []tools.TargetArgs{
		{Location: "42,-76", RadiusKm: math.Inf(1)},
		{Location: "42,-76", MinFrequency: math.NaN()},
	} {
		_, err := buildTargetChecklist(context.Background(), args, nil, func(context.Context, tools.TargetArgs) ([]it.RecentObservation, error) {
			t.Fatal("loader called for invalid arguments")
			return nil, nil
		})
		if err == nil {
			t.Fatal("expected error for non-finite numeric argument")
		}
	}
}

func TestRecentObservationSourceOfflineDoesNotRequireLiveClient(t *testing.T) {
	t.Setenv("WINGIT_EBIRD_TOKEN", "")
	path := filepath.Join(t.TempDir(), "recent.json")
	if err := os.WriteFile(path, []byte(`[]`), 0o600); err != nil {
		t.Fatal(err)
	}

	source := recentObservationSource{
		offlinePath: path,
		liveClient: newLiveClientProvider(func() (*ebird.Client, error) {
			return ebird.NewClientFromEnv()
		}),
	}
	rows, err := source.load(context.Background(), tools.TargetArgs{Location: "offline demo"})
	if err != nil {
		t.Fatalf("offline load error = %v", err)
	}
	if len(rows) != 0 {
		t.Fatalf("loaded %d observations, want 0", len(rows))
	}
}

func TestLiveClientProviderSharesPacingAcrossCalls(t *testing.T) {
	for _, concurrent := range []bool{false, true} {
		name := "sequential"
		if concurrent {
			name = "concurrent"
		}
		t.Run(name, func(t *testing.T) {
			const minInterval = 120 * time.Millisecond
			requestTimes := make(chan time.Time, 2)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requestTimes <- time.Now()
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`[]`))
			}))
			defer server.Close()

			factoryCalls := 0
			provider := newLiveClientProvider(func() (*ebird.Client, error) {
				factoryCalls++
				return ebird.NewClient("test-token", ebird.WithBaseURL(server.URL), ebird.WithMinInterval(minInterval)), nil
			})
			source := recentObservationSource{liveClient: provider}
			call := func() error {
				_, err := buildTargetChecklist(context.Background(), tools.TargetArgs{Location: "42.47,-76.45"}, nil, source.load)
				return err
			}

			if concurrent {
				var wg sync.WaitGroup
				errs := make(chan error, 2)
				for range 2 {
					wg.Add(1)
					go func() {
						defer wg.Done()
						errs <- call()
					}()
				}
				wg.Wait()
				close(errs)
				for err := range errs {
					if err != nil {
						t.Fatalf("buildTargetChecklist() error = %v", err)
					}
				}
			} else {
				if err := call(); err != nil {
					t.Fatalf("first buildTargetChecklist() error = %v", err)
				}
				if err := call(); err != nil {
					t.Fatalf("second buildTargetChecklist() error = %v", err)
				}
			}

			first, second := <-requestTimes, <-requestTimes
			if delta := second.Sub(first); delta < minInterval-40*time.Millisecond {
				t.Fatalf("requests were paced too closely: %v < %v", delta, minInterval-40*time.Millisecond)
			}
			if factoryCalls != 1 {
				t.Fatalf("client factory called %d times, want 1", factoryCalls)
			}
		})
	}
}

func TestLiveClientProviderCancellationDoesNotSendRequest(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer server.Close()

	provider := newLiveClientProvider(func() (*ebird.Client, error) {
		return ebird.NewClient("test-token", ebird.WithBaseURL(server.URL), ebird.WithMinInterval(time.Second)), nil
	})
	source := recentObservationSource{liveClient: provider}
	args := tools.TargetArgs{Location: "42.47,-76.45"}
	if _, err := buildTargetChecklist(context.Background(), args, nil, source.load); err != nil {
		t.Fatalf("first buildTargetChecklist() error = %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	start := time.Now()
	_, err := buildTargetChecklist(ctx, args, nil, source.load)
	if err == nil {
		t.Fatal("canceled call returned no error")
	}
	if elapsed := time.Since(start); elapsed > 200*time.Millisecond {
		t.Fatalf("canceled call took %v, want prompt cancellation", elapsed)
	}
	if got := requests.Load(); got != 1 {
		t.Fatalf("server received %d requests, want only the first", got)
	}
}
