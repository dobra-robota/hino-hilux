// Command export renders every route in the shared registry to static files
// under dist/ and copies the assets tree in. dist/ is the artifact published to
// GitHub Pages. This is the phase-1 deploy path; the same registry also drives
// cmd/server, so there is no duplicated view or routing logic.
package main

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/dobra-robota/hino-hilux/internal/site"
)

const distDir = "dist"

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := os.RemoveAll(distDir); err != nil {
		return err
	}

	ctx := context.Background()
	for _, r := range site.Routes() {
		out := filepath.Join(distDir, filepath.FromSlash(r.Output))
		if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
			return err
		}
		f, err := os.Create(out)
		if err != nil {
			return err
		}
		if err := r.Page.Render(ctx, f); err != nil {
			f.Close()
			return err
		}
		if err := f.Close(); err != nil {
			return err
		}
		log.Printf("wrote %s", out)
	}

	if err := copyDir("assets", filepath.Join(distDir, "assets")); err != nil {
		return err
	}
	log.Printf("copied assets -> %s", filepath.Join(distDir, "assets"))

	// static/ holds site-root files (robots.txt, later favicon, CNAME, etc.)
	// copied to the dist root, not under /assets/.
	if err := copyDir("static", distDir); err != nil {
		return err
	}
	log.Printf("copied static -> %s", distDir)
	return nil
}

// copyDir recursively copies src to dst. Missing src is not an error.
func copyDir(src, dst string) error {
	info, err := os.Stat(src)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return copyFile(src, dst, info)
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}
	for _, e := range entries {
		s := filepath.Join(src, e.Name())
		d := filepath.Join(dst, e.Name())
		if e.IsDir() {
			if err := copyDir(s, d); err != nil {
				return err
			}
			continue
		}
		fi, err := e.Info()
		if err != nil {
			return err
		}
		if err := copyFile(s, d, fi); err != nil {
			return err
		}
	}
	return nil
}

func copyFile(src, dst string, info os.FileInfo) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}
