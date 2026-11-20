package services_test

import (
	"encoding/xml"
	"slices"
	"strings"
	"testing"

	"github.com/jayleaton/myPortfolio/apps/web/internal/content"
	"github.com/jayleaton/myPortfolio/apps/web/internal/services"
)

func portfolioPageRoutes(portfolio services.Portfolio) []services.PageRoute {
	pages := portfolio.Site.Pages
	return []services.PageRoute{
		{Path: "/saas", Page: pages.SaaS},
		{Path: "/local-ai", Page: pages.LocalAI},
		{Path: "/socials", Page: pages.Socials},
		{Path: "/archive", Page: pages.Archive},
	}
}

func TestBuildShareMetaEveryPageUnique(test *testing.T) {
	test.Parallel()
	portfolio := content.Source{}.Portfolio()
	origin := "https://www.jayleaton.com/"
	imagePath := "/static/og.png?v=deadbeef"
	subpages := portfolioPageRoutes(portfolio)
	allPaths := []string{"/", "/saas", "/local-ai", "/socials", "/archive"}

	titles := make(map[string]string, len(allPaths))
	descriptions := make(map[string]string, len(allPaths))
	for _, routePath := range allPaths {
		meta := services.BuildShareMeta(portfolio, origin, routePath, imagePath, subpages)
		if meta.Title == "" || meta.Description == "" || meta.ShareTitle == "" || meta.ShareText == "" {
			test.Fatalf("path %s has empty metadata: %+v", routePath, meta)
		}
		if want := "https://www.jayleaton.com" + routePath; meta.Canonical != want {
			test.Errorf("path %s canonical = %q, want %q", routePath, meta.Canonical, want)
		}
		if want := "https://www.jayleaton.com" + imagePath; meta.Image != want {
			test.Errorf("path %s image = %q, want %q", routePath, meta.Image, want)
		}
		wantType := "website"
		if routePath == "/" {
			wantType = "profile"
		}
		if meta.OGType != wantType {
			test.Errorf("path %s og:type = %q, want %q", routePath, meta.OGType, wantType)
		}
		if other, duplicate := titles[meta.Title]; duplicate {
			test.Errorf("duplicate title %q on %s and %s", meta.Title, other, routePath)
		}
		titles[meta.Title] = routePath
		if other, duplicate := descriptions[meta.Description]; duplicate {
			test.Errorf("duplicate description %q on %s and %s", meta.Description, other, routePath)
		}
		descriptions[meta.Description] = routePath
	}
}

func TestPersonLDIncludesSameAsWithoutDuplicates(test *testing.T) {
	test.Parallel()
	portfolio := content.Source{}.Portfolio()
	var buildThings services.Link
	for _, item := range portfolio.Network {
		if item.Name == "BuildThings" {
			buildThings = item
			break
		}
	}
	extra := []services.Link{
		{Name: "GitHub", URL: "https://github.com/jayleaton"},
		buildThings,
		{Name: "Empty", URL: "   "},
	}
	meta := services.BuildShareMeta(portfolio, "https://www.jayleaton.com", "/", "/static/og.png?v=1", portfolioPageRoutes(portfolio))
	doc := services.PersonLD(portfolio, meta, extra)
	if doc["@context"] != "https://schema.org" || doc["@type"] != "Person" || doc["name"] != portfolio.Name {
		test.Fatalf("PersonLD = %#v", doc)
	}
	if doc["url"] != "https://www.jayleaton.com/" || doc["image"] != "https://www.jayleaton.com/static/og.png?v=1" {
		test.Fatalf("PersonLD urls = %#v", doc)
	}
	sameAs, ok := doc["sameAs"].([]string)
	if !ok {
		test.Fatalf("sameAs type = %T", doc["sameAs"])
	}
	required := []string{
		"https://x.com/jayleaton",
		"https://github.com/jayleaton",
		"https://youtube.com/@build-things",
		"https://buildthings.co",
		"https://ko-fi.com/jayleaton",
	}
	for _, want := range required {
		if !slices.Contains(sameAs, want) {
			test.Errorf("sameAs %v missing %q", sameAs, want)
		}
	}
	seen := make(map[string]bool, len(sameAs))
	for _, entry := range sameAs {
		if seen[entry] {
			test.Errorf("duplicate sameAs entry %q", entry)
		}
		seen[entry] = true
	}
}

func TestRenderSitemapAndRobots(test *testing.T) {
	test.Parallel()
	paths := []string{"/", "/saas", "/local-ai", "/socials", "/archive"}
	raw, err := services.RenderSitemap("https://www.jayleaton.com/", paths)
	if err != nil {
		test.Fatalf("RenderSitemap: %v", err)
	}
	if !strings.HasPrefix(string(raw), xml.Header) {
		test.Fatalf("missing XML header: %s", raw)
	}
	var parsed struct {
		XMLName xml.Name `xml:"urlset"`
		Xmlns   string   `xml:"xmlns,attr"`
		URLs    []struct {
			Loc string `xml:"loc"`
		} `xml:"url"`
	}
	if err := xml.Unmarshal(raw, &parsed); err != nil {
		test.Fatalf("xml.Unmarshal: %v\n%s", err, raw)
	}
	if parsed.Xmlns != "http://www.sitemaps.org/schemas/sitemap/0.9" {
		test.Errorf("xmlns = %q", parsed.Xmlns)
	}
	if len(parsed.URLs) != len(paths) {
		test.Fatalf("url count = %d, want %d", len(parsed.URLs), len(paths))
	}
	for index, routePath := range paths {
		want := "https://www.jayleaton.com" + routePath
		if parsed.URLs[index].Loc != want {
			test.Errorf("url[%d] = %q, want %q", index, parsed.URLs[index].Loc, want)
		}
	}

	robots := services.RenderRobots("https://www.jayleaton.com/", "/sitemap.xml")
	wantRobots := "User-agent: *\nAllow: /\nSitemap: https://www.jayleaton.com/sitemap.xml\n"
	if robots != wantRobots {
		test.Errorf("RenderRobots = %q, want %q", robots, wantRobots)
	}
	if bare := services.RenderRobots("", "/sitemap.xml"); bare != "User-agent: *\nAllow: /\n" {
		test.Errorf("RenderRobots(empty) = %q", bare)
	}
}
