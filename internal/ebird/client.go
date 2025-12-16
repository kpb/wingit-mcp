package ebird

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/kpb/wingit-mcp/internal/types"
)

const (
	defaultBaseURL = "https://api.ebird.org"
)

var (
	ErrUnauthorized = errors.New("ebird: unauthorized (bad token?)")
	ErrRateLimited  = errors.New("ebird: rate limited")
	ErrBadRequest   = errors.New("ebird: bad request")
)

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
	userAgent  string
}

type Option func(*Client)

func WithBaseURL(baseURL string) Option {
	return func(c *Client) { c.baseURL = strings.TrimRight(baseURL, "/") }
}

func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.httpClient = hc }
}

func WithUserAgent(ua string) Option {
	return func(c *Client) { c.userAgent = ua }
}

func NewClient(token string, opts ...Option) *Client {
	c := &Client{
		baseURL: defaultBaseURL,
		token:   strings.TrimSpace(token),
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		userAgent: "wingit-mcp/0.2.0 (+https://github.com/kpb/wingit-mcp)",
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// RecentNearby calls:
// GET /v2/data/obs/geo/recent?lat=...&lng=...&dist=...&back=...&maxResults=...&sort=...
//
// eBird constraints: back 1-30, dist 0-50km
func (c *Client) RecentNearby(
	ctx context.Context,
	lat, lng float64,
	distKm float64,
	backDays int,
	maxResults int,
) ([]types.RecentObservation, error) {
	if c.token == "" {
		return nil, fmt.Errorf("%w: missing API token", ErrUnauthorized)
	}

	// Clamp to eBird documented ranges (engine may normalize too).
	if backDays < 1 {
		backDays = 1
	} else if backDays > 30 {
		backDays = 30
	}
	if distKm < 0 {
		distKm = 0
	} else if distKm > 50 {
		distKm = 50
	}
	if maxResults < 1 {
		maxResults = 0 // omit
	}

	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("ebird: invalid baseURL: %w", err)
	}
	u.Path = strings.TrimRight(u.Path, "/") + "/v2/data/obs/geo/recent"

	q := u.Query()
	// Docs say “to 2 decimal places”
	q.Set("lat", strconv.FormatFloat(lat, 'f', 2, 64))
	q.Set("lng", strconv.FormatFloat(lng, 'f', 2, 64))
	q.Set("dist", strconv.FormatFloat(distKm, 'f', -1, 64))
	q.Set("back", strconv.Itoa(backDays))
	q.Set("sort", "date") // stable default; parameterize later
	if maxResults > 0 {
		q.Set("maxResults", strconv.Itoa(maxResults))
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("ebird: new request: %w", err)
	}
	req.Header.Set("X-eBirdApiToken", c.token) // required
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ebird: request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read limited body for error messages; decode stream for success.
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
		msg := strings.TrimSpace(string(b))

		switch resp.StatusCode {
		case http.StatusUnauthorized, http.StatusForbidden:
			return nil, fmt.Errorf("%w: %s", ErrUnauthorized, msg)
		case http.StatusTooManyRequests:
			return nil, fmt.Errorf("%w: %s", ErrRateLimited, msg)
		case http.StatusBadRequest:
			return nil, fmt.Errorf("%w: %s", ErrBadRequest, msg)
		default:
			return nil, fmt.Errorf("ebird: http %d: %s", resp.StatusCode, msg)
		}
	}

	var out []types.RecentObservation
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&out); err != nil {
		return nil, fmt.Errorf("ebird: decode response: %w", err)
	}
	return out, nil
}
