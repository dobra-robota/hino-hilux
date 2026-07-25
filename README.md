# hino-hilux

Brand/landing site for an **independent** truck workshop specialising in fast-turnaround
Hino 300 & Hino 700 servicing. Built with Go + [`templ`](https://templ.guide), rendered to
static HTML, and deployed to GitHub Pages.

Not affiliated with, endorsed by, or connected to Hino Motors Ltd. / Hino South Africa.
Model names are used only to describe the vehicles serviced.

Full design + decisions: [`docs/specs/website.md`](docs/specs/website.md).

## Architecture

One codebase, two entrypoints sharing a single route registry
(`internal/site/routes.go`) so the phase-1 static site and the future dynamic app never
diverge:

- `cmd/server` — live `net/http` server (local dev now; production runtime later).
- `cmd/export` — renders every route to `dist/` for GitHub Pages.

Views are typed `templ` components in `views/`. Styling is one hand-written stylesheet
(`assets/css/site.css`). Site-root files (e.g. `robots.txt`) live in `static/`.

## Prerequisites

- Go **1.25+** (the pinned `templ` version requires it; the toolchain auto-downloads if needed).

## Develop

`*_templ.go` is generated (gitignored), so generate first — the Makefile does it for you:

```sh
make dev      # generate + run the server at http://localhost:8080
make build    # generate + export the static site to dist/
make generate # just regenerate *_templ.go after editing a .templ file
make clean    # remove dist/
```

To preview the exact project-Pages output locally:

```sh
SITE_BASE_PATH=/hino-hilux make build
```

## Deploy

Pushing to `master` runs `.github/workflows/pages.yml`, which generates, exports (with
`SITE_BASE_PATH=/<repo>`), and publishes `dist/` to GitHub Pages.

One-time setup: **Settings → Pages → Build and deployment → Source = "GitHub Actions".**

> `static/robots.txt` (AI-crawler allowlist) only governs crawling once the site is on a
> custom domain — on the default `*.github.io/hino-hilux/` URL, `robots.txt` at the
> sub-path is ignored by crawlers.
