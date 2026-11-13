// Package content is the portfolio's copy and links as typed Go values: the
// featured work, the network of sites, socials, the archive and the
// channels whose live numbers the site shows. Edit data.go to change the
// site; nothing else hard-codes copy. Services read it through the
// services.ContentSource port (Source implements it).
package content

// Portfolio is everything the home page shows that is not fetched live.
type Portfolio struct {
	Name     string
	Headline string
	Bio      string
	// XHandle is the X account whose follower count the site shows live.
	XHandle string
	// YouTube is the channel whose latest videos the site shows.
	YouTube  Channel
	Featured []Project
	Network  []Link
	Socials  []Link
	Archive  []Project
	// FallbackVideos show until the channel has been fetched once.
	FallbackVideos []Video
	// Site is the page's own copy: title, share text and section labels.
	Site Site
}

// Site is the copy around the content: metadata for search and share
// cards, the pages and the labels of their sections.
type Site struct {
	// Title and Description are the home page's; the other pages use their
	// Page's Name and Intro.
	Title       string
	Description string
	// ShareTitle and ShareText are the home page's X / Open Graph card.
	ShareTitle, ShareText string
	// Kicker sits above the home page's headline.
	Kicker string
	// Hardware is the bottom plate of the stack figure.
	Hardware Hardware
	Pages    Pages
	Sections Sections
	Footer   string
}

// Hardware is the machine the published numbers were measured on: its
// name, the memory and chip of one unit, the units and the link between them.
type Hardware struct {
	Name, Memory, Chip, Link string
	Units                    []string
}

// Pages are the pages the home page opens onto, in nav order.
type Pages struct {
	LocalAI, SaaS, Socials, Archive Page
}

// Page is one page's copy: its number and name in the nav, the headline
// and intro at its top (the intro is also its search description) and the
// line on its home-page tile.
type Page struct {
	Number, Name, Heading, Intro, Teaser string
}

// Sections are the labels of the sections inside the pages: Work and
// Receipts on the local AI page, Products and Network on the SaaS page,
// Uplinks and Videos on the socials page.
type Sections struct {
	Work, Receipts, Products, Network, Uplinks, Videos Section
}

// Section is one section's anchor, number, title and one-line intro.
type Section struct {
	ID, Number, Title, Intro string
}

// Project is one piece of work. Repo and Image may be empty; Image is a
// file name under the UI's static assets (img/<name>.webp).
type Project struct {
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Tagline     string `json:"tagline"`
	Description string `json:"description"`
	Role        string `json:"role"`
	// Area is the page a featured project belongs to (AreaLocalAI or AreaSaaS).
	Area  string `json:"area,omitempty"`
	URL   string `json:"url"`
	Repo  string `json:"repo,omitempty"`
	Image string `json:"image,omitempty"`
	// Visual is the kind of work (one of the Visual* constants). It places
	// a featured project on a layer of the hero's stack figure.
	Visual     string   `json:"visual,omitempty"`
	Tags       []string `json:"tags"`
	Highlights []string `json:"highlights,omitempty"`
	Metrics    []Metric `json:"metrics,omitempty"`
}

// The areas of featured work: each has its own page.
const (
	AreaLocalAI = "local-ai" // engines, recipes and routing
	AreaSaaS    = "saas"     // products
)

// The kinds of featured work.
const (
	VisualTensor = "tensor" // an inference engine: the stack's engine layer
	VisualRack   = "rack"   // GPU recipes: the recipes layer
	VisualRouter = "router" // model routing: the router layer
	VisualRadar  = "radar"  // a product (security)
	VisualFrame  = "frame"  // a product (imaging)
)

// Metric is a published, sourced number (never a live one).
type Metric struct {
	Value  string `json:"value"`
	Label  string `json:"label"`
	Source string `json:"source"`
}

// Link is a site in the network or a social profile.
type Link struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Image       string `json:"image,omitempty"`
}

// Channel is a YouTube channel: its ID feeds the latest-videos fetch.
type Channel struct {
	ID     string `json:"id"`
	Handle string `json:"handle"`
	URL    string `json:"url"`
}

// Video is one YouTube upload.
type Video struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// Source serves the portfolio (services.ContentSource).
type Source struct{}

// Portfolio returns the site's content. Callers get their own copy of the
// top-level value; treat the slices as read-only.
func (Source) Portfolio() Portfolio { return portfolio }
