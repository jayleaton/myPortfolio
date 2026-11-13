package services_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/jayleaton/myPortfolio/apps/web/internal/services"
)

// The follower parsers read the recorded responses and refuse anything
// that is not a plausible count for the account asked for.
func TestFollowerParsers(test *testing.T) {
	test.Parallel()
	cases := map[string]struct {
		parse   func([]byte, string) (int, error)
		body    string
		handle  string
		want    int
		wantErr error
	}{
		"fxtwitter fixture":           {parse: services.ParseFxTwitterFollowers, body: "@fxtwitter.json", handle: "jayleaton", want: 2367},
		"fxtwitter handle case":       {parse: services.ParseFxTwitterFollowers, body: "@fxtwitter.json", handle: "JayLEaton", want: 2367},
		"fxtwitter other account":     {parse: services.ParseFxTwitterFollowers, body: "@fxtwitter.json", handle: "someone", wantErr: services.ErrWrongAccount},
		"fxtwitter not found":         {parse: services.ParseFxTwitterFollowers, body: `{"code":404,"message":"User not found"}`, handle: "jayleaton", wantErr: services.ErrUnexpectedBody},
		"fxtwitter no followers":      {parse: services.ParseFxTwitterFollowers, body: `{"code":200,"user":{"screen_name":"jayleaton"}}`, handle: "jayleaton", wantErr: services.ErrUnexpectedBody},
		"fxtwitter zero is a count":   {parse: services.ParseFxTwitterFollowers, body: `{"code":200,"user":{"screen_name":"jayleaton","followers":0}}`, handle: "jayleaton", want: 0},
		"fxtwitter negative":          {parse: services.ParseFxTwitterFollowers, body: `{"code":200,"user":{"screen_name":"jayleaton","followers":-1}}`, handle: "jayleaton", wantErr: services.ErrUnexpectedBody},
		"fxtwitter cloudflare page":   {parse: services.ParseFxTwitterFollowers, body: "<!DOCTYPE html><title>Just a moment...</title>", handle: "jayleaton", wantErr: services.ErrUnexpectedBody},
		"vxtwitter fixture":           {parse: services.ParseVxTwitterFollowers, body: "@vxtwitter.json", handle: "jayleaton", want: 2367},
		"vxtwitter no count":          {parse: services.ParseVxTwitterFollowers, body: `{"screen_name":"jayleaton"}`, handle: "jayleaton", wantErr: services.ErrUnexpectedBody},
		"vxtwitter server error page": {parse: services.ParseVxTwitterFollowers, body: "Internal Server Error", handle: "jayleaton", wantErr: services.ErrUnexpectedBody},
		"vxtwitter implausible count": {parse: services.ParseVxTwitterFollowers, body: `{"screen_name":"jayleaton","followers_count":99999999999}`, handle: "jayleaton", wantErr: services.ErrUnexpectedBody},
	}
	for name, testCase := range cases {
		test.Run(name, func(test *testing.T) {
			test.Parallel()
			body := []byte(testCase.body)
			if file, ok := strings.CutPrefix(testCase.body, "@"); ok {
				body = fixture(test, file)
			}
			got, err := testCase.parse(body, testCase.handle)
			if !errors.Is(err, testCase.wantErr) || (testCase.wantErr == nil && err != nil) || got != testCase.want {
				test.Fatalf("parse = %d, %v; want %d, %v", got, err, testCase.want, testCase.wantErr)
			}
		})
	}
}

func TestFollowerSourceURLs(test *testing.T) {
	test.Parallel()
	if got := services.FollowerSources[0].URL("jay leaton"); got != "https://api.fxtwitter.com/jay%20leaton" {
		test.Fatalf("fxtwitter URL = %s", got)
	}
}

// The video parsers read the recorded channel page and feed in order,
// dropping malformed, untitled and repeated entries.
func TestVideoParsers(test *testing.T) {
	test.Parallel()
	cases := map[string]struct {
		parse   func([]byte) ([]services.Video, error)
		body    string
		want    []string
		wantErr error
	}{
		"channel page fixture": {
			parse: services.ParseChannelPageVideos, body: "@youtube-videos.html",
			want: []string{"WCrknrNsmjQ", "Td33FCCShwg", "Xk5-4kuoIKQ", "MYTghpyGeWs", "wgUuVoZhDJQ"},
		},
		"channel page, older videoRenderer shape": {
			parse: services.ParseChannelPageVideos,
			body: `<script>var ytInitialData = {"items":[` +
				`{"videoRenderer":{"videoId":"AAAAAAAAAAA","title":{"runs":[{"text":"Runs title"}]}}},` +
				`{"videoRenderer":{"videoId":"BBBBBBBBBBB","title":{"simpleText":"Simple title"}}},` +
				`{"videoRenderer":{"videoId":"AAAAAAAAAAA","title":{"simpleText":"Repeat"}}},` +
				`{"lockupViewModel":{"contentId":"PLAYLIST000","contentType":"LOCKUP_CONTENT_TYPE_PLAYLIST"}}` +
				`]};</script>`,
			want: []string{"AAAAAAAAAAA", "BBBBBBBBBBB"},
		},
		"channel page without data": {parse: services.ParseChannelPageVideos, body: "<html>Before you continue to YouTube</html>", wantErr: services.ErrUnexpectedBody},
		"channel page broken data":  {parse: services.ParseChannelPageVideos, body: "var ytInitialData = {nope", wantErr: services.ErrUnexpectedBody},
		"channel page no videos":    {parse: services.ParseChannelPageVideos, body: `var ytInitialData = {"contents":{}};`, wantErr: services.ErrUnexpectedBody},
		"feed fixture":              {parse: services.ParseFeedVideos, body: "@youtube-feed.xml", want: []string{"WCrknrNsmjQ", "Td33FCCShwg"}},
		"feed 404 page":             {parse: services.ParseFeedVideos, body: "<!DOCTYPE html><html lang=en><title>Error 404 (Not Found)!!1</title>", wantErr: services.ErrUnexpectedBody},
	}
	for name, testCase := range cases {
		test.Run(name, func(test *testing.T) {
			test.Parallel()
			body := []byte(testCase.body)
			if file, ok := strings.CutPrefix(testCase.body, "@"); ok {
				body = fixture(test, file)
			}
			got, err := testCase.parse(body)
			if !errors.Is(err, testCase.wantErr) || (testCase.wantErr == nil && err != nil) {
				test.Fatalf("parse = %v, want %v", err, testCase.wantErr)
			}
			ids := make([]string, 0, len(got))
			for _, video := range got {
				if video.Title == "" {
					test.Errorf("video %s has no title", video.ID)
				}
				ids = append(ids, video.ID)
			}
			if strings.Join(ids, ",") != strings.Join(testCase.want, ",") {
				test.Fatalf("ids = %v, want %v", ids, testCase.want)
			}
		})
	}
	feed, _ := services.ParseFeedVideos(fixture(test, "youtube-feed.xml"))
	if feed[1].Title != "The Agent Portability Problem & How I Fixed It" {
		test.Fatalf("feed title = %q, want entities decoded", feed[1].Title)
	}
}

func TestVideoSourceURLs(test *testing.T) {
	test.Parallel()
	if got := services.VideoSources[1].URL("UC-x"); got != "https://www.youtube.com/feeds/videos.xml?channel_id=UC-x" {
		test.Fatalf("feed URL = %s", got)
	}
}
