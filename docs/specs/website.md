# Spec: hino-hilux brand site

Status: draft (spec-driven). Owner: repo maintainer. Source of truth for the build.

## Problem

`dobra-robota/hino-hilux` is empty. An **independent** truck-service business (the owner's
uncle's company) needs a web presence. Two horizons:

- **Phase 1 (now):** a landing / portfolio page positioning the business as **fast-turnaround
  Hino 300 & Hino 700 service experts** for fleet operators. Free to host, trivial to deploy.
- **Phase 2 (later):** interactive **service booking and ordering**, needing a request-time backend.

The phase-1 build must **not** be a throwaway: it must be structured so phase 2 adds dynamic
endpoints and a real runtime **without rewriting** the view or routing layer.

## Product context

- **Business:** independent specialist workshop. **Not affiliated with Hino Motors / Hino SA.**
  Uses model names (Hino 300, Hino 700) **nominatively** to describe what it services.
- **Audience:** fleet operators / businesses running Hino trucks (B2B lean).
- **Positioning:** speed + specialist expertise — "Hino 300 (light-duty) & Hino 700
  (heavy-duty) experts, very quick turnaround." Minimises fleet downtime.
- **Primary CTA (phase 1):** phone call (`tel:` link) + physical address + embedded map.
  No form/backend needed (booking is phase 2).
- **Copy:** realistic expert copy drafted now as a working starting point; owner edits later.
- **Reference for look/structure only:** https://www.hino.co.za/ — echo its clean industrial
  idiom (full-bleed truck hero, high whitespace, strong sans-serif, model-focused cards).
  Do NOT copy its logo/assets or imply affiliation.

## Goal

A **Go web application, server-first**, whose primary artifact is an HTTP server, plus a
**static-export mode** for phase-1 hosting.

- Views are typed **`templ` components** (Go); page content authored as **Markdown**.
- A **shared route registry** is consumed by two entrypoints:
  - `cmd/server` — a `net/http` server (local dev now; production runtime in phase 2).
  - `cmd/export` — renders every registered route to `dist/` for **GitHub Pages** (phase 1).
- Phase 1 ships **pure static HTML/CSS** — **no HTMX, no JS build**. HTMX arrives in phase 2
  with the booking/ordering endpoints (see D3).
- CI (**GitHub Actions**) runs the export and publishes `dist/` to Pages.

Observable end state: pushing to `master` publishes the updated static site at the project's
GitHub Pages URL in one CI run; the **same** codebase runs as a live server via
`go run ./cmd/server`, and in phase 2 on Cloudflare Containers / EC2.

## Non-goals (phase-1 blast-radius fence)

- **No dynamic endpoints, no HTMX, no JS in phase 1.** No booking/ordering handlers, forms,
  DB, or auth **yet** — architecture must accommodate them (D6).
- **No Node/JS build toolchain.** CI stays Go-only (plus Pages actions). No npm/Tailwind/bundler.
- **No CMS / admin UI.** Content changes are git commits.
- **No i18n / multi-locale** this iteration.
- **No custom domain / DNS** this iteration (default `*.github.io` URL). Later `CNAME` + DNS.
- **No production runtime provisioning** (Cloudflare/EC2) this iteration — phase 2.
- **No Hino branding/assets.** Own identity only; model names used nominatively.

## Decisions (resolved)

- **D1 — Server-first Go app with static export, not an SSG.** Views render through `templ`
  components (`github.com/a-h/templ`). A shared **route registry** maps each URL path to a
  component render func. `cmd/server` mounts them on `http.ServeMux`; `cmd/export` iterates
  the registry and writes each to `dist/`. Rejected Hugo / pure SSG (thrown away at phase 2)
  and frameworks (chi/echo/gin): stdlib `net/http` (Go 1.22 method+pattern routing) keeps the
  app portable to any host. `templ` chosen over stdlib `html/template` for type-safety and
  phase-2 HTMX ergonomics (owner's pick); `html/template` remains a valid lighter fallback.

- **D2 — Content: Markdown + front matter at build time.** Each `content/**/*.md` starts with
  a `---`-fenced **YAML front matter** block (`title`, `slug`, `description`, `order`, `draft`)
  followed by the Markdown body. The loader (`internal/site/content.go`) splits the leading
  `---...---` block, decodes it into a typed `FrontMatter` struct via `gopkg.in/yaml.v3`, and
  renders the remaining body to HTML with `github.com/yuin/goldmark`. (goldmark parses Markdown
  only — front matter is handled explicitly by the loader.) Path mirrors output:
  `content/index.md → /` → `dist/index.html`; `content/services.md → /services/` →
  `dist/services/index.html`.

- **D3 — HTMX deferred to phase 2.** Phase 1 is pure static HTML from Go templates; adding
  HTMX now buys nothing on a static host (nothing to call). In phase 2 the same `templ`
  components + `net/http` server gain `hx-get`/`hx-post` endpoints and a vendored `htmx.min.js`
  include — a one-line add, not a rewrite. Rejected wiring `hx-boost` now as pointless
  complexity on a static site.

- **D4 — Styling: hand-written CSS, single stylesheet** (`assets/css/site.css`, copied to
  `dist/assets/`). May **closely resemble** the reference (https://www.hino.co.za/):
  industrial/automotive with a Hino-like red accent on charcoal + white, strong sans-serif,
  full-bleed truck hero, generous whitespace, model-focused cards. The line not crossed:
  no Hino logo/official assets and no affiliation claim (see R3). All brand tokens (colours,
  fonts) centralised so a rebrand is a one-file change. Tailwind rejected (needs Node → non-goal).

- **D5 — Deploy: official GitHub Pages Actions flow.** `actions/configure-pages`,
  `actions/upload-pages-artifact` (uploads `dist/`), `actions/deploy-pages`. Source =
  "GitHub Actions". Concurrency-guarded so only the latest `master` push publishes.

- **D6 — Migration path is a design constraint.** The route registry + `templ` components are
  the portable core. Phase 2 = run `cmd/server` on **Cloudflare Containers or EC2/Fly/Cloud
  Run** (NOT Cloudflare Workers — Go-on-Workers is WASM-only and painful), add HTMX + dynamic
  routes beside the static ones. No view/routing rewrite. Static pages may stay on Pages while
  dynamic lives elsewhere.

- **D7 — AI-crawler allowlist via `robots.txt`.** `static/robots.txt` (copied to the site
  root) disallows all crawlers except an allowlist of trusted AI agents: OpenAI (`GPTBot`,
  `OAI-SearchBot`, `ChatGPT-User`), Anthropic (`ClaudeBot`, `Claude-SearchBot`, `Claude-User`),
  Google AI (`Google-Extended` + `Googlebot` — Gemini has no standalone crawler and grounds on
  Google's index; allowing `Googlebot` also restores normal Google Search), DeepSeek
  (`DeepSeekBot`) and Moonshot/Kimi (`Kimi-User`). The DeepSeek and Kimi tokens are from the
  community `ai-robots-txt` list, not first-party docs — marked unverified in the file.
  **Limitation:** on a GitHub Pages *project* site, `robots.txt` is only honoured at the host
  root (`<owner>.github.io/robots.txt`), NOT the sub-path — so this policy does not take effect
  on the default project URL. It governs crawling only once the site is on a **custom domain**
  (or a root user-site serves the file). Shipped now so it is ready; see R6.

- **D8 — Site-root static files + registry-derived nav.** A `static/` dir holds root files
  (robots.txt now; favicon/CNAME later), copied to `dist/` root by the exporter and served at
  `/` by the server. The nav is derived from the route registry (a page appears in the menu
  only if registered with a `NavLabel`), so no slice ever ships dead nav links.

## Site map (phase 1)

- **Home** (`/`) — hero (headline: Hino 300 & 700 experts, fast turnaround), quick value
  props, primary phone CTA, teaser of services.
- **Services** (`/services/`) — what they service: Hino 300 light-duty, Hino 700 heavy-duty,
  fleet servicing/maintenance, quick-turnaround repairs. Nominative model use.
- **About** (`/about/`) — independent specialist story, expertise, why fleets trust them.
- **Contact** (`/contact/`) — phone (`tel:`), address, hours, embedded map, non-affiliation note.

## Architecture

```
content/               Markdown + front matter (the copy)
  index.md
  services.md
  about.md
  contact.md
views/                 templ components (Go) -> *_templ.go
  base.templ           <html> shell: head, nav, footer (incl. non-affiliation notice)
  page.templ           renders a content page into base
internal/site/
  routes.go            the shared route registry: []Route{ Path, Render }
  content.go           markdown + front matter loading (goldmark)
cmd/server/main.go     net/http server: mounts registry on ServeMux (dev + phase-2 prod)
cmd/export/main.go     static export: iterate registry -> dist/ + copy assets
assets/                 copied to dist/assets/
  css/site.css         brand tokens + layout
  img/...              hero/truck imagery (placeholder now)
static/                 copied to dist/ ROOT (site-root files)
  robots.txt           AI-crawler allowlist (see D7)
dist/                  build output (gitignored); the Pages artifact
```

Build/run flows:

- **Local dev:** `templ generate && go run ./cmd/server` → serves the site at `localhost`.
- **Static export (CI):** `templ generate && go run ./cmd/export` → writes `dist/`; then
  upload + deploy to Pages.

## Acceptance criteria

Shared core / dual runtime:
- [ ] One route registry (`internal/site/routes.go`) is the only source of routes; both
      `cmd/server` and `cmd/export` consume it (no duplicated route lists). Nav is derived
      from the registry, so no page links to a route that does not exist.
- [ ] After `make generate` (regenerates gitignored `*_templ.go`), `go run ./cmd/server`
      serves every registered page over HTTP locally.
- [ ] `templ generate && go run ./cmd/export` produces `dist/` with zero errors, same pages as
      the server, at the D2 output paths.

Content / rendering:
- [ ] Markdown body renders to HTML (headings, links, lists, emphasis) inside the layout.
- [ ] `<title>` and `<meta name="description">` come from each page's front matter.
- [ ] `draft: true` pages are excluded from `dist/` (and the server).
- [ ] Nav lists top-level pages by front-matter `order`; the current page is marked active.
- [ ] `assets/` is copied to `dist/assets/`; links resolve under the Pages sub-path.
- [ ] No HTMX/JS is loaded; pages are fully functional as static HTML/CSS.

Product:
- [ ] Home hero states the Hino 300 & Hino 700 fast-turnaround expert positioning.
- [ ] Every page exposes the phone `tel:` CTA; Contact shows address + embedded map.
- [ ] A non-affiliation disclaimer ("independent; not affiliated with Hino") appears in the footer.
- [ ] No Hino logos/official assets are used.

Deploy:
- [ ] A push to `master` runs the workflow (`templ generate` → export → upload → deploy) green
      end to end, with `pages: write` + `id-token: write` perms and a concurrency group.
- [ ] The published Pages URL serves the home page and all non-draft pages with working CSS + nav.
- [ ] `static/robots.txt` is served at `dist/` root (and `/robots.txt` on the server) with the
      AI allowlist; policy is effective only under a custom domain (R6/D7).

## Risks & rollback

- **R1 — Pages sub-path breaks absolute asset/nav URLs.** Project Pages serve under `/<repo>/`.
  Mitigation (implemented): an explicit `BasePath` (from `SITE_BASE_PATH`) prefixes every
  asset/nav URL — empty locally, `/hino-hilux` in CI. Verify the deployed URL, not just `dist/`.
- **R2 — `templ`/Go version drift in CI.** Mitigation: pin the `templ` version, install the
  matching CLI in CI; `templ generate` in CI is authoritative.
- **R3 — Trademark / affiliation confusion.** Mitigation: no Hino logo/assets, own identity,
  nominative model use only, explicit footer disclaimer. Acceptance covers this.
- **R4 — Over-engineering phase 1.** Server + exporter split is deliberately minimal (stdlib
  `net/http`, few files). Justified solely by the phase-2 migration requirement.
- **R5 — Map embed pulls a third-party script.** A Google Maps `<iframe>` is static-safe (no
  build dep) but is an external request; acceptable, or swap for a static map image + address link.
- **R6 — `robots.txt` inert on the default project Pages URL.** REP honours `robots.txt` only
  at the host root, so `<owner>.github.io/hino-hilux/robots.txt` is ignored (see D7). The AI
  allowlist takes effect only under a custom domain / root user-site. Blocked on the custom
  domain, which is already a deferred non-goal. File shipped and ready.
- **Rollback:** static, versioned site. Revert the commit; CI redeploys the prior `dist/`.

## Proposed slices (tracer bullets)

Each slice leaves `master` green and deployable.

- **S1 — Walking skeleton to Pages, server-first.** `views/base.templ` + one page, the route
  registry, `cmd/server` (serves it) and `cmd/export` (writes `dist/index.html`), plus the full
  GitHub Actions deploy workflow. Acceptance: a live Pages URL renders the page AND
  `go run ./cmd/server` serves the same locally. (Proves the pipeline + dual-runtime seam —
  the highest-risk parts — first.)
- **S2 — Markdown content pipeline.** `goldmark` + front matter; `content/*.md` → routes; draft
  filtering; output-path mapping; server and exporter both use it.
- **S3 — Layout, nav, assets, real pages + copy.** `page.templ`, shared nav/footer (with
  non-affiliation notice + phone CTA), `assets/` copy + industrial `site.css`, and drafted copy
  for home/services/about/contact per the site map and positioning.
- **S4 — Contact & conversion.** `tel:` CTA everywhere, contact page with address/hours/map,
  SEO meta, favicon, 404 page.
- **S5 — Polish.** Responsive CSS pass, hero/placeholder imagery, accessibility check
  (landmarks, alt text, contrast).

## Phase 2 (out of scope now; recorded so phase 1 stays compatible)

Booking + ordering: add HTMX (`htmx.min.js`) + dynamic routes to the registry with
`hx-get`/`hx-post` handlers, a data store, and a production runtime (Cloudflare Containers /
EC2). Static marketing pages may stay on Pages or move behind the server. No view/routing
rewrite expected.
