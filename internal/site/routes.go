// Package site holds the shared route registry consumed by both cmd/server
// (live HTTP) and cmd/export (static generation). It is the single source of
// truth for what pages exist and what renders them.
package site

import (
	"os"
	"strings"

	"github.com/a-h/templ"

	"github.com/dobra-robota/hino-hilux/views"
)

// Route maps a URL path to the component that renders it and the file it
// exports to under dist/.
type Route struct {
	Path   string          // URL path, e.g. "/" or "/services/"
	Output string          // dist-relative output file, e.g. "index.html"
	Page   templ.Component // the fully-composed page component
}

// pageSpec declares one page. NavLabel != "" means the page appears in the
// top nav; keeping it empty (or simply not declaring a page) is what prevents
// dead nav links; the nav is derived from the pages that actually exist.
type pageSpec struct {
	Path        string
	Output      string
	NavLabel    string
	Title       string
	Description string
	Build       func(views.Layout) templ.Component
}

// pages is the ordered registry. Add a page here and it is served, exported,
// and (if NavLabel is set) linked in the nav, with no other file to touch.
var pages = []pageSpec{
	{
		Path:        "/",
		Output:      "index.html",
		NavLabel:    "Home",
		Title:       "Hino & Hilux Specialists, a family workshop",
		Description: "A family-run workshop with 50+ years of combined experience servicing Hino trucks and Hilux bakkies, with fast, honest turnaround for fleet operators.",
		Build:       views.HomePage,
	},
	{
		Path:        "/services/",
		Output:      "services/index.html",
		NavLabel:    "Services",
		Title:       "Services: Hino & Hilux servicing, repairs & fleet work",
		Description: "Specialist servicing, repairs and fleet maintenance for Hino 300 and Hino 700 trucks and the Toyota Hilux, with fast, honest turnaround.",
		Build:       views.ServicesPage,
	},
	{
		Path:        "/about/",
		Output:      "about/index.html",
		NavLabel:    "About us",
		Title:       "About us: a family Hino & Hilux workshop",
		Description: "A father-and-two-sons workshop with 50+ years of combined experience specialising in Hino trucks and the Toyota Hilux.",
		Build:       views.AboutPage,
	},
}

// basePath is the URL prefix pages live under. Empty for the local server;
// set to "/hino-hilux" (via SITE_BASE_PATH) for the GitHub Pages project site.
func basePath() string {
	return strings.TrimRight(os.Getenv("SITE_BASE_PATH"), "/")
}

// Routes returns every page in the site, with nav derived from the registry so
// the menu never links to a page that does not exist.
func Routes() []Route {
	base := basePath()

	var nav []views.NavLink
	for _, p := range pages {
		if p.NavLabel != "" {
			nav = append(nav, views.NavLink{Label: p.NavLabel, Path: p.Path})
		}
	}

	routes := make([]Route, 0, len(pages))
	for _, p := range pages {
		l := views.Layout{
			Title:       p.Title,
			Description: p.Description,
			BasePath:    base,
			Current:     p.Path,
			Nav:         nav,
		}
		routes = append(routes, Route{Path: p.Path, Output: p.Output, Page: p.Build(l)})
	}
	return routes
}
