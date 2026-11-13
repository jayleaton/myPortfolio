package ui

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/starfederation/datastar-go/datastar"

	"github.com/jayleaton/myPortfolio/apps/web/internal/apiclient"
	"github.com/jayleaton/myPortfolio/apps/web/internal/config"
	"github.com/jayleaton/myPortfolio/apps/web/ui/static"
)

// liveFollowersPoll is how often a page's follower stream re-reads the
// count. The read is in-memory (the service refreshes from X every 15
// minutes), so a short poll only bounds how late a new count shows.
const liveFollowersPoll = 20 * time.Second

// portfolioView renders one of the portfolio's pages from the page data.
type portfolioView func(page apiclient.PortfolioPage, meta ShareMeta) templ.Component

// portfolioPage is the handler of the portfolio page at path. Every page
// renders from the same in-memory data (apiclient.Portfolio), so none of
// them waits on X or YouTube.
func (ui *UI) portfolioPage(path string, view portfolioView) handlerFunc {
	return func(writer http.ResponseWriter, request *http.Request, _ apiclient.Caller) {
		page, err := ui.api.Portfolio(request.Context())
		if err != nil {
			ui.apiFailed(writer, request, err)
			return
		}
		ui.render(writer, request, http.StatusOK, view(page, ui.shareMeta(request, page, path)))
	}
}

// portfolioRoutes registers the portfolio's pages.
func (ui *UI) portfolioRoutes(page func(apiclient.Policy, string, handlerFunc)) {
	page(apiclient.Public, "GET /{$}", ui.portfolioPage(config.AppRoutes.Home, HomePage))
	page(apiclient.Public, "GET "+config.AppRoutes.LocalAI, ui.portfolioPage(config.AppRoutes.LocalAI, LocalAIPage))
	page(apiclient.Public, "GET "+config.AppRoutes.SaaS, ui.portfolioPage(config.AppRoutes.SaaS, SaaSPage))
	page(apiclient.Public, "GET "+config.AppRoutes.Socials, ui.portfolioPage(config.AppRoutes.Socials, SocialsPage))
	page(apiclient.Public, "GET "+config.AppRoutes.Archive, ui.portfolioPage(config.AppRoutes.Archive, ArchivePage))
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
	if request.URL.Query().Get(liveStreamParam) == "" {
		ui.liveHandshake(ctx, sse) // see live.go
		return
	}
	var sent apiclient.Followers
	first := true
	for {
		followers, err := ui.api.PortfolioFollowers(ctx)
		if err != nil {
			ui.log.WarnContext(ctx, "ui: live followers", "err", err)
			return // the page reconnects
		}
		if first || followersChanged(sent, followers) {
			if err := sse.MarshalAndPatchSignals(followerSignals(followers)); err != nil {
				return
			}
			sent, first = followers, false
		}
		select {
		case <-ctx.Done():
			return
		case <-ui.closing:
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

// ShareMeta is what search engines and share cards (X, Open Graph) read
// about the page at Path: its own title, description and canonical URL.
type ShareMeta struct {
	Path               string
	Title, Description string
	// ShareTitle and ShareText are the card's text.
	ShareTitle, ShareText string
	Canonical             string
	Image                 string
	XHandle               string
}

// ogImageWidth and ogImageHeight are static.OGImageFile's size, the
// summary_large_image card's 1.91:1.
const (
	ogImageWidth  = 1200
	ogImageHeight = 630
)

// shareMeta is the metadata of the page at path. It builds absolute URLs
// from PublicURL (share cards need them), or from the request's host when
// it is unset. The home page uses the site's own title and card text; any
// other page is named after its nav entry and described by its intro.
func (ui *UI) shareMeta(request *http.Request, page apiclient.PortfolioPage, path string) ShareMeta {
	origin := ui.site
	if origin == nil {
		scheme := "http"
		if request.TLS != nil {
			scheme = "https"
		}
		origin = &url.URL{Scheme: scheme, Host: request.Host}
	}
	base := strings.TrimSuffix(origin.String(), "/")
	site := page.Portfolio.Site
	meta := ShareMeta{
		Path: path, Title: site.Title, Description: site.Description, ShareTitle: site.ShareTitle, ShareText: site.ShareText,
		Canonical: base + path, Image: base + static.URL(static.OGImageFile), XHandle: page.Portfolio.XHandle,
	}
	for _, item := range navPages(page) {
		if item.Path == path {
			meta.Title = item.Page.Name + " · " + page.Portfolio.Name
			meta.Description, meta.ShareTitle, meta.ShareText = item.Page.Intro, meta.Title, item.Page.Intro
		}
	}
	return meta
}

// navItem is one page in the nav: where it lives and its copy.
type navItem struct {
	Path string
	Page apiclient.Page
}

// navPages are the pages beside home, in nav order.
func navPages(page apiclient.PortfolioPage) []navItem {
	pages := page.Portfolio.Site.Pages
	return []navItem{
		{Path: config.AppRoutes.SaaS, Page: pages.SaaS},
		{Path: config.AppRoutes.LocalAI, Page: pages.LocalAI},
		{Path: config.AppRoutes.Socials, Page: pages.Socials},
		{Path: config.AppRoutes.Archive, Page: pages.Archive},
	}
}

// The areas of featured work (content.Area*): each has its own page.
const (
	areaLocalAI = "local-ai"
	areaSaaS    = "saas"
)

// inArea are the featured projects of one area, in content order.
func inArea(featured []apiclient.Project, area string) []apiclient.Project {
	var out []apiclient.Project
	for _, project := range featured {
		if project.Area == area {
			out = append(out, project)
		}
	}
	return out
}

// projectHref is a featured project's card: the page of its area, at its slug.
func projectHref(project apiclient.Project) string {
	path := config.AppRoutes.LocalAI
	if project.Area == areaSaaS {
		path = config.AppRoutes.SaaS
	}
	return path + "#" + project.Slug
}

// shipped counts the products and sites that are live: the featured SaaS
// products and the network's entries. A network entry that is one of the
// social profiles (the YouTube channel) is a channel, not a product or
// site, so it is not counted.
func shipped(page apiclient.PortfolioPage) int {
	profiles := make(map[string]bool, len(page.Portfolio.Socials))
	for _, link := range page.Portfolio.Socials {
		profiles[link.URL] = true
	}
	count := len(inArea(page.Portfolio.Featured, areaSaaS))
	for _, link := range page.Portfolio.Network {
		if !profiles[link.URL] {
			count++
		}
	}
	return count
}

// names joins the projects' names for a one-line summary.
func names(projects []apiclient.Project) string {
	out := make([]string, 0, len(projects))
	for _, project := range projects {
		out = append(out, project.Name)
	}
	return strings.Join(out, " · ")
}

// The "Measured" view: the helpers below turn the content's projects and
// metrics into what the templates draw (decode tokens, the runs the hero
// replays, the stack figure's layers, evidence panels and the receipts
// ledger), so the templates stay dumb.

// Decode tokens: an illustrative tokeniser. The leading space rides with
// its word, a long word splits in two and trailing punctuation stands alone.
const (
	longWord   = 8  // a word of this many runes splits into two tokens
	longerWord = 10 // a word longer than this splits one rune later
	shortCut   = 4  // where a long word splits
	longCut    = 5  // where a longer word splits
)

// decodeTokens splits text into the tokens the hero lights up one by one.
func decodeTokens(text string) []string {
	var out []string
	for index, word := range strings.Fields(text) {
		prefix := ""
		if index > 0 {
			prefix = " "
		}
		core := []rune(strings.TrimRight(word, ".,"))
		tail := word[len(string(core)):]
		if len(core) >= longWord {
			cut := shortCut
			if len(core) > longerWord {
				cut = longCut
			}
			out = append(out, prefix+string(core[:cut]), string(core[cut:]))
		} else {
			out = append(out, prefix+string(core))
		}
		if tail != "" {
			out = append(out, tail)
		}
	}
	return out
}

// tokenSpans renders text as one span per decode token with nothing
// between the spans, so a split word still reads as one word.
func tokenSpans(text string) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, writer io.Writer) error {
		var out strings.Builder
		for _, token := range decodeTokens(text) {
			out.WriteString(`<span class="pf-tok">` + templ.EscapeString(token) + `</span>`)
		}
		_, err := io.WriteString(writer, out.String())
		return err
	})
}

// rateUnit marks a metric as a decode rate the hero can replay.
const rateUnit = "tok/s"

// splitMetric splits a metric's value into its number and unit:
// "87.9 tok/s" is ("87.9", "tok/s"), "~0.9 s" is ("~0.9", "s"), "50" is ("50", "").
func splitMetric(value string) (string, string) {
	number, unit, _ := strings.Cut(value, " ")
	return number, unit
}

// amount is a metric's number as a float ("~0.9" is 0.9; 0 when it is not one).
func amount(number string) float64 {
	parsed, err := strconv.ParseFloat(strings.TrimPrefix(number, "~"), 64)
	if err != nil {
		return 0
	}
	return parsed
}

// shortLabel is a metric's label up to its first comma: what was measured,
// without the conditions.
func shortLabel(label string) string {
	short, _, _ := strings.Cut(label, ",")
	return short
}

// subject names what a metric of project measured: the project's tag its
// label mentions ("DeepSeek"), or the project itself.
func subject(project apiclient.Project, metric apiclient.Metric) string {
	for _, tag := range project.Tags {
		if strings.Contains(metric.Label, tag) {
			return tag
		}
	}
	return project.Name
}

// run is one published decode rate the hero replays its text at.
type run struct {
	Name, Rate, Label, Source string
}

// decodeRuns are the featured projects' decode-rate metrics, in page order.
func decodeRuns(featured []apiclient.Project) []run {
	var out []run
	for _, project := range featured {
		out = append(out, projectRuns(project)...)
	}
	return out
}

func projectRuns(project apiclient.Project) []run {
	var out []run
	for _, metric := range project.Metrics {
		if number, unit := splitMetric(metric.Value); unit == rateUnit {
			out = append(out, run{Name: subject(project, metric), Rate: number, Label: metric.Label, Source: metric.Source})
		}
	}
	return out
}

// The stack figure's layers. A featured project sits on the layer of its
// kind (content.Visual*); anything else is a product.
const (
	layerEngine   = 1
	layerRecipes  = 2
	layerRouter   = 3
	layerProducts = 4
)

var layerNames = map[int]string{layerEngine: "Engine", layerRecipes: "Recipes", layerRouter: "Router", layerProducts: "Products"}

// siliconLayer labels the hardware plate under the project layers.
const siliconLayer = "Silicon"

func layerOf(project apiclient.Project) int {
	switch project.Visual {
	case "tensor":
		return layerEngine
	case "rack":
		return layerRecipes
	case "router":
		return layerRouter
	default:
		return layerProducts
	}
}

// levelTag is a layer's short id: "L2".
func levelTag(level int) string { return "L" + strconv.Itoa(level) }

// layerName is a project's layer name: "Recipes".
func layerName(project apiclient.Project) string { return layerNames[layerOf(project)] }

// plateText is the line under a project's name on its plate: what its
// headline number measured, or its tagline when the plate shows a decode
// run or no number.
func plateText(project apiclient.Project) string {
	if len(project.Metrics) == 0 || len(projectRuns(project)) > 0 {
		return project.Tagline
	}
	return shortLabel(project.Metrics[0].Label)
}

// stackLayer is one plate of the stack figure.
type stackLayer struct {
	Level    int
	Name     string
	Projects []apiclient.Project
}

// stackLayers groups the featured projects into plates, top layer first.
// A layer with no project is left out.
func stackLayers(featured []apiclient.Project) []stackLayer {
	var out []stackLayer
	for level := layerProducts; level >= layerEngine; level-- {
		plate := stackLayer{Level: level, Name: layerNames[level]}
		for _, project := range featured {
			if layerOf(project) == level {
				plate.Projects = append(plate.Projects, project)
			}
		}
		if len(plate.Projects) > 0 {
			out = append(out, plate)
		}
	}
	return out
}

// Evidence: what a featured card's right-hand panel shows.

// onOneAxis reports whether metrics can share one bar axis: two or more,
// all decode rates. Other numbers that happen to share a unit (two
// durations of different jobs) measure different things and stay readouts.
func onOneAxis(metrics []apiclient.Metric) bool {
	const fewest = 2
	if len(metrics) < fewest {
		return false
	}
	for _, metric := range metrics {
		if _, unit := splitMetric(metric.Value); unit != rateUnit {
			return false
		}
	}
	return true
}

// Bar axes end on a round number: the next multiple of axisStep above the
// largest value, drawn with axisTicks gaps.
const (
	axisStep  = 25
	axisTicks = 4
)

// axisMax is the end of the bar axis for metrics.
func axisMax(metrics []apiclient.Metric) int {
	top := 0.0
	for _, metric := range metrics {
		number, _ := splitMetric(metric.Value)
		top = max(top, amount(number))
	}
	return (int(top)/axisStep + 1) * axisStep
}

// axisLabels are the axis's tick labels, 0 to limit.
func axisLabels(limit int) []string {
	out := make([]string, 0, axisTicks+1)
	for tick := 0; tick <= axisTicks; tick++ {
		out = append(out, strconv.Itoa(limit*tick/axisTicks))
	}
	return out
}

// barHeight is the bar's height in its SVG's units (the width is the axis).
const barHeight = 10

// barBox is the viewBox of one bar: as wide as the axis, so a bar's width
// is its value.
func barBox(limit int) string { return fmt.Sprintf("0 0 %d %d", limit, barHeight) }

// barWidth is a metric's bar width in axis units.
func barWidth(metric apiclient.Metric) string {
	number, _ := splitMetric(metric.Value)
	return strconv.FormatFloat(amount(number), 'f', -1, 64)
}

// receipt is one row of the receipts ledger: a published number, what it
// measured and the project it belongs to.
type receipt struct {
	Project      apiclient.Project
	Metric       apiclient.Metric
	Number, Unit string
}

// receipts lists every featured metric, in page order.
func receipts(featured []apiclient.Project) []receipt {
	var out []receipt
	for _, project := range featured {
		for _, metric := range project.Metrics {
			number, unit := splitMetric(metric.Value)
			out = append(out, receipt{Project: project, Metric: metric, Number: number, Unit: unit})
		}
	}
	return out
}

// Small formatting helpers for the templates.

func itoa(value int) string { return strconv.Itoa(value) }

// boolAttr is an ARIA state's value: "true" or "false".
func boolAttr(value bool) string { return strconv.FormatBool(value) }

// groupDigits writes 2367 as "2,367" (the server-rendered count before the
// script rolls it).
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

// ogType is the page's Open Graph type: the home page is the person's
// profile, the others are pages of the site.
func ogType(meta ShareMeta) string {
	if meta.Path == config.AppRoutes.Home {
		return "profile"
	}
	return "website"
}
