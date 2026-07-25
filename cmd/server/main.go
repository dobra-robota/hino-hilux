// Command server runs the site as a live HTTP server. In phase 1 it is a local
// dev convenience; in phase 2 it becomes the production runtime that also hosts
// the dynamic HTMX booking/ordering endpoints. It serves exactly the pages in
// the shared route registry.
package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/a-h/templ"

	"github.com/dobra-robota/hino-hilux/internal/site"
)

func main() {
	addr := ":" + port()

	mux := http.NewServeMux()
	for _, r := range site.Routes() {
		// Exact-match each path ("{$}") so "/" doesn't shadow subtrees.
		mux.Handle("GET "+r.Path+"{$}", templ.Handler(r.Page))
	}
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	// Site-root static files (robots.txt, etc.). Registered per file so they
	// don't shadow page routes. NOTE: on a GitHub Pages *project* site,
	// robots.txt at the sub-path is ignored by crawlers (host-root only);
	// it only governs crawling once the site is on a custom domain.
	if entries, err := os.ReadDir("static"); err == nil {
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			name := e.Name()
			mux.Handle("GET /"+name, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.ServeFile(w, r, filepath.Join("static", name))
			}))
		}
	}

	log.Printf("serving on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func port() string {
	if p := os.Getenv("PORT"); p != "" {
		return p
	}
	return "8080"
}
