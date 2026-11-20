package services

import (
	"encoding/xml"
	"strings"
)

// ShareMeta holds the search and share-card metadata for one public page:
// its own title, description, canonical URL and Open Graph card fields.
type ShareMeta struct {
	Path               string
	Title, Description string
	// ShareTitle and ShareText are the card's text.
	ShareTitle, ShareText string
	Canonical             string
	Image                 string
	OGType                string
	XHandle               string
}

// PageRoute pairs a public route path with its page copy.
type PageRoute struct {
	Path string
	Page Page
}

// BuildShareMeta resolves the search and share metadata for path under origin.
// The home page ("/") uses the site's own title, description and profile type;
// every listed subpage uses its own name, intro and website type.
func BuildShareMeta(portfolio Portfolio, origin, path, imagePath string, pages []PageRoute) ShareMeta {
	base := strings.TrimSuffix(origin, "/")
	site := portfolio.Site
	meta := ShareMeta{
		Path:        path,
		Title:       site.Title,
		Description: site.Description,
		ShareTitle:  site.ShareTitle,
		ShareText:   site.ShareText,
		Canonical:   base + path,
		Image:       base + imagePath,
		OGType:      "profile",
		XHandle:     portfolio.XHandle,
	}
	for _, item := range pages {
		if item.Path == path {
			meta.Title = item.Page.Name + " · " + portfolio.Name
			meta.Description = item.Page.Intro
			meta.ShareTitle = meta.Title
			meta.ShareText = item.Page.Intro
			meta.OGType = "website"
			break
		}
	}
	return meta
}

// PersonSameAs collects the distinct profile and company URLs from socials
// followed by extra links, preserving order.
func PersonSameAs(socials, extra []Link) []string {
	seen := make(map[string]bool, len(socials)+len(extra))
	out := make([]string, 0, len(socials)+len(extra))
	appendLinks := func(links []Link) {
		for _, link := range links {
			trimmed := strings.TrimSpace(link.URL)
			if trimmed == "" || seen[trimmed] {
				continue
			}
			seen[trimmed] = true
			out = append(out, trimmed)
		}
	}
	appendLinks(socials)
	appendLinks(extra)
	return out
}

// PersonLD builds the schema.org Person JSON-LD document for the portfolio.
func PersonLD(portfolio Portfolio, meta ShareMeta, extra []Link) map[string]any {
	return map[string]any{
		"@context":    "https://schema.org",
		"@type":       "Person",
		"name":        portfolio.Name,
		"url":         meta.Canonical,
		"image":       meta.Image,
		"description": portfolio.Bio,
		"sameAs":      PersonSameAs(portfolio.Socials, extra),
	}
}

const sitemapNamespace = "http://www.sitemaps.org/schemas/sitemap/0.9"

type sitemapURLSet struct {
	XMLName xml.Name     `xml:"urlset"`
	Xmlns   string       `xml:"xmlns,attr"`
	URLs    []sitemapURL `xml:"url"`
}

type sitemapURL struct {
	Loc string `xml:"loc"`
}

// SitemapURLs joins each public route path onto origin.
func SitemapURLs(origin string, paths []string) []string {
	base := strings.TrimSuffix(origin, "/")
	out := make([]string, 0, len(paths))
	for _, routePath := range paths {
		out = append(out, base+routePath)
	}
	return out
}

// RenderSitemap encodes a standard XML sitemap for paths under origin.
func RenderSitemap(origin string, paths []string) ([]byte, error) {
	locs := SitemapURLs(origin, paths)
	entries := make([]sitemapURL, 0, len(locs))
	for _, loc := range locs {
		entries = append(entries, sitemapURL{Loc: loc})
	}
	body, err := xml.MarshalIndent(sitemapURLSet{
		Xmlns: sitemapNamespace,
		URLs:  entries,
	}, "", "  ")
	if err != nil {
		return nil, err
	}
	return append([]byte(xml.Header), append(body, '\n')...), nil
}

// RenderRobots returns robots.txt allowing crawling and advertising the sitemap.
func RenderRobots(origin, sitemapPath string) string {
	base := strings.TrimSuffix(origin, "/")
	if base == "" || sitemapPath == "" {
		return "User-agent: *\nAllow: /\n"
	}
	return "User-agent: *\nAllow: /\nSitemap: " + base + sitemapPath + "\n"
}
