package services_test

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sync"
	"testing"
	"testing/synctest"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/jayleaton/myPortfolio/apps/web/internal/content"
	"github.com/jayleaton/myPortfolio/apps/web/internal/services"
	"github.com/jayleaton/myPortfolio/apps/web/internal/services/mocks"
)

var (
	anyCtx      = gomock.Any()
	fetchedAt   = time.Date(2026, time.November, 13, 4, 30, 0, 0, time.UTC)
	errOffline  = errors.New("dial tcp: offline")
	errDatabase = errors.New("database down")
	fallback    = []services.Video{{ID: "fallback001", Title: "Fallback one"}, {ID: "fallback002", Title: "Fallback two"}}
)

func fixture(test *testing.T, name string) []byte {
	test.Helper()
	body, err := os.ReadFile("testdata/" + name)
	if err != nil {
		test.Fatal(err)
	}
	return body
}

// testContent is a ContentSource with the handle and channel the fixtures
// were recorded for.
type testContent struct{ handle, channel string }

func (source testContent) Portfolio() services.Portfolio {
	return services.Portfolio{
		Name:           "JJ Eaton",
		XHandle:        source.handle,
		YouTube:        content.Channel{ID: source.channel},
		FallbackVideos: fallback,
	}
}

var jj = testContent{handle: "jayleaton", channel: "UC-AGtYCONF0Yu0uK_fL1j1Q"}

type portfolioPorts struct {
	fetch *mocks.MockFetcherMockRecorder
	cache *mocks.MockCacheStoreMockRecorder
}

func newPortfolioService(test *testing.T, source services.ContentSource) (*services.PortfolioService, portfolioPorts) {
	controller := gomock.NewController(test)
	fetch, cache := mocks.NewMockFetcher(controller), mocks.NewMockCacheStore(controller)
	service := services.NewPortfolioService(source, fetch, cache, func() time.Time { return fetchedAt })
	return service, portfolioPorts{fetch: fetch.EXPECT(), cache: cache.EXPECT()}
}

func sourceURL(name, key string) string {
	for _, source := range services.FollowerSources {
		if source.Name == name {
			return source.URL(key)
		}
	}
	for _, source := range services.VideoSources {
		if source.Name == name {
			return source.URL(key)
		}
	}
	panic("no source " + name)
}

func cached(test *testing.T, key string, value any) services.CacheEntry {
	test.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		test.Fatal(err)
	}
	return services.CacheEntry{Key: key, Value: data, UpdatedAt: fetchedAt}
}

// Before any refresh the page has the content, the fallback videos and an
// unknown follower count (the UI hides it).
func TestPortfolioPageBeforeRefresh(test *testing.T) {
	test.Parallel()
	service, _ := newPortfolioService(test, jj)
	page := service.Page(ctx)
	if page.Portfolio.Name != "JJ Eaton" || page.Followers.Known || page.Followers.Handle != "jayleaton" {
		test.Fatalf("page = %+v", page)
	}
	if page.Videos.Live || len(page.Videos.Items) != len(fallback) {
		test.Fatalf("videos = %+v, want the fallback list", page.Videos)
	}
	page.Videos.Items[0].Title = "changed by a caller"
	if service.Videos(ctx).Items[0].Title != "Fallback one" {
		test.Fatal("Videos shares its slice with callers")
	}
}

// RefreshFollowers keeps the first source that answers with a count for
// the right account, persists it, and keeps the last good count when every
// source fails.
func TestRefreshFollowers(test *testing.T) {
	test.Parallel()
	fx, vx := sourceURL("fxtwitter", "jayleaton"), sourceURL("vxtwitter", "jayleaton")
	cases := map[string]struct {
		expect    func(test *testing.T, ports portfolioPorts)
		wantErr   error
		wantCount int
		wantFrom  string
	}{
		"fxtwitter answers": {
			expect: func(test *testing.T, ports portfolioPorts) {
				ports.fetch.Fetch(anyCtx, fx).Return(fixture(test, "fxtwitter.json"), nil)
				ports.cache.PutCache(anyCtx, gomock.Any()).DoAndReturn(func(_ context.Context, entry services.CacheEntry) error {
					var stored services.Followers
					if err := json.Unmarshal(entry.Value, &stored); err != nil || entry.Key != "x_followers" || stored.Count != 2367 || !entry.UpdatedAt.Equal(fetchedAt) {
						test.Errorf("stored %s = %s (%v)", entry.Key, entry.Value, err)
					}
					return nil
				})
			},
			wantCount: 2367, wantFrom: "fxtwitter",
		},
		"fxtwitter down, vxtwitter answers": {
			expect: func(test *testing.T, ports portfolioPorts) {
				ports.fetch.Fetch(anyCtx, fx).Return(nil, errOffline)
				ports.fetch.Fetch(anyCtx, vx).Return(fixture(test, "vxtwitter.json"), nil)
				ports.cache.PutCache(anyCtx, gomock.Any()).Return(nil)
			},
			wantCount: 2367, wantFrom: "vxtwitter",
		},
		"fxtwitter changed its shape": {
			expect: func(test *testing.T, ports portfolioPorts) {
				ports.fetch.Fetch(anyCtx, fx).Return([]byte(`<html>Just a moment...</html>`), nil)
				ports.fetch.Fetch(anyCtx, vx).Return(nil, errOffline)
			},
			wantErr: services.ErrNoSource,
		},
		"stored but not persisted": {
			expect: func(test *testing.T, ports portfolioPorts) {
				ports.fetch.Fetch(anyCtx, fx).Return(fixture(test, "fxtwitter.json"), nil)
				ports.cache.PutCache(anyCtx, gomock.Any()).Return(errDatabase)
			},
			wantErr: errDatabase, wantCount: 2367, wantFrom: "fxtwitter",
		},
	}
	for name, testCase := range cases {
		test.Run(name, func(test *testing.T) {
			test.Parallel()
			service, ports := newPortfolioService(test, jj)
			testCase.expect(test, ports)
			err := service.RefreshFollowers(ctx)
			if !errors.Is(err, testCase.wantErr) || (testCase.wantErr == nil && err != nil) {
				test.Fatalf("RefreshFollowers = %v, want %v", err, testCase.wantErr)
			}
			got := service.Followers(ctx)
			if got.Known != (testCase.wantCount > 0) || got.Count != testCase.wantCount || got.Source != testCase.wantFrom {
				test.Fatalf("Followers = %+v", got)
			}
		})
	}
}

// A failed refresh after a good one keeps the good count.
func TestRefreshFollowersKeepsLastGood(test *testing.T) {
	test.Parallel()
	service, ports := newPortfolioService(test, jj)
	ports.fetch.Fetch(anyCtx, sourceURL("fxtwitter", "jayleaton")).Return(fixture(test, "fxtwitter.json"), nil)
	ports.cache.PutCache(anyCtx, gomock.Any()).Return(nil)
	if err := service.RefreshFollowers(ctx); err != nil {
		test.Fatal(err)
	}
	ports.fetch.Fetch(anyCtx, gomock.Any()).Return(nil, errOffline).Times(len(services.FollowerSources))
	if err := service.RefreshFollowers(ctx); !errors.Is(err, services.ErrNoSource) || !errors.Is(err, errOffline) {
		test.Fatalf("RefreshFollowers = %v", err)
	}
	if got := service.Followers(ctx); !got.Known || got.Count != 2367 {
		test.Fatalf("Followers = %+v, want the last good count", got)
	}
}

func TestRefreshVideos(test *testing.T) {
	test.Parallel()
	page, feed := sourceURL("channel page", jj.channel), sourceURL("rss", jj.channel)
	cases := map[string]struct {
		expect    func(test *testing.T, ports portfolioPorts)
		wantErr   error
		wantFirst string
		wantLive  bool
	}{
		"channel page": {
			expect: func(test *testing.T, ports portfolioPorts) {
				ports.fetch.Fetch(anyCtx, page).Return(fixture(test, "youtube-videos.html"), nil)
				ports.cache.PutCache(anyCtx, gomock.Any()).DoAndReturn(func(_ context.Context, entry services.CacheEntry) error {
					var stored services.Videos
					if err := json.Unmarshal(entry.Value, &stored); err != nil || entry.Key != "youtube_videos" || len(stored.Items) != services.VideoLimit {
						test.Errorf("stored %s = %s (%v)", entry.Key, entry.Value, err)
					}
					return nil
				})
			},
			wantFirst: "WCrknrNsmjQ", wantLive: true,
		},
		"page blocked, feed answers": {
			expect: func(test *testing.T, ports portfolioPorts) {
				ports.fetch.Fetch(anyCtx, page).Return([]byte("<html>consent</html>"), nil)
				ports.fetch.Fetch(anyCtx, feed).Return(fixture(test, "youtube-feed.xml"), nil)
				ports.cache.PutCache(anyCtx, gomock.Any()).Return(nil)
			},
			wantFirst: "WCrknrNsmjQ", wantLive: true,
		},
		"both down keeps the fallback": {
			expect: func(test *testing.T, ports portfolioPorts) {
				ports.fetch.Fetch(anyCtx, page).Return(nil, errOffline)
				ports.fetch.Fetch(anyCtx, feed).Return(nil, errOffline)
			},
			wantErr: services.ErrNoSource, wantFirst: "fallback001",
		},
	}
	for name, testCase := range cases {
		test.Run(name, func(test *testing.T) {
			test.Parallel()
			service, ports := newPortfolioService(test, jj)
			testCase.expect(test, ports)
			if err := service.RefreshVideos(ctx); !errors.Is(err, testCase.wantErr) || (testCase.wantErr == nil && err != nil) {
				test.Fatalf("RefreshVideos = %v, want %v", err, testCase.wantErr)
			}
			got := service.Videos(ctx)
			if got.Live != testCase.wantLive || got.Items[0].ID != testCase.wantFirst || len(got.Items) > services.VideoLimit {
				test.Fatalf("Videos = %+v", got)
			}
		})
	}
}

// Without a handle or channel there is nothing to fetch.
func TestRefreshWithoutAccounts(test *testing.T) {
	test.Parallel()
	service, _ := newPortfolioService(test, testContent{})
	if err := service.RefreshFollowers(ctx); err != nil {
		test.Fatal(err)
	}
	if err := service.RefreshVideos(ctx); err != nil {
		test.Fatal(err)
	}
}

// Restore serves what was persisted before a restart, ignores missing and
// foreign entries, and reports broken ones.
func TestRestore(test *testing.T) {
	test.Parallel()
	good := services.Followers{Handle: "jayleaton", Count: 2300, Known: true, Source: "fxtwitter", UpdatedAt: fetchedAt}
	videos := services.Videos{Items: []services.Video{{ID: "WCrknrNsmjQ", Title: "Cached"}}, Live: true, UpdatedAt: fetchedAt}
	cases := map[string]struct {
		followers, videos func(test *testing.T) (services.CacheEntry, error)
		wantErr           error
		wantCount         int
		wantVideo         string
	}{
		"both cached": {
			followers: func(test *testing.T) (services.CacheEntry, error) { return cached(test, "x_followers", good), nil },
			videos:    func(test *testing.T) (services.CacheEntry, error) { return cached(test, "youtube_videos", videos), nil },
			wantCount: 2300, wantVideo: "WCrknrNsmjQ",
		},
		"nothing cached": {
			followers: func(*testing.T) (services.CacheEntry, error) {
				return services.CacheEntry{}, services.ErrRecordNotFound
			},
			videos: func(*testing.T) (services.CacheEntry, error) {
				return services.CacheEntry{}, services.ErrRecordNotFound
			},
			wantVideo: "fallback001",
		},
		"another account's count": {
			followers: func(test *testing.T) (services.CacheEntry, error) {
				other := good
				other.Handle = "someone"
				return cached(test, "x_followers", other), nil
			},
			videos: func(*testing.T) (services.CacheEntry, error) {
				return services.CacheEntry{}, services.ErrRecordNotFound
			},
			wantVideo: "fallback001",
		},
		"broken and unreachable": {
			followers: func(*testing.T) (services.CacheEntry, error) {
				return services.CacheEntry{Key: "x_followers", Value: []byte("{")}, nil
			},
			videos:  func(*testing.T) (services.CacheEntry, error) { return services.CacheEntry{}, errDatabase },
			wantErr: errDatabase, wantVideo: "fallback001",
		},
	}
	for name, testCase := range cases {
		test.Run(name, func(test *testing.T) {
			test.Parallel()
			service, ports := newPortfolioService(test, jj)
			ports.cache.GetCache(anyCtx, "x_followers").Return(testCase.followers(test))
			ports.cache.GetCache(anyCtx, "youtube_videos").Return(testCase.videos(test))
			err := service.Restore(ctx)
			if !errors.Is(err, testCase.wantErr) || (testCase.wantErr == nil && err != nil) {
				test.Fatalf("Restore = %v, want %v", err, testCase.wantErr)
			}
			if got := service.Followers(ctx); got.Count != testCase.wantCount || got.Known != (testCase.wantCount > 0) {
				test.Fatalf("Followers = %+v", got)
			}
			if got := service.Videos(ctx); got.Items[0].ID != testCase.wantVideo {
				test.Fatalf("Videos = %+v", got)
			}
		})
	}
	test.Run("broken entry names its key", func(test *testing.T) {
		test.Parallel()
		service, ports := newPortfolioService(test, jj)
		ports.cache.GetCache(anyCtx, "x_followers").Return(services.CacheEntry{Value: []byte("{")}, nil)
		ports.cache.GetCache(anyCtx, "youtube_videos").Return(services.CacheEntry{}, services.ErrRecordNotFound)
		if err := service.Restore(ctx); err == nil || !strings.Contains(err.Error(), "x_followers") {
			test.Fatalf("Restore = %v", err)
		}
	})
}

// Run restores, refreshes both values at once, then again on each tick,
// reporting failures without stopping, until its context ends.
func TestRun(test *testing.T) {
	test.Parallel()
	synctest.Test(test, func(test *testing.T) {
		service, ports := newPortfolioService(test, jj)
		ports.cache.GetCache(anyCtx, gomock.Any()).Return(services.CacheEntry{}, errDatabase).Times(2)
		var mu sync.Mutex
		fetches := map[string]int{}
		ports.fetch.Fetch(anyCtx, gomock.Any()).DoAndReturn(func(_ context.Context, url string) ([]byte, error) {
			mu.Lock()
			defer mu.Unlock()
			fetches[url]++
			if url == sourceURL("fxtwitter", "jayleaton") {
				return fixture(test, "fxtwitter.json"), nil
			}
			return nil, errOffline
		}).AnyTimes()
		ports.cache.PutCache(anyCtx, gomock.Any()).Return(nil).AnyTimes()
		reported := map[string]int{}
		runCtx, cancel := context.WithCancel(test.Context())
		done := make(chan struct{})
		go func() {
			defer close(done)
			service.Run(runCtx, services.RefreshIntervals{Followers: time.Minute, Videos: 2 * time.Minute}, func(task string, _ error) {
				mu.Lock()
				defer mu.Unlock()
				reported[task]++
			})
		}()
		synctest.Wait()
		if !service.Followers(ctx).Known {
			test.Fatal("Run did not refresh followers at start")
		}
		time.Sleep(2 * time.Minute) //nolint:forbidigo // synctest's fake clock: advances instantly
		synctest.Wait()
		cancel()
		<-done
		mu.Lock()
		defer mu.Unlock()
		if got := fetches[sourceURL("fxtwitter", "jayleaton")]; got != 3 {
			test.Errorf("follower fetches = %d, want 3 (start + 2 ticks)", got)
		}
		if reported["restore"] != 1 || reported["videos"] != 2 || reported["followers"] != 0 {
			test.Errorf("reported = %v", reported)
		}
	})
}

// Run fills in default intervals and a nil reporter, and returns as soon as
// its context ends.
func TestRunDefaults(test *testing.T) {
	test.Parallel()
	synctest.Test(test, func(test *testing.T) {
		service, ports := newPortfolioService(test, testContent{})
		ports.cache.GetCache(anyCtx, gomock.Any()).Return(services.CacheEntry{}, services.ErrRecordNotFound).Times(2)
		runCtx, cancel := context.WithCancel(test.Context())
		cancel()
		service.Run(runCtx, services.RefreshIntervals{}, nil)
	})
}
