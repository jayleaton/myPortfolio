package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// FollowerSource is one free, unauthenticated place to read an X account's
// public follower count. The paid X API is deliberately not used.
type FollowerSource struct {
	Name  string
	URL   func(handle string) string
	Parse func(body []byte, handle string) (int, error)
}

// FollowerSources are tried in order. fxtwitter's public JSON API answers
// today; vxtwitter has the same shape of API and is kept as a fallback for
// when one of them is down or blocked.
var FollowerSources = []FollowerSource{
	{Name: "fxtwitter", URL: func(handle string) string { return "https://api.fxtwitter.com/" + url.PathEscape(handle) }, Parse: ParseFxTwitterFollowers},
	{Name: "vxtwitter", URL: func(handle string) string { return "https://api.vxtwitter.com/" + url.PathEscape(handle) }, Parse: ParseVxTwitterFollowers},
}

// maxFollowers bounds a parsed count: anything above is a parse of the
// wrong field, not a follower count (the largest accounts have ~2e8).
const maxFollowers = 2_000_000_000

// Errors a follower parser returns.
var (
	ErrUnexpectedBody = errors.New("unexpected response body")
	ErrWrongAccount   = errors.New("response is for another account")
)

// ParseFxTwitterFollowers reads user.followers from
// GET https://api.fxtwitter.com/<handle> (testdata/fxtwitter.json).
func ParseFxTwitterFollowers(body []byte, handle string) (int, error) {
	var response struct {
		Code int `json:"code"`
		User *struct {
			ScreenName string `json:"screen_name"`
			Followers  *int   `json:"followers"`
		} `json:"user"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return 0, fmt.Errorf("%w: %w", ErrUnexpectedBody, err)
	}
	if response.Code != 0 && response.Code != 200 {
		return 0, fmt.Errorf("%w: code %d", ErrUnexpectedBody, response.Code)
	}
	if response.User == nil || response.User.Followers == nil {
		return 0, fmt.Errorf("%w: no user.followers", ErrUnexpectedBody)
	}
	return checkFollowers(response.User.ScreenName, handle, *response.User.Followers)
}

// ParseVxTwitterFollowers reads followers_count from
// GET https://api.vxtwitter.com/<handle> (testdata/vxtwitter.json).
func ParseVxTwitterFollowers(body []byte, handle string) (int, error) {
	var response struct {
		ScreenName string `json:"screen_name"`
		Followers  *int   `json:"followers_count"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return 0, fmt.Errorf("%w: %w", ErrUnexpectedBody, err)
	}
	if response.Followers == nil {
		return 0, fmt.Errorf("%w: no followers_count", ErrUnexpectedBody)
	}
	return checkFollowers(response.ScreenName, handle, *response.Followers)
}

// checkFollowers accepts count only for the account asked for and within
// the plausible range.
func checkFollowers(screenName, handle string, count int) (int, error) {
	if !strings.EqualFold(screenName, handle) {
		return 0, fmt.Errorf("%w: %q", ErrWrongAccount, screenName)
	}
	if count < 0 || count > maxFollowers {
		return 0, fmt.Errorf("%w: follower count %d", ErrUnexpectedBody, count)
	}
	return count, nil
}
