package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// fetchTimeout bounds one outbound read of a public page (the X and
	// YouTube sources); a slow source is skipped, never waited on.
	fetchTimeout = 10 * time.Second
	// maxFetchBody caps a fetched page: a YouTube channel page is ~1 MiB.
	maxFetchBody = 4 << 20
	// fetchUserAgent is a generic browser User-Agent: the sources answer
	// browsers, and it names neither the site nor its owner.
	fetchUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0 Safari/537.36"
)

// Fetcher implements services.Fetcher over HTTP for the fixed public URLs
// the portfolio reads (never user-supplied URLs: those go through
// packages/storage/safefetch).
type Fetcher struct {
	client *http.Client
}

// NewFetcher returns a fetcher with its own short-timeout client.
func NewFetcher() *Fetcher {
	return &Fetcher{client: &http.Client{Timeout: fetchTimeout}}
}

// Fetch GETs url and returns the body; any non-2xx status is an error.
func (fetcher *Fetcher) Fetch(ctx context.Context, url string) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", fetchUserAgent)
	request.Header.Set("Accept-Language", "en")
	response, err := fetcher.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("GET %s: %s", url, response.Status)
	}
	body, err := io.ReadAll(io.LimitReader(response.Body, maxFetchBody+1))
	if err != nil {
		return nil, err
	}
	if len(body) > maxFetchBody {
		return nil, fmt.Errorf("GET %s: body over %d bytes", url, maxFetchBody)
	}
	return body, nil
}
