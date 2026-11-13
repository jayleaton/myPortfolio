package services

//go:generate go tool mockgen -source=portfolio.go -destination=mocks/portfolio.go -package=mocks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jayleaton/myPortfolio/apps/web/internal/content"
	"github.com/jayleaton/myPortfolio/packages/database/model"
)

// The portfolio's content types, re-exported so the API and UI share them.
type (
	Portfolio = content.Portfolio
	Project   = content.Project
	Link      = content.Link
	Metric    = content.Metric
	Video     = content.Video
	Site      = content.Site
	Section   = content.Section
	// CacheEntry is one persisted cache row (last good live values).
	CacheEntry = model.CacheEntry
)

// ContentSource is where the portfolio's copy comes from (content.Source).
type ContentSource interface {
	Portfolio() Portfolio
}

// Fetcher reads a public URL and returns its body. The adapter
// (server.Fetcher) sends a generic User-Agent, bounds the time and the body
// size, and returns an error for any non-2xx status.
type Fetcher interface {
	Fetch(ctx context.Context, url string) ([]byte, error)
}

// CacheStore persists the last good live values across restarts
// (database.CacheRepository). GetCache returns ErrRecordNotFound for a
// missing key.
type CacheStore interface {
	GetCache(ctx context.Context, key string) (CacheEntry, error)
	PutCache(ctx context.Context, entry CacheEntry) error
}

// Followers is the X follower count as last fetched. Known is false until
// a fetch has succeeded once (here or, through the cache, before a
// restart); the UI hides the number then.
type Followers struct {
	Handle    string    `json:"handle"`
	Count     int       `json:"count"`
	Known     bool      `json:"known"`
	Source    string    `json:"source,omitempty"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Videos is the channel's latest uploads as last fetched. Live is false
// while they are the content's fallback list.
type Videos struct {
	Items     []Video   `json:"items"`
	Live      bool      `json:"live"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// PortfolioPage is everything the home page renders.
type PortfolioPage struct {
	Portfolio Portfolio `json:"portfolio"`
	Videos    Videos    `json:"videos"`
	Followers Followers `json:"followers"`
}

// Cache keys of the persisted live values.
const (
	followersCacheKey = "x_followers"
	videosCacheKey    = "youtube_videos"
)

// ErrNoSource is returned by a refresh when every source failed; the last
// good value stays.
var ErrNoSource = errors.New("portfolio: every source failed")

// PortfolioService serves the portfolio and keeps its live values (X
// followers, latest videos) fresh in memory. Reads never touch the network:
// Run refreshes in the background, and a failed refresh keeps the last good
// value.
type PortfolioService struct {
	source ContentSource
	fetch  Fetcher
	cache  CacheStore
	now    func() time.Time

	mu        sync.RWMutex
	followers Followers
	videos    Videos
}

// NewPortfolioService returns the service; now defaults to time.Now.
func NewPortfolioService(source ContentSource, fetch Fetcher, cache CacheStore, now func() time.Time) *PortfolioService {
	if now == nil {
		now = time.Now
	}
	portfolio := source.Portfolio()
	return &PortfolioService{
		source:    source,
		fetch:     fetch,
		cache:     cache,
		now:       now,
		followers: Followers{Handle: portfolio.XHandle},
		videos:    Videos{Items: limitVideos(portfolio.FallbackVideos)},
	}
}

// Page is the home page's data.
func (service *PortfolioService) Page(ctx context.Context) PortfolioPage {
	return PortfolioPage{Portfolio: service.source.Portfolio(), Videos: service.Videos(ctx), Followers: service.Followers(ctx)}
}

// Followers is the last good X follower count.
func (service *PortfolioService) Followers(_ context.Context) Followers {
	service.mu.RLock()
	defer service.mu.RUnlock()
	return service.followers
}

// Videos is the last good list of latest uploads (or the fallback list).
func (service *PortfolioService) Videos(_ context.Context) Videos {
	service.mu.RLock()
	defer service.mu.RUnlock()
	videos := service.videos
	videos.Items = append([]Video(nil), videos.Items...)
	return videos
}

// Restore loads the values persisted before a restart. A missing entry is
// not an error; a broken one is reported and skipped.
func (service *PortfolioService) Restore(ctx context.Context) error {
	var errs []error
	var followers Followers
	if err := service.load(ctx, followersCacheKey, &followers); err != nil {
		errs = append(errs, err)
	} else if followers.Known && followers.Handle == service.source.Portfolio().XHandle {
		service.mu.Lock()
		service.followers = followers
		service.mu.Unlock()
	}
	var videos Videos
	if err := service.load(ctx, videosCacheKey, &videos); err != nil {
		errs = append(errs, err)
	} else if videos.Live && len(videos.Items) > 0 {
		videos.Items = limitVideos(videos.Items)
		service.mu.Lock()
		service.videos = videos
		service.mu.Unlock()
	}
	return errors.Join(errs...)
}

// load decodes the cache entry at key into into; a missing key leaves it zero.
func (service *PortfolioService) load(ctx context.Context, key string, into any) error {
	entry, err := service.cache.GetCache(ctx, key)
	if errors.Is(err, ErrRecordNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("portfolio: read cache %s: %w", key, err)
	}
	if err := json.Unmarshal(entry.Value, into); err != nil {
		return fmt.Errorf("portfolio: decode cache %s: %w", key, err)
	}
	return nil
}

// store persists value at key. The in-memory value is already updated, so a
// failure only costs the value after a restart.
func (service *PortfolioService) store(ctx context.Context, key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("portfolio: encode cache %s: %w", key, err)
	}
	if err := service.cache.PutCache(ctx, CacheEntry{Key: key, Value: data, UpdatedAt: service.now()}); err != nil {
		return fmt.Errorf("portfolio: write cache %s: %w", key, err)
	}
	return nil
}

// RefreshFollowers tries each follower source in order and keeps the first
// count one returns. When all fail it returns ErrNoSource (joined with each
// source's error) and the last good count stays.
func (service *PortfolioService) RefreshFollowers(ctx context.Context) error {
	handle := service.source.Portfolio().XHandle
	if handle == "" {
		return nil
	}
	errs := []error{ErrNoSource}
	for _, source := range FollowerSources {
		body, err := service.fetch.Fetch(ctx, source.URL(handle))
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", source.Name, err))
			continue
		}
		count, err := source.Parse(body, handle)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", source.Name, err))
			continue
		}
		followers := Followers{Handle: handle, Count: count, Known: true, Source: source.Name, UpdatedAt: service.now()}
		service.mu.Lock()
		service.followers = followers
		service.mu.Unlock()
		return service.store(ctx, followersCacheKey, followers)
	}
	return errors.Join(errs...)
}

// RefreshVideos tries each video source in order and keeps the first
// non-empty list. When all fail it returns ErrNoSource and the last good
// list (or the fallback list) stays.
func (service *PortfolioService) RefreshVideos(ctx context.Context) error {
	channel := service.source.Portfolio().YouTube.ID
	if channel == "" {
		return nil
	}
	errs := []error{ErrNoSource}
	for _, source := range VideoSources {
		body, err := service.fetch.Fetch(ctx, source.URL(channel))
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", source.Name, err))
			continue
		}
		items, err := source.Parse(body)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", source.Name, err))
			continue
		}
		videos := Videos{Items: limitVideos(items), Live: true, UpdatedAt: service.now()}
		service.mu.Lock()
		service.videos = videos
		service.mu.Unlock()
		return service.store(ctx, videosCacheKey, videos)
	}
	return errors.Join(errs...)
}

// RefreshIntervals are how often Run refreshes each live value.
type RefreshIntervals struct {
	Followers time.Duration
	Videos    time.Duration
}

// DefaultRefreshIntervals: followers every 15 minutes, videos hourly. Both
// sources are free public endpoints; staying well under their limits keeps
// them working.
var DefaultRefreshIntervals = RefreshIntervals{Followers: followersEvery, Videos: time.Hour}

// followersEvery is the follower count's refresh interval.
const followersEvery = 15 * time.Minute

// Run restores the persisted values, refreshes both live values at once,
// then on every interval until ctx ends. Failures go to report (task is
// "restore", "followers" or "videos") and never stop the loop.
func (service *PortfolioService) Run(ctx context.Context, intervals RefreshIntervals, report func(task string, err error)) {
	if report == nil {
		report = func(string, error) {}
	}
	if intervals.Followers <= 0 {
		intervals.Followers = DefaultRefreshIntervals.Followers
	}
	if intervals.Videos <= 0 {
		intervals.Videos = DefaultRefreshIntervals.Videos
	}
	if err := service.Restore(ctx); err != nil {
		report("restore", err)
	}
	refreshFollowers := func() {
		if err := service.RefreshFollowers(ctx); err != nil && ctx.Err() == nil {
			report("followers", err)
		}
	}
	refreshVideos := func() {
		if err := service.RefreshVideos(ctx); err != nil && ctx.Err() == nil {
			report("videos", err)
		}
	}
	refreshFollowers()
	refreshVideos()
	followersTick := time.NewTicker(intervals.Followers)
	defer followersTick.Stop()
	videosTick := time.NewTicker(intervals.Videos)
	defer videosTick.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-followersTick.C:
			refreshFollowers()
		case <-videosTick.C:
			refreshVideos()
		}
	}
}
