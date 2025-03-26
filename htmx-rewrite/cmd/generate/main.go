package main

import (
	"bytes"
	"context"
	"fmt"
	"loficode/templates/pages/home"
	"loficode/templates/pages/notfound"
	"loficode/templates/pages/post"
	"loficode/templates/pages/privacypolicy"
	"loficode/templates/pages/tos"

	"io/fs"
	"loficode/templates/components"
	"os"
	"path/filepath"
	"strings"

	"github.com/a-h/templ"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {
	var posts []components.Post

	err := os.MkdirAll("public/posts", 0755)
	if err != nil {
		return fmt.Errorf("failed to create public/posts directory: %w", err)
	}

	err = filepath.WalkDir("cms/_posts", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk directory: %w", err)
		}
		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			fmt.Printf("Skipping: %s\n", path)
			return nil
		}
		p, err := parseMarkdownFile(path)
		if err != nil {
			return fmt.Errorf("failed to parse markdown file: %w", err)
		}
		posts = append(posts, *p)
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to walk directory: %w", err)
	}

	fmt.Printf("Parsed %d posts\n", len(posts))

	staticPages := map[string]templ.Component{
		"index.html":          home.Home(),
		"tos.html":            tos.TermsOfService(),
		"privacy-policy.html": privacypolicy.PrivacyPolicy(),
		"404.html":            notfound.NotFound(),
	}
	for _, p := range posts {
		htmlOut := filepath.Join("posts", p.Slug+".html")
		staticPages[htmlOut] = post.Post(p)
	}

	for name, component := range staticPages {
		err := renderStaticPage(name, component)
		if err != nil {
			return fmt.Errorf("failed to render static page: %w", err)
		}
	}

	return nil
}

func renderStaticPage(name string, component templ.Component) error {
	path := filepath.Join("public", name)
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create %s: %w", path, err)
	}
	defer f.Close()

	if err := component.Render(context.Background(), f); err != nil {
		return fmt.Errorf("failed to render %s: %w", path, err)
	}
	fmt.Println("Wrote:", path)
	return nil
}

func parseMarkdownFile(path string) (*components.Post, error) {
	input, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	context := parser.NewContext()
	md := goldmark.New(
		goldmark.WithExtensions(meta.Meta),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	var buf bytes.Buffer
	err = md.Convert(input, &buf, parser.WithContext(context))
	if err != nil {
		return nil, err
	}

	metaData := meta.Get(context)
	post := &components.Post{
		Content: buf.String(),
	}

	if title, ok := metaData["title"].(string); ok {
		post.Name = title
	}
	if slug, ok := metaData["slug"].(string); ok {
		post.Slug = slug
	}

	return post, nil
}
