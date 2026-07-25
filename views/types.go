package views

// NavLink is a single top-level navigation entry.
type NavLink struct {
	Label string
	Path  string // site-root-relative path, e.g. "/" or "/services/"
}

// Layout carries everything the shared page shell needs to render one page.
// BasePath makes URLs correct under GitHub Pages' project sub-path (e.g.
// "/hino-hilux") while staying empty for the local server.
type Layout struct {
	Title       string
	Description string
	BasePath    string    // no trailing slash; "" locally, "/hino-hilux" on Pages
	Current     string    // the current page's Path, for active-nav marking
	Nav         []NavLink
}

// Href returns a navigable URL for a site-root-relative path, prefixed with BasePath.
func (l Layout) Href(p string) string { return l.BasePath + p }
