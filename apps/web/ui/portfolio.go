package ui

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/starfederation/datastar-go/datastar"

	"github.com/jayleaton/myPortfolio/apps/web/internal/apiclient"
	"github.com/jayleaton/myPortfolio/apps/web/ui/static"
)

// liveFollowersPoll is how often a page's follower stream re-reads the
// count. The read is in-memory (the service refreshes from X every 15
// minutes), so a short poll only bounds how late a new count shows.
const liveFollowersPoll = 20 * time.Second

func (ui *UI) home(writer http.ResponseWriter, request *http.Request, _ apiclient.Caller) {
	page, err := ui.api.Portfolio(request.Context())
	if err != nil {
		ui.apiFailed(writer, request, err)
		return
	}
	ui.render(writer, request, http.StatusOK, PortfolioPage(page, ui.shareMeta(request, page)))
}

// followerSignals are the Datastar signals the follower count binds to.
func followerSignals(followers apiclient.Followers) map[string]any {
	return map[string]any{"followers": followers.Count, "followersKnown": followers.Known}
}

// liveFollowers streams the X follower count: the current value at once,
// then a patch whenever the service's last good value changes. The page
// animates the number to each new value.
func (ui *UI) liveFollowers(writer http.ResponseWriter, request *http.Request, _ apiclient.Caller) {
	ctx := request.Context()
	sse := datastar.NewSSE(writer, request)
	tick := time.NewTicker(liveFollowersPoll)
	defer tick.Stop()
	var sent apiclient.Followers
	first := true
	for {
		followers, err := ui.api.PortfolioFollowers(ctx)
		if err != nil {
			ui.log.WarnContext(ctx, "ui: live followers", "err", err)
			return
		}
		if first || followers.Count != sent.Count || followers.Known != sent.Known {
			if err := sse.MarshalAndPatchSignals(followerSignals(followers)); err != nil {
				return
			}
			sent, first = followers, false
		}
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
		}
	}
}

// robots allows everything and points crawlers at the canonical origin.
func (ui *UI) robots(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.Header().Set("Cache-Control", "public, max-age=86400")
	_, _ = fmt.Fprint(writer, "User-agent: *\nAllow: /\n")
}

// ShareMeta is what search engines and share cards (X, Open Graph) read.
type ShareMeta struct {
	Canonical string
	Image     string
	XHandle   string
}

// ogImageWidth and ogImageHeight are static.OGImageFile's size, the
// summary_large_image card's 1.91:1.
const (
	ogImageWidth  = 1200
	ogImageHeight = 630
)

// shareMeta builds absolute URLs from PublicURL (share cards need them),
// or from the request's host when it is unset.
func (ui *UI) shareMeta(request *http.Request, page apiclient.PortfolioPage) ShareMeta {
	origin := ui.site
	if origin == nil {
		scheme := "http"
		if request.TLS != nil {
			scheme = "https"
		}
		origin = &url.URL{Scheme: scheme, Host: request.Host}
	}
	base := strings.TrimSuffix(origin.String(), "/")
	return ShareMeta{Canonical: base + "/", Image: base + static.URL(static.OGImageFile), XHandle: page.Portfolio.XHandle}
}

// The architecture diagram: the core (JJ) in the middle, the featured
// projects around it, and the two live feeds (X, YouTube) above and below.
// Coordinates are computed here so templates stay dumb; the CSS animates
// the packets along the edges (paused for prefers-reduced-motion).

// diagramNode is one box of the diagram.
type diagramNode struct {
	Href, Label, Sub, Kind   string
	Left, Top, Width, Height int
}

// diagramEdge is a curve from the core to a node, and its packets' timing.
type diagramEdge struct {
	ID, Path      string
	Duration      string
	Begin         [2]string
	Reverse       bool
	Live          bool
	LabelX, Label string
}

// diagram is one layout (wide for desktop, tall for phones).
type diagram struct {
	Width, Height int
	Core          diagramNode
	Nodes         []diagramNode
	Feeds         []diagramNode
	Edges         []diagramEdge
}

func (graph diagram) ViewBox() string { return fmt.Sprintf("0 0 %d %d", graph.Width, graph.Height) }

// Layout constants of the two diagrams, in viewBox units.
const (
	wideWidth, wideHeight = 1000, 600
	wideNodeW, wideNodeH  = 200, 72
	wideCoreW, wideCoreH  = 220, 96
	wideMargin            = 24
	wideFeedW, wideFeedH  = 300, 48
	tallWidth, tallHeight = 360, 700
	tallNodeW, tallNodeH  = 156, 64
	tallCoreW, tallCoreH  = 200, 80
	tallGap               = 16
	tallFeedW, tallFeedH  = 300, 44
	tallTop               = 12
	tallCoreTop           = 96
	tallRowsTop           = 236
	tallRowStep           = 92
	half                  = 2
	packetSeconds         = 2.6
	packetStep            = 0.35
	packetSpeeds          = 3 // edges cycle through three packet speeds
)

// nodeSub is the second line of a node: its first published metric, or its role.
func nodeSub(project apiclient.Project) string {
	if len(project.Metrics) > 0 {
		return project.Metrics[0].Value
	}
	return strings.ToLower(project.Role)
}

func packetTiming(index int) (string, [2]string) {
	duration := packetSeconds + float64(index%packetSpeeds)*packetStep
	format := func(seconds float64) string { return strconv.FormatFloat(seconds, 'f', 2, 64) + "s" }
	return format(duration), [2]string{format(float64(index) * packetStep), format(float64(index)*packetStep + duration/half)}
}

func curve(fromX, fromY, toX, toY int, vertical bool) string {
	if vertical {
		mid := (fromY + toY) / half
		return fmt.Sprintf("M%d %d C%d %d %d %d %d %d", fromX, fromY, fromX, mid, toX, mid, toX, toY)
	}
	mid := (fromX + toX) / half
	return fmt.Sprintf("M%d %d C%d %d %d %d %d %d", fromX, fromY, mid, fromY, mid, toY, toX, toY)
}

// feedNodes are the X and YouTube feeds' boxes.
func feedNodes(page apiclient.PortfolioPage) (diagramNode, diagramNode) {
	return diagramNode{Href: "#uplinks", Label: "x.com/@" + page.Portfolio.XHandle, Sub: "followers", Kind: "feed"},
		diagramNode{Href: "#videos", Label: "youtube/@" + page.Portfolio.YouTube.Handle, Sub: "latest uploads", Kind: "feed"}
}

// wideDiagram puts the featured projects in two columns either side of the core.
func wideDiagram(page apiclient.PortfolioPage) diagram {
	graph := diagram{Width: wideWidth, Height: wideHeight}
	graph.Core = diagramNode{Href: "#work", Label: "jj.eaton", Sub: "core · router", Kind: "core",
		Left: (wideWidth - wideCoreW) / half, Top: (wideHeight - wideCoreH) / half, Width: wideCoreW, Height: wideCoreH}
	featured := page.Portfolio.Featured
	left := (len(featured) + 1) / half
	columns := [][]apiclient.Project{featured[:left], featured[left:]}
	index := 0
	for side, column := range columns {
		left := wideMargin
		if side == 1 {
			left = wideWidth - wideMargin - wideNodeW
		}
		step := wideHeight / (len(column) + 1)
		for row, project := range column {
			node := diagramNode{Href: "#" + project.Slug, Label: project.Name, Sub: nodeSub(project), Kind: project.Visual,
				Left: left, Top: step*(row+1) - wideNodeH/half, Width: wideNodeW, Height: wideNodeH}
			graph.Nodes = append(graph.Nodes, node)
			fromX, toX := graph.Core.Left, node.Left+node.Width
			if side == 1 {
				fromX, toX = graph.Core.Left+graph.Core.Width, node.Left
			}
			duration, begin := packetTiming(index)
			graph.Edges = append(graph.Edges, diagramEdge{ID: fmt.Sprintf("edge-w-%d", index), Path: curve(fromX, graph.Core.Top+graph.Core.Height/half, toX, node.Top+node.Height/half, false),
				Duration: duration, Begin: begin, Reverse: index%2 == 1})
			index++
		}
	}
	xFeed, youtubeFeed := feedNodes(page)
	xFeed.Left, xFeed.Top, xFeed.Width, xFeed.Height = (wideWidth-wideFeedW)/half, wideMargin, wideFeedW, wideFeedH
	youtubeFeed.Left, youtubeFeed.Top, youtubeFeed.Width, youtubeFeed.Height = (wideWidth-wideFeedW)/half, wideHeight-wideMargin-wideFeedH, wideFeedW, wideFeedH
	graph.Feeds = []diagramNode{xFeed, youtubeFeed}
	center := wideWidth / half
	graph.Edges = append(graph.Edges,
		diagramEdge{ID: "edge-w-x", Path: curve(center, xFeed.Top+xFeed.Height, center, graph.Core.Top, true), Duration: "1.8s", Begin: [2]string{"0s", "0.9s"}, Live: true},
		diagramEdge{ID: "edge-w-yt", Path: curve(center, youtubeFeed.Top, center, graph.Core.Top+graph.Core.Height, true), Duration: "2.2s", Begin: [2]string{"0.3s", "1.4s"}, Live: true},
	)
	return graph
}

// tallDiagram stacks the feed, the core and the projects in rows of two.
func tallDiagram(page apiclient.PortfolioPage) diagram {
	featured := page.Portfolio.Featured
	rows := (len(featured) + 1) / half
	graph := diagram{Width: tallWidth, Height: tallRowsTop + rows*tallRowStep + tallFeedH + tallGap}
	if graph.Height < tallHeight/half {
		graph.Height = tallHeight / half
	}
	xFeed, youtubeFeed := feedNodes(page)
	xFeed.Left, xFeed.Top, xFeed.Width, xFeed.Height = (tallWidth-tallFeedW)/half, tallTop, tallFeedW, tallFeedH
	graph.Core = diagramNode{Href: "#work", Label: "jj.eaton", Sub: "core · router", Kind: "core",
		Left: (tallWidth - tallCoreW) / half, Top: tallCoreTop, Width: tallCoreW, Height: tallCoreH}
	center := tallWidth / half
	graph.Edges = append(graph.Edges, diagramEdge{ID: "edge-t-x", Path: curve(center, xFeed.Top+xFeed.Height, center, graph.Core.Top, true), Duration: "1.6s", Begin: [2]string{"0s", "0.8s"}, Live: true})
	for index, project := range featured {
		row, column := index/half, index%half
		left := tallGap + column*(tallNodeW+tallGap)
		if index == len(featured)-1 && column == 0 {
			left = (tallWidth - tallNodeW) / half // a lone last node is centred
		}
		node := diagramNode{Href: "#" + project.Slug, Label: project.Name, Sub: nodeSub(project), Kind: project.Visual,
			Left: left, Top: tallRowsTop + row*tallRowStep, Width: tallNodeW, Height: tallNodeH}
		graph.Nodes = append(graph.Nodes, node)
		duration, begin := packetTiming(index)
		graph.Edges = append(graph.Edges, diagramEdge{ID: fmt.Sprintf("edge-t-%d", index), Path: curve(center, graph.Core.Top+graph.Core.Height, node.Left+node.Width/half, node.Top, true),
			Duration: duration, Begin: begin, Reverse: index%2 == 1})
	}
	youtubeFeed.Left, youtubeFeed.Top, youtubeFeed.Width, youtubeFeed.Height = (tallWidth-tallFeedW)/half, graph.Height-tallFeedH-tallTop, tallFeedW, tallFeedH
	graph.Feeds = []diagramNode{xFeed, youtubeFeed}
	return graph
}

// Small formatting helpers for the templates.

func itoa(value int) string { return strconv.Itoa(value) }

// groupDigits writes 2367 as "2,367" (the server-rendered count before the
// script animates it).
func groupDigits(value int) string {
	digits := strconv.Itoa(value)
	const group = 3
	var out strings.Builder
	for index, digit := range digits {
		if index > 0 && (len(digits)-index)%group == 0 {
			out.WriteByte(',')
		}
		out.WriteRune(digit)
	}
	return out.String()
}

// host is a URL's host without www., for the mono labels.
func host(link string) string {
	parsed, err := url.Parse(link)
	if err != nil || parsed.Host == "" {
		return link
	}
	return strings.TrimPrefix(parsed.Host, "www.") + strings.TrimSuffix(parsed.Path, "/")
}

// pad2 is a two-digit index (01, 02...).
func pad2(index int) string { return fmt.Sprintf("%02d", index+1) }

func thumbnail(video apiclient.Video) string {
	return "https://i.ytimg.com/vi/" + url.PathEscape(video.ID) + "/hqdefault.jpg"
}

func watchURL(video apiclient.Video) string {
	return "https://www.youtube.com/watch?" + url.Values{"v": {video.ID}}.Encode()
}

// socialFor finds a social link by name ("" when absent).
func socialFor(links []apiclient.Link, name string) apiclient.Link {
	for _, link := range links {
		if link.Name == name {
			return link
		}
	}
	return apiclient.Link{}
}

// otherSocials are the socials besides X (which gets its own live card).
func otherSocials(links []apiclient.Link) []apiclient.Link {
	out := make([]apiclient.Link, 0, len(links))
	for _, link := range links {
		if link.Name != "X" {
			out = append(out, link)
		}
	}
	return out
}

// telemetry is one decorative readout in the hero: a published number the
// script jitters around (data-base, data-jitter), never presented as live.
type telemetry struct {
	Label, Value, Unit, Base, Jitter, Decimals string
}

var heroTelemetry = []telemetry{
	{Label: "decode", Value: "79.9", Unit: "tok/s", Base: "79.9", Jitter: "1.6", Decimals: "1"},
	{Label: "voice turn", Value: "0.92", Unit: "s p50", Base: "0.92", Jitter: "0.03", Decimals: "2"},
	{Label: "image 1024²", Value: "6.3", Unit: "s warm", Base: "6.3", Jitter: "0.2", Decimals: "1"},
	{Label: "glm-5.3", Value: "72.4", Unit: "tok/s", Base: "72.4", Jitter: "1.2", Decimals: "1"},
}

// routerCenterX and routerCenterY are the router visual's hub.
const routerCenterX, routerCenterY = 160, 90

// routerPath is the router visual's curve from its hub to a model at x, y.
func routerPath(targetX, targetY int) string {
	return curve(routerCenterX, routerCenterY, targetX, targetY, false)
}

// personLD is the page's schema.org Person (search results, link previews).
func personLD(page apiclient.PortfolioPage, meta ShareMeta) map[string]any {
	sameAs := make([]string, 0, len(page.Portfolio.Socials))
	for _, link := range page.Portfolio.Socials {
		sameAs = append(sameAs, link.URL)
	}
	return map[string]any{
		"@context":    "https://schema.org",
		"@type":       "Person",
		"name":        page.Portfolio.Name,
		"url":         meta.Canonical,
		"image":       meta.Image,
		"description": page.Portfolio.Bio,
		"sameAs":      sameAs,
	}
}

// routerTargets are the router visual's model nodes.
var routerTargets = [][2]int{{40, 34}, {40, 90}, {40, 146}, {280, 34}, {280, 90}, {280, 146}}
