# templ is pinned and run via `go run` so no separate install / PATH setup is needed.
TEMPL := go run github.com/a-h/templ/cmd/templ@v0.3.1020

.PHONY: generate dev build clean

# Generate *_templ.go from .templ sources (required before build/run).
generate:
	$(TEMPL) generate

# Run the live server locally at http://localhost:8080 (PORT overrides).
dev: generate
	go run ./cmd/server

# Build the static site into dist/. Set SITE_BASE_PATH for project Pages
# (CI sets it to /<repo>); leave empty for a custom-domain / root deploy.
build: generate
	go run ./cmd/export

clean:
	rm -rf dist
