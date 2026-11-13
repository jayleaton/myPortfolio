// Package static embeds the UI's browser assets: the generated Tailwind CSS,
// the pinned Datastar bundle (no CDN), the portfolio's script, fonts and
// images. It serves them versioned, and the text ones gzipped.
package static

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"io/fs"
	"mime"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/jayleaton/myPortfolio/apps/web/ui/routes"
)

// Datastar bundle pin. Update all three together; TestDatastarPin checks the
// embedded file against DatastarSHA256.
// Source: https://raw.githubusercontent.com/starfederation/datastar/v1.0.4/bundles/datastar.js
const (
	DatastarVersion = "1.0.4"
	DatastarFile    = "datastar-1.0.4.js"
	DatastarSHA256  = "727844adfc825ee651fb93c544a2a739986f9a21820a94524b35f0cac470cf91"
)

// FaviconFile is the original template's favicon, also served at /favicon.ico.
const FaviconFile = "favicon.ico"

// CSSFile is written by the Tailwind standalone CLI (see apps/web/ui/README.md).
const CSSFile = "app.css"

// PortfolioScript is the portfolio page's only script of its own: the
// hero's decode replay, the follower odometer, the menu and the reveals.
const PortfolioScript = "portfolio.js"

// OGImageFile is the 1200×630 share card (X, Open Graph).
const OGImageFile = "og.png"

// The self-hosted fonts (SIL Open Font License; the license texts sit
// beside them in fonts/). DisplayFontFile is preloaded by the page:
// portfolio.css names the same URL, version included, so keep its ?v= in
// step with URL(DisplayFontFile) when the file changes.
const (
	DisplayFontFile = "fonts/archivo.woff2"
	MonoFontFile    = "fonts/jetbrains-mono.woff2"
)

// Image is the URL of an image under img/ (content.Project.Image).
func Image(name string) string { return URL(name) }

//go:embed app.css datastar-1.0.4.js favicon.ico portfolio.js og.png fonts img
var FS embed.FS

var versions = func() map[string]string {
	out := map[string]string{}
	_ = fs.WalkDir(FS, ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}
		file, err := FS.ReadFile(path)
		if err != nil {
			return err
		}
		sum := sha256.Sum256(file)
		out[path] = hex.EncodeToString(sum[:8])
		return nil
	})
	return out
}()

// gzipped holds the text assets (CSS, JavaScript) compressed once at start:
// they are most of the page's bytes and shrink to about a fifth. Fonts and
// images are already compressed and are served as they are.
var gzipped = func() map[string][]byte {
	out := map[string][]byte{}
	for name := range versions {
		if extension := path.Ext(name); extension != ".css" && extension != ".js" {
			continue
		}
		file, err := FS.ReadFile(name)
		if err != nil {
			continue
		}
		var packed bytes.Buffer
		writer, err := gzip.NewWriterLevel(&packed, gzip.BestCompression)
		if err != nil {
			continue
		}
		if _, err := writer.Write(file); err != nil || writer.Close() != nil {
			continue
		}
		out[name] = packed.Bytes()
	}
	return out
}()

// acceptsGzip reports whether the client listed gzip in Accept-Encoding
// without refusing it (q=0).
func acceptsGzip(request *http.Request) bool {
	for _, coding := range strings.Split(request.Header.Get("Accept-Encoding"), ",") {
		name, params, _ := strings.Cut(strings.TrimSpace(coding), ";")
		if strings.EqualFold(strings.TrimSpace(name), "gzip") {
			weight := strings.ReplaceAll(strings.TrimSpace(params), " ", "")
			return weight != "q=0" && weight != "q=0.0"
		}
	}
	return false
}

// URL returns the cache-busted URL of an embedded asset.
func URL(name string) string {
	return routes.Static + name + "?v=" + versions[name]
}

// Handler serves the embedded assets under routes.Static. Requests carrying
// the content version are cached immutably. CSS and JavaScript go out
// gzipped to clients that accept it. Only the embedded files are served: no
// directory listings, 404 for anything else.
func Handler() http.Handler {
	files := http.StripPrefix(routes.Static, http.FileServerFS(FS))
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		name := strings.TrimPrefix(request.URL.Path, routes.Static)
		if _, ok := versions[name]; !ok {
			http.NotFound(writer, request)
			return
		}
		if value := request.URL.Query().Get("v"); value != "" && value == versions[name] {
			writer.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			writer.Header().Set("Cache-Control", "no-cache")
		}
		packed, ok := gzipped[name]
		if !ok {
			files.ServeHTTP(writer, request)
			return
		}
		writer.Header().Add("Vary", "Accept-Encoding")
		if !acceptsGzip(request) {
			files.ServeHTTP(writer, request)
			return
		}
		writer.Header().Set("Content-Encoding", "gzip")
		writer.Header().Set("Content-Type", mime.TypeByExtension(path.Ext(name)))
		http.ServeContent(writer, request, name, time.Time{}, bytes.NewReader(packed))
	})
}
