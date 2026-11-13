package services

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"maps"
	"net/url"
	"regexp"
	"slices"
	"strings"
)

// VideoSource is one free place to read a YouTube channel's latest uploads.
type VideoSource struct {
	Name  string
	URL   func(channelID string) string
	Parse func(body []byte) ([]Video, error)
}

// VideoSources are tried in order. The channel's videos page carries the
// latest uploads in its ytInitialData; the RSS feed is the documented
// source but has answered 404 for channels since late 2026, so it is the
// fallback.
var VideoSources = []VideoSource{
	{Name: "channel page", URL: func(id string) string { return "https://www.youtube.com/channel/" + url.PathEscape(id) + "/videos" }, Parse: ParseChannelPageVideos},
	{Name: "rss", URL: func(id string) string {
		return "https://www.youtube.com/feeds/videos.xml?" + url.Values{"channel_id": {id}}.Encode()
	}, Parse: ParseFeedVideos},
}

// VideoLimit is how many latest videos the site shows.
const VideoLimit = 4

// videoID is the shape of a YouTube video id; anything else is dropped, so
// a parsed id is always safe in a watch or thumbnail URL.
var videoID = regexp.MustCompile(`^[A-Za-z0-9_-]{11}$`)

// initialDataMarker starts the JSON YouTube embeds in its pages.
var initialDataMarker = []byte("var ytInitialData = ")

// ParseChannelPageVideos reads the uploads from a channel's /videos page
// (testdata/youtube-videos.html): every lockupViewModel of a video, and the
// older videoRenderer shape, in page order.
func ParseChannelPageVideos(body []byte) ([]Video, error) {
	start := bytes.Index(body, initialDataMarker)
	if start < 0 {
		return nil, fmt.Errorf("%w: no ytInitialData", ErrUnexpectedBody)
	}
	var data any
	// Decode reads one JSON value and ignores the script after it.
	if err := json.NewDecoder(bytes.NewReader(body[start+len(initialDataMarker):])).Decode(&data); err != nil {
		return nil, fmt.Errorf("%w: ytInitialData: %w", ErrUnexpectedBody, err)
	}
	var videos []Video
	walkJSON(data, func(object map[string]any) {
		if lockup, ok := object["lockupViewModel"].(map[string]any); ok {
			if lockup["contentType"] == "LOCKUP_CONTENT_TYPE_VIDEO" {
				id, _ := lockup["contentId"].(string)
				videos = append(videos, Video{ID: id, Title: stringAt(lockup, "metadata", "lockupMetadataViewModel", "title", "content")})
			}
		}
		if renderer, ok := object["videoRenderer"].(map[string]any); ok {
			id, _ := renderer["videoId"].(string)
			title := stringAt(renderer, "title", "simpleText")
			if runs, ok := mapAt(renderer, "title")["runs"].([]any); ok && title == "" && len(runs) > 0 {
				title, _ = runs[0].(map[string]any)["text"].(string)
			}
			videos = append(videos, Video{ID: id, Title: title})
		}
	})
	return validVideos(videos)
}

// ParseFeedVideos reads the entries of a channel's Atom feed
// (https://www.youtube.com/feeds/videos.xml?channel_id=..., testdata/youtube-feed.xml).
func ParseFeedVideos(body []byte) ([]Video, error) {
	var feed struct {
		Entries []struct {
			VideoID string `xml:"http://www.youtube.com/xml/schemas/2015 videoId"`
			Title   string `xml:"http://www.w3.org/2005/Atom title"`
		} `xml:"http://www.w3.org/2005/Atom entry"`
	}
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUnexpectedBody, err)
	}
	videos := make([]Video, 0, len(feed.Entries))
	for _, entry := range feed.Entries {
		videos = append(videos, Video{ID: entry.VideoID, Title: entry.Title})
	}
	return validVideos(videos)
}

// validVideos keeps well-formed, titled videos once each, in order, and
// fails when none are left.
func validVideos(videos []Video) ([]Video, error) {
	seen := map[string]bool{}
	out := make([]Video, 0, len(videos))
	for _, video := range videos {
		video.Title = strings.TrimSpace(video.Title)
		if !videoID.MatchString(video.ID) || video.Title == "" || seen[video.ID] {
			continue
		}
		seen[video.ID] = true
		out = append(out, video)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("%w: no videos", ErrUnexpectedBody)
	}
	return out, nil
}

// limitVideos copies at most VideoLimit videos.
func limitVideos(videos []Video) []Video {
	return append([]Video(nil), videos[:min(len(videos), VideoLimit)]...)
}

// walkJSON calls visit on every object in value, depth first, in order.
func walkJSON(value any, visit func(map[string]any)) {
	switch node := value.(type) {
	case map[string]any:
		visit(node)
		// Sorted keys: a map's order is random, the page's order is not.
		for _, key := range slices.Sorted(maps.Keys(node)) {
			walkJSON(node[key], visit)
		}
	case []any:
		for _, child := range node {
			walkJSON(child, visit)
		}
	}
}

func mapAt(object map[string]any, keys ...string) map[string]any {
	for _, key := range keys {
		object, _ = object[key].(map[string]any)
	}
	return object
}

func stringAt(object map[string]any, keys ...string) string {
	last := len(keys) - 1
	text, _ := mapAt(object, keys[:last]...)[keys[last]].(string)
	return text
}
